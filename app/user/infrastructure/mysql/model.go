package mysql

import (
	"time"
)

type User struct {
	ID        int64      `json:"id"                   gorm:"primarykey"`
	Username  string     `json:"username"             gorm:"type:varchar(32);uniqueIndex;not null"`
	Password  string     `json:"-"                    gorm:"type:varchar(128);not null"`
	AvatarURL string     `json:"avatar_url"           gorm:"type:varchar(256)"`
	CreatedAt time.Time  `json:"created_at"           gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at"           gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
