package service

import (
	"github.com/yxrxy/videoHub/app/user/domain/repository"
)

type UserService struct {
	db    repository.UserDB
	cache repository.UserCache
}

func NewUserService(db repository.UserDB, cache repository.UserCache) *UserService {
	if db == nil || cache == nil {
		panic("userService`s db or cache should not be nil")
	}
	svc := &UserService{
		db:    db,
		cache: cache,
	}
	return svc
}
