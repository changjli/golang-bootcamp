package entities

import "time"

type Wallet struct {
	ID        string    `json:"id"`         // Unique identifier for the wallet
	UserID    string    `json:"user_id"`    // The ID of the user who owns the wallet
	Balance   float64   `json:"balance"`    // The current balance of the wallet
	CreatedAt time.Time `json:"created_at"` // Timestamp when the wallet was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the wallet was last updated
}
