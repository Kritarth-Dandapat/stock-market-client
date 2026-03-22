// Package main is the entrypoint for the portfolio service.
// Syncs NAV data, calculates holdings, and exposes portfolio summaries.
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

	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/middleware"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/models"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/config"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/logger"
)

func main() {
	log := logger.New()
	cfg := config.Load()

	h := &portfolioHandler{cfg: cfg, log: log}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck
	})
	mux.Handle("GET /portfolio", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.getPortfolio)))
	mux.Handle("GET /portfolio/holdings", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.getHoldings)))
	mux.Handle("GET /funds", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.getFunds)))
	mux.Handle("GET /funds/{schemeCode}", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.getFundByScheme)))

	rl := middleware.NewRateLimiter(120, time.Minute)
	var handler http.Handler = mux
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
		log.Info("portfolio-service starting", slog.String("port", cfg.Server.Port))
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

type portfolioHandler struct {
	cfg *config.Config
	log *slog.Logger
}

type portfolioSummary struct {
	TotalInvested   float64          `json:"total_invested"`
	CurrentValue    float64          `json:"current_value"`
	AbsoluteReturn  float64          `json:"absolute_return"`
	XIRR            float64          `json:"xirr"`
	Holdings        []models.Holding `json:"holdings"`
}

func (h *portfolioHandler) getPortfolio(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	h.log.Info("get portfolio", slog.String("user_id", claims.UserID))

	// In production: aggregate holdings from DB, fetch latest NAV
	summary := portfolioSummary{
		TotalInvested:  0,
		CurrentValue:   0,
		AbsoluteReturn: 0,
		XIRR:           0,
		Holdings:       []models.Holding{},
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *portfolioHandler) getHoldings(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	h.log.Info("get holdings", slog.String("user_id", claims.UserID))

	writeJSON(w, http.StatusOK, []models.Holding{})
}

func (h *portfolioHandler) getFunds(w http.ResponseWriter, r *http.Request) {
	// In production: query fund scheme DB, support filtering/search
	writeJSON(w, http.StatusOK, []models.FundScheme{})
}

func (h *portfolioHandler) getFundByScheme(w http.ResponseWriter, r *http.Request) {
	schemeCode := r.PathValue("schemeCode")
	if schemeCode == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "schemeCode is required"})
		return
	}
	// In production: fetch scheme details + NAV history from DB
	writeJSON(w, http.StatusOK, models.FundScheme{SchemeCode: schemeCode})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
