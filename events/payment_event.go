package events

import "time"

type PaymentRequestedEvent struct {
	PaymentID   string    `json:"payment_id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	AuthToken   string    `json:"auth_token"`
	Description string    `json:"description"`
	MerchantID  string    `json:"merchant_id"`
}
