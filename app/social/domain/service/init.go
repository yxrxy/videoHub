package service

import (
	"github.com/yxrxy/videoHub/app/social/domain/repository"
)

type SocialService struct {
	db        repository.SocialDB
	cache     repository.SocialCache
	wsService repository.SocialWebSocket
}

func NewSocialService(db repository.SocialDB, cache repository.SocialCache, wsService repository.SocialWebSocket) *SocialService {
	if db == nil || cache == nil || wsService == nil {
		panic("socialService`s db or cache or wsService should not be nil")
	}
	svc := &SocialService{
		db:    db,
		cache: cache,
	}
	return svc
}
