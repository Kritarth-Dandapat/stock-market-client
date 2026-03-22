// Package main is the entrypoint for the KYC service.
// Handles CKYC lookup via CVL KRA and V-CIP video KYC initiation.
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

	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/crypto"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/middleware"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/models"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/config"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/logger"
	"github.com/google/uuid"
)

func main() {
	log := logger.New()
	cfg := config.Load()

	h := &kycHandler{cfg: cfg, log: log}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck
	})
	mux.Handle("POST /kyc/initiate", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.initiate)))
	mux.Handle("GET /kyc/status", middleware.JWTMiddleware(cfg.JWT.AccessSecret)(http.HandlerFunc(h.status)))

	rl := middleware.NewRateLimiter(20, time.Minute)
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
		log.Info("kyc-service starting", slog.String("port", cfg.Server.Port))
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

type kycHandler struct {
	cfg *config.Config
	log *slog.Logger
}

type kycInitiateRequest struct {
	PAN         string `json:"pan"`
	FullName    string `json:"full_name"`
	DateOfBirth string `json:"date_of_birth"` // YYYY-MM-DD
}

func (h *kycHandler) initiate(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req kycInitiateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.PAN == "" || req.FullName == "" || req.DateOfBirth == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "pan, full_name and date_of_birth are required"})
		return
	}

	// Encrypt PAN before persisting — key must be 32 bytes from KMS in production.
	key := make([]byte, 32) // placeholder — load from config.Crypto.PIIEncryptionKey in production
	encryptedPAN, err := crypto.EncryptPII(req.PAN, key)
	if err != nil {
		h.log.Error("PAN encryption failed", slog.Any("err", err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	record := &models.KYCRecord{
		ID:           uuid.New(),
		UserID:       uuid.MustParse(claims.UserID),
		EncryptedPAN: encryptedPAN,
		PANMasked:    crypto.MaskPAN(req.PAN),
		FullName:     req.FullName,
		Status:       models.KYCStatusInProgress,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	// In production: persist record, call CVL KRA CKYC lookup, initiate V-CIP if needed.
	h.log.Info("KYC initiated", slog.String("user_id", claims.UserID), slog.String("pan_masked", record.PANMasked))

	writeJSON(w, http.StatusAccepted, map[string]string{
		"kyc_id":     record.ID.String(),
		"pan_masked": record.PANMasked,
		"status":     string(record.Status),
	})
}

func (h *kycHandler) status(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	// In production: fetch KYC record from DB by user_id
	h.log.Info("KYC status requested", slog.String("user_id", claims.UserID))
	writeJSON(w, http.StatusOK, map[string]string{
		"user_id": claims.UserID,
		"status":  string(models.KYCStatusPending),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
