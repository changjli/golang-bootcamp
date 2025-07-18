package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           int    `gorm:"type:int;primary_key;"`
	Name         string `gorm:"type:varchar(255);not null"`
	Email        string `gorm:"type:varchar(255);unique;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
}
