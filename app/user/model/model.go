package model

import (
	"time"
)

type User struct {
	ID        int64     `gorm:"primarykey"`
	Username  string    `gorm:"type:varchar(32);uniqueIndex;not null"`
	Password  string    `gorm:"type:varchar(128);not null"`
	AvatarURL string    `gorm:"type:varchar(256)"`
	MFASecret string    `gorm:"type:varchar(32)"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
