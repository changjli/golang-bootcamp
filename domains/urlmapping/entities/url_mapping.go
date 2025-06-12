package entities

import "time"

type URLMapping struct {
	ID        uint      `gorm:"primaryKey"`
	ShortCode string    `gorm:"size:16;not null;uniqueIndex"`
	LongURL   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt *time.Time
}
