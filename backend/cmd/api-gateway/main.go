// Package main is the entrypoint for the API gateway service.
// The gateway handles JWT validation, rate limiting, audit logging,
// and proxies requests to the appropriate downstream Go services.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kritarth-Dandapat/stock-market-client/backend/internal/middleware"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/config"
	"github.com/Kritarth-Dandapat/stock-market-client/backend/pkg/logger"
)

func main() {
	log := logger.New()
	cfg := config.Load()

	mux := http.NewServeMux()

	// Health check — unauthenticated
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck
	})

	// Apply middleware chain: security headers → rate limiter → audit logger → JWT (on protected routes)
	rl := middleware.NewRateLimiter(100, time.Minute)
	al := middleware.NewAuditLogger(log)

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
		log.Info("api-gateway starting", slog.String("port", cfg.Server.Port))
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

	log.Info("shutting down gracefully")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error", slog.Any("err", err))
	}
}
