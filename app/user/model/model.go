package model

import (
	"time"
)

type User struct {
	ID        int64      `gorm:"primarykey"`
	Username  string     `gorm:"type:varchar(32);uniqueIndex;not null"`
	Password  string     `gorm:"type:varchar(128);not null"`
	AvatarURL string     `gorm:"type:varchar(256)"`
	MFASecret string     `gorm:"type:varchar(32)"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
