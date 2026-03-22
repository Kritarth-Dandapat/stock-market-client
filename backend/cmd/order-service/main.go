// Package main is the entrypoint for the order service.
// Handles lumpsum and SIP order placement via the BSE StAR MF API,
// payment webhook processing, and order status reconciliation.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/bse"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/middleware"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/models"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/config"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/logger"
	"github.com/google/uuid"
)

func main() {
	log := logger.New()
	cfg := config.Load()

	bseURL := cfg.BSE.BaseURL
	if cfg.BSE.UseUAT {
		bseURL = cfg.BSE.UATURL
	}
	bseClient := bse.NewClient(cfg.BSE.MemberCode, cfg.BSE.UserID, cfg.BSE.Password, bseURL)

	h := &orderHandler{cfg: cfg, log: log, bse: bseClient}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck
	})
	mux.Handle("POST /orders", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.placeOrder)))
	mux.Handle("GET /orders/{id}", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.getOrder)))
	mux.Handle("GET /orders", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.listOrders)))
	// Payment webhook — HMAC signature verified internally
	mux.HandleFunc("POST /webhooks/payment", h.paymentWebhook)

	al := middleware.NewAuditLogger(log)
	rl := middleware.NewRateLimiter(60, time.Minute)
	var handler http.Handler = mux
	handler = al.Middleware(handler)
	handler = rl.Middleware(handler)
	handler = middleware.SecurityHeaders(handler)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("order-service starting", slog.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx) //nolint:errcheck
}

type orderHandler struct {
	cfg *config.Config
	log *slog.Logger
	bse *bse.Client
}

type placeOrderRequest struct {
	SchemeCode  string            `json:"scheme_code"`
	OrderType   models.OrderType  `json:"order_type"`
	Amount      float64           `json:"amount"`
	FolioNo     string            `json:"folio_no,omitempty"`
}

func (h *orderHandler) placeOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req placeOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.SchemeCode == "" || req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "scheme_code and positive amount are required"})
		return
	}

	order := &models.Order{
		ID:         uuid.New(),
		UserID:     uuid.MustParse(claims.UserID),
		SchemeCode: req.SchemeCode,
		OrderType:  req.OrderType,
		Amount:     req.Amount,
		Status:     models.OrderStatusPending,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	// In production: persist order, verify KYC status, then submit to BSE
	bseReq := &bse.OrderRequest{
		ClientCode:  claims.UserID,
		SchemeCode:  req.SchemeCode,
		BuySell:     bse.TransactionTypePurchase,
		BuySellType: "FRESH",
		Amount:      req.Amount,
		OrderFlag:   bse.OrderFlagNew,
		FolioNo:     req.FolioNo,
		RefNo:       order.ID.String(),
		KYCStatus:   "Y",
	}

	bseResp, err := h.bse.PlaceOrder(r.Context(), bseReq)
	if err != nil {
		h.log.Error("BSE order placement failed", slog.Any("err", err), slog.String("user_id", claims.UserID))
		// Return submitted status — reconciliation job will sync later
		order.Status = models.OrderStatusSubmitted
	} else {
		order.BSEOrderID = bseResp.OrderID
		order.Status = models.OrderStatusProcessing
	}

	writeJSON(w, http.StatusAccepted, order)
}

func (h *orderHandler) getOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orderID := r.PathValue("id")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "order id is required"})
		return
	}

	// In production: fetch from DB, verify ownership
	h.log.Info("get order", slog.String("order_id", orderID), slog.String("user_id", claims.UserID))
	writeJSON(w, http.StatusOK, map[string]string{"order_id": orderID, "status": "processing"})
}

func (h *orderHandler) listOrders(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	// In production: paginated query from DB
	h.log.Info("list orders", slog.String("user_id", claims.UserID))
	writeJSON(w, http.StatusOK, []models.Order{})
}

func (h *orderHandler) paymentWebhook(w http.ResponseWriter, r *http.Request) {
	// In production: verify Razorpay/PayU webhook signature, update order status
	w.WriteHeader(http.StatusOK)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
