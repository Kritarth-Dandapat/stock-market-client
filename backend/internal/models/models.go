package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated platform user.
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	Phone        string     `json:"phone" db:"phone"`
	PasswordHash string     `json:"-" db:"password_hash"`
	IsVerified   bool       `json:"is_verified" db:"is_verified"`
	MFAEnabled   bool       `json:"mfa_enabled" db:"mfa_enabled"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`
}

// KYCStatus represents the KYC verification state for a user.
type KYCStatus string

const (
	KYCStatusPending    KYCStatus = "pending"
	KYCStatusInProgress KYCStatus = "in_progress"
	KYCStatusVerified   KYCStatus = "verified"
	KYCStatusRejected   KYCStatus = "rejected"
)

// KYCRecord stores KYC details for a user (PAN / Aadhaar stored encrypted).
type KYCRecord struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	EncryptedPAN  string     `json:"-" db:"encrypted_pan"`
	PANMasked     string     `json:"pan_masked" db:"pan_masked"`
	FullName      string     `json:"full_name" db:"full_name"`
	DateOfBirth   time.Time  `json:"date_of_birth" db:"date_of_birth"`
	Status        KYCStatus  `json:"status" db:"status"`
	CKYCNumber    string     `json:"ckyc_number,omitempty" db:"ckyc_number"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// OrderType distinguishes lumpsum purchases from SIP registrations.
type OrderType string

const (
	OrderTypeLumpsum OrderType = "lumpsum"
	OrderTypeSIP     OrderType = "sip"
	OrderTypeRedeem  OrderType = "redeem"
)

// OrderStatus represents the lifecycle state of an order.
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusSubmitted  OrderStatus = "submitted"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusFailed     OrderStatus = "failed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order represents a mutual fund purchase or redemption order.
type Order struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	UserID         uuid.UUID   `json:"user_id" db:"user_id"`
	SchemeCode     string      `json:"scheme_code" db:"scheme_code"`
	SchemeName     string      `json:"scheme_name" db:"scheme_name"`
	OrderType      OrderType   `json:"order_type" db:"order_type"`
	Amount         float64     `json:"amount" db:"amount"`
	Units          *float64    `json:"units,omitempty" db:"units"`
	NAV            *float64    `json:"nav,omitempty" db:"nav"`
	Status         OrderStatus `json:"status" db:"status"`
	BSEOrderID     string      `json:"bse_order_id,omitempty" db:"bse_order_id"`
	PaymentID      string      `json:"payment_id,omitempty" db:"payment_id"`
	FailureReason  string      `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
}

// SIPMandate stores a SIP / NACH mandate registered with BSE.
type SIPMandate struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	SchemeCode     string     `json:"scheme_code" db:"scheme_code"`
	SchemeName     string     `json:"scheme_name" db:"scheme_name"`
	Amount         float64    `json:"amount" db:"amount"`
	Frequency      string     `json:"frequency" db:"frequency"`
	StartDate      time.Time  `json:"start_date" db:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty" db:"end_date"`
	MandateID      string     `json:"mandate_id" db:"mandate_id"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Holding represents a user's current mutual fund holding.
type Holding struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	SchemeCode  string    `json:"scheme_code" db:"scheme_code"`
	SchemeName  string    `json:"scheme_name" db:"scheme_name"`
	ISIN        string    `json:"isin" db:"isin"`
	Units       float64   `json:"units" db:"units"`
	AvgNAV      float64   `json:"avg_nav" db:"avg_nav"`
	CurrentNAV  float64   `json:"current_nav" db:"current_nav"`
	InvestedAmt float64   `json:"invested_amount" db:"invested_amount"`
	CurrentAmt  float64   `json:"current_amount" db:"current_amount"`
	XIRR        *float64  `json:"xirr,omitempty" db:"xirr"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AuditLog records every significant mutation on the platform.
type AuditLog struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Action     string    `json:"action" db:"action"`
	Resource   string    `json:"resource" db:"resource"`
	ResourceID string    `json:"resource_id" db:"resource_id"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	Metadata   string    `json:"metadata" db:"metadata"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// FundScheme holds basic metadata for a mutual fund scheme.
type FundScheme struct {
	SchemeCode string  `json:"scheme_code"`
	ISIN       string  `json:"isin"`
	SchemeName string  `json:"scheme_name"`
	AMCCode    string  `json:"amc_code"`
	Category   string  `json:"category"`
	NAV        float64 `json:"nav"`
	NAVDate    string  `json:"nav_date"`
	MinPurchase float64 `json:"min_purchase"`
	MinSIP      float64 `json:"min_sip"`
}
