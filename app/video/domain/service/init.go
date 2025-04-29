package service

import (
	"context"

	"github.com/yxrxy/videoHub/app/user/domain/repository"
	videorepo "github.com/yxrxy/videoHub/app/video/domain/repository"
)

type VideoService struct {
	db     videorepo.VideoDB
	cache  videorepo.VideoCache
	mq     videorepo.VideoMQ
	es     videorepo.VideoElastic
	userDB repository.UserDB
}

func NewVideoService(db videorepo.VideoDB, cache videorepo.VideoCache, mq videorepo.VideoMQ, es videorepo.VideoElastic, userDB repository.UserDB) *VideoService {
	if db == nil || cache == nil || mq == nil || es == nil || userDB == nil {
		panic("videoService`s db or cache or mq or es or userDB should not be nil")
	}
	svc := &VideoService{
		db:     db,
		cache:  cache,
		mq:     mq,
		es:     es,
		userDB: userDB,
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
