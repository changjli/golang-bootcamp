package entities

import (
	"core-service/domains/users/entities"

	"gorm.io/gorm"
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
	gorm.Model
	ID          int               `gorm:"type:uuid;primary_key;"`
	UserID      int               `gorm:"type:uuid;not null"`
	Type        TransactionType   `gorm:"type:string;not null"`
	Amount      float64           `gorm:"type:numeric(15,2);not null"`
	Status      TransactionStatus `gorm:"type:string;not null"`
	FromUserID  int               `gorm:"type:uuid"`
	ToUserID    int               `gorm:"type:uuid"`
	Description string            `gorm:"type:varchar(255)"`
	User        entities.User     `gorm:"foreignKey:UserID"`
	FromUser    entities.User     `gorm:"foreignKey:FromUserID"`
	ToUser      entities.User     `gorm:"foreignKey:ToUserID"`
}
