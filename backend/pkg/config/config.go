package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	BSE      BSEConfig
	Crypto   CryptoConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// JWTConfig holds JWT signing and expiry settings.
type JWTConfig struct {
	AccessSecret        string
	RefreshSecret       string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
}

// BSEConfig holds BSE StAR MF API credentials.
type BSEConfig struct {
	MemberCode string
	UserID     string
	Password   string
	BaseURL    string
	UATURL     string
	UseUAT     bool
}

// CryptoConfig holds AES encryption key (32 bytes for AES-256).
type CryptoConfig struct {
	// PIIEncryptionKey is a 32-byte hex-encoded AES-256 key stored in secrets manager.
	PIIEncryptionKey string
}

// Load reads configuration from environment variables with sane defaults.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             getEnv("DATABASE_DSN", ""),
			MaxOpenConns:    getInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:       getEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret:      getEnv("JWT_REFRESH_SECRET", ""),
			AccessTokenExpiry:  getDuration("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDuration("JWT_REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
		},
		BSE: BSEConfig{
			MemberCode: getEnv("BSE_MEMBER_CODE", ""),
			UserID:     getEnv("BSE_USER_ID", ""),
			Password:   getEnv("BSE_PASSWORD", ""),
			BaseURL:    getEnv("BSE_BASE_URL", "https://www.bsestarmf.in/RMFSWebService/BseStarMFWebService.svc"),
			UATURL:     getEnv("BSE_UAT_URL", "https://bsestarmfuat.in/RMFSWebService/BseStarMFWebService.svc"),
			UseUAT:     getBool("BSE_USE_UAT", true),
		},
		Crypto: CryptoConfig{
			PIIEncryptionKey: getEnv("PII_ENCRYPTION_KEY", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}

func getBool(key string, defaultValue bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue
	}
	return b
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultValue
	}
	return d
}
