// Package main is the entrypoint for the auth service.
// Handles user registration, OTP-based login, JWT issuance and refresh.
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
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/config"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func main() {
	log := logger.New()
	cfg := config.Load()

	h := &authHandler{cfg: cfg, log: log}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /auth/register", h.register)
	mux.HandleFunc("POST /auth/login/otp/request", h.requestOTP)
	mux.HandleFunc("POST /auth/login/otp/verify", h.verifyOTP)
	mux.HandleFunc("POST /auth/token/refresh", h.refreshToken)

	rl := middleware.NewRateLimiter(30, time.Minute)
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
		log.Info("auth-service starting", slog.String("port", cfg.Server.Port))
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

type authHandler struct {
	cfg *config.Config
	log *slog.Logger
}

type registerRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type otpRequest struct {
	Phone string `json:"phone"`
}

type otpVerifyRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (h *authHandler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Email == "" || req.Phone == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and phone are required"})
		return
	}

	// In production: persist user to DB, send email verification, etc.
	userID := uuid.New().String()
	h.log.Info("user registered", slog.String("user_id", userID), slog.String("email", req.Email))

	writeJSON(w, http.StatusCreated, map[string]string{
		"user_id": userID,
		"message": "Registration successful. Please verify your phone number via OTP.",
	})
}

func (h *authHandler) requestOTP(w http.ResponseWriter, r *http.Request) {
	var req otpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Phone == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "phone is required"})
		return
	}

	// In production: generate OTP, store in Redis with TTL, dispatch via AWS SNS / Gupshup.
	h.log.Info("OTP requested", slog.String("phone", req.Phone))
	writeJSON(w, http.StatusOK, map[string]string{"message": "OTP sent"})
}

func (h *authHandler) verifyOTP(w http.ResponseWriter, r *http.Request) {
	var req otpVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Phone == "" || req.OTP == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "phone and otp are required"})
		return
	}

	// In production: validate OTP from Redis, look up user ID by phone.
	userID := uuid.New().String()

	accessToken, err := h.issueToken(userID, h.cfg.JWT.AccessSecret, h.cfg.JWT.AccessTokenExpiry)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not issue token"})
		return
	}
	refreshToken, err := h.issueToken(userID, h.cfg.JWT.RefreshSecret, h.cfg.JWT.RefreshTokenExpiry)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not issue refresh token"})
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(h.cfg.JWT.AccessTokenExpiry.Seconds()),
	})
}

func (h *authHandler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "refresh_token is required"})
		return
	}

	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(body.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.JWT.RefreshSecret), nil
	})
	if err != nil || !token.Valid {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
		return
	}

	accessToken, err := h.issueToken(claims.UserID, h.cfg.JWT.AccessSecret, h.cfg.JWT.AccessTokenExpiry)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not issue token"})
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(h.cfg.JWT.AccessTokenExpiry.Seconds()),
	})
}

func (h *authHandler) issueToken(userID, secret string, expiry time.Duration) (string, error) {
	claims := &middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
