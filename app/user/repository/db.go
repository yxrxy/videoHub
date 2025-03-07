package repository

import (
	"context"

	"github.com/yxrrxy/videoHub/app/user/model"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *User) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrUserNotExist
		}
		return nil, err
	}
	return &user, nil
}

func (u *User) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrUserNotExist
		}
		return nil, err
	}
	return &user, nil
}

func (u *User) UpdateAvatar(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).
		Update("avatar_url", user.AvatarURL).Error
}

func (u *User) UpdateMFASecret(ctx context.Context, id int64, secret string) error {
	return u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Update("mfa_secret", secret).Error
}

func (u *User) ExistByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := u.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
