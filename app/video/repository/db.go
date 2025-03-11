package repository

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/yxrrxy/videoHub/app/video/model"
	"github.com/yxrrxy/videoHub/config"
)

type Video struct {
	db *gorm.DB
}

func NewVideo(db *gorm.DB) *Video {
	return &Video{db: db}
}

func (v *Video) CreateVideo(ctx context.Context, video *model.Video) error {
	return v.db.WithContext(ctx).Create(video).Error
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
