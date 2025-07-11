package entities

import (
	"time"
)

// --- Custom Types for Enums ---

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
	PaymentExpired   PaymentStatus = "expired"
)

// --- GORM Entity ---

// Payment represents the payment model for GORM in the payment-service.
type Payment struct {
	ID          string        `gorm:"type:uuid;primary_key;"`
	UserID      int           `gorm:"type:uuid;not null;index"`
	MerchantID  string        `gorm:"type:varchar(255);not null"`
	Amount      float64       `gorm:"type:numeric(15,2);not null"`
	Status      PaymentStatus `gorm:"type:payment_status;not null;default:'pending';index"`
	Description string        `gorm:"type:varchar(255)"`
	CreatedAt   time.Time     `gorm:"index"`
	UpdatedAt   time.Time
}
