package service

import (
	"context"

	"github.com/yxrxy/videoHub/app/video/domain/repository"
)

type VideoService struct {
	db    repository.VideoDB
	cache repository.VideoCache
	mq    repository.VideoMQ
}

func NewVideoService(db repository.VideoDB, cache repository.VideoCache, mq repository.VideoMQ) *VideoService {
	if db == nil || cache == nil || mq == nil {
		panic("videoService`s db or cache or mq should not be nil")
	}
	svc := &VideoService{
		db:    db,
		cache: cache,
		mq:    mq,
	}
	svc.init()
	return svc
}

func (s *VideoService) init() {
	s.initConsumer()
}

func (s *VideoService) initConsumer() {
	go s.ConsumeProcessVideo(context.Background())
}
