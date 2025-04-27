package repository

import (
	"context"

	"github.com/yxrxy/videoHub/app/user/domain/model"
)

type UserDB interface {
	IsUserExist(ctx context.Context, username string) (bool, error)
	IsUserExistById(ctx context.Context, id int64) (bool, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUserInfo(ctx context.Context, username string) (*model.User, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	UpdateAvatar(ctx context.Context, user *model.User) error
}

type UserCache interface {
	IsExist(ctx context.Context, key string) bool
	SetUserAccessToken(ctx context.Context, key string, token string) error
	SetUserRefreshToken(ctx context.Context, key string, token string) error
	DeleteUserToken(ctx context.Context, key string) error
	GetToken(ctx context.Context, key string) (string, error)
}
