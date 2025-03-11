package model

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID          int64     `gorm:"primarykey"`
	UserID      int64     `gorm:"not null"`
	Title       string    `gorm:"type:varchar(128);not null"`
	VideoURL    string    `gorm:"type:varchar(256);not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
