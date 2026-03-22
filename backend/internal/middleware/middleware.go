package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is an unexported type for request context keys.
type contextKey string

const claimsKey contextKey = "jwt_claims"

// Claims are the custom JWT claims embedded in every access token.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTMiddleware validates a Bearer JWT in the Authorization header.
// On success the parsed Claims are stored on the request context.
func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing Authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				writeError(w, http.StatusUnauthorized, "invalid Authorization header format")
				return
			}

			tokenStr := parts[1]
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext retrieves JWT claims stored by JWTMiddleware.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsKey).(*Claims)
	return c, ok
}

// RateLimiter is a simple in-memory token-bucket rate limiter keyed by IP.
// For production use, replace with a Redis-backed implementation.
type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new in-memory rate limiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Middleware returns an HTTP middleware that enforces the rate limit.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			ip = strings.Split(fwd, ",")[0]
		}

		now := time.Now()
		cutoff := now.Add(-rl.window)

		existing := rl.requests[ip]
		valid := make([]time.Time, 0, len(existing))
		for _, t := range existing {
			if t.After(cutoff) {
				valid = append(valid, t)
			}
		}
		rl.requests[ip] = valid

		if len(rl.requests[ip]) >= rl.limit {
			w.Header().Set("Retry-After", rl.window.String())
			writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		rl.requests[ip] = append(rl.requests[ip], now)
		next.ServeHTTP(w, r)
	})
}

// AuditLogger logs every mutation request to stdout as structured JSON.
// In production, replace with a write to the immutable audit_logs table.
type AuditLogger struct {
	logger interface {
		Info(msg string, args ...any)
	}
}

// NewAuditLogger creates an AuditLogger backed by the provided slog-compatible logger.
func NewAuditLogger(logger interface{ Info(msg string, args ...any) }) *AuditLogger {
	return &AuditLogger{logger: logger}
}

// Middleware logs metadata for non-GET/HEAD/OPTIONS requests.
func (al *AuditLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
			userID := ""
			if claims, ok := ClaimsFromContext(r.Context()); ok {
				userID = claims.UserID
			}
			al.logger.Info("audit",
				"user_id", userID,
				"method", r.Method,
				"path", r.URL.Path,
				"ip", r.RemoteAddr,
				"user_agent", r.Header.Get("User-Agent"),
				"time", time.Now().UTC().Format(time.RFC3339),
			)
		}
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders adds production-grade HTTP security headers to every response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'")
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg}) //nolint:errcheck
}
