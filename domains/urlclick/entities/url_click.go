package entities

import "time"

type URLClick struct {
	ID        uint      `gorm:"primaryKey"`
	MappingID uint      `gorm:"not null;index"`
	ClickedAt time.Time `gorm:"autoCreateTime"`
	IPAddress string    `gorm:"size:45"`
	UserAgent string    `gorm:"type:text"`
}
