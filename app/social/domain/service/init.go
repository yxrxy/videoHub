package service

import (
	"github.com/yxrxy/videoHub/app/social/domain/repository"
)

type SocialService struct {
	db    repository.SocialDB
	cache repository.SocialCache
}

func NewSocialService(db repository.SocialDB, cache repository.SocialCache) *SocialService {
	if db == nil || cache == nil {
		panic("socialService`s db or cache should not be nil")
	}
	svc := &SocialService{
		db:    db,
		cache: cache,
	}
	return svc
}
