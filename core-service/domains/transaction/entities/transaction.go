package entities

import (
	"time"
)

// TransactionType defines the type of a transaction.
type TransactionType string

const (
	TopUp       TransactionType = "topup"
	TransferIn  TransactionType = "transfer_in"
	TransferOut TransactionType = "transfer_out"
	Payment     TransactionType = "payment"
)

// TransactionStatus defines the status of a transaction.
type TransactionStatus string

const (
	Completed TransactionStatus = "completed"
	Pending   TransactionStatus = "pending"
	Failed    TransactionStatus = "failed"
	Expired   TransactionStatus = "expired"
)

// Transaction represents the core transaction entity in the system.
// It's a single source of truth for all transaction records.
type Transaction struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	UserID      string            `json:"user_id"` // The user who owns this transaction record
	Type        TransactionType   `json:"type"`
	Amount      float64           `json:"amount"`
	Status      TransactionStatus `json:"status"`
	FromUserID  string            `json:"from_user_id"` // Source of funds (for transfers)
	ToUserID    string            `json:"to_user_id"`   // Destination of funds (for transfers)
	MerchantID  string            `json:"merchant_id"`  // Merchant ID for payments
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
