package entities

import (
	"core-service/domains/users/entities"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model               // Provides ID, CreatedAt, UpdatedAt, DeletedAt
	UserID     int           `gorm:"not null;unique"` // Must match the User's ID type
	Balance    float64       `gorm:"type:numeric(15,2);not null;default:0.00"`
	User       entities.User `gorm:"foreignKey:UserID"` // Defines the relationship
}
