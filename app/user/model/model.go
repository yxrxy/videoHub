package model

import (
	"time"
)

type User struct {
	ID        int64      `gorm:"primarykey" json:"id"`
	Username  string     `gorm:"type:varchar(32);uniqueIndex;not null" json:"username"`
	Password  string     `gorm:"type:varchar(128);not null" json:"-"`
	AvatarURL string     `gorm:"type:varchar(256)" json:"avatar_url"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}
