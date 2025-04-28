package service

import (
	"context"

	"github.com/yxrxy/videoHub/app/video/domain/repository"
)

type VideoService struct {
	db    repository.VideoDB
	cache repository.VideoCache
	mq    repository.VideoMQ
	es    repository.VideoElastic
}

func NewVideoService(db repository.VideoDB, cache repository.VideoCache, mq repository.VideoMQ, es repository.VideoElastic) *VideoService {
	if db == nil || cache == nil || mq == nil || es == nil {
		panic("videoService`s db or cache or mq should not be nil")
	}
	svc := &VideoService{
		db:    db,
		cache: cache,
		mq:    mq,
		es:    es,
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
