package mysql

import (
	"context"
	"errors"

	"github.com/yxrxy/videoHub/app/user/domain/model"
	"github.com/yxrxy/videoHub/app/user/domain/repository"
	"github.com/yxrxy/videoHub/pkg/errno"
	"gorm.io/gorm"
)

type userDB struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) repository.UserDB {
	return &userDB{db: db}
}

func (u *userDB) IsUserExist(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := u.db.WithContext(ctx).Model(&User{}).
		Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsUserExistById 检查用户是否存在
func (u *userDB) IsUserExistById(ctx context.Context, id int64) (bool, error) {
	var count int64
	if err := u.db.WithContext(ctx).Model(&User{}).
		Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *userDB) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	dbUser := &User{
		Username:  user.Username,
		Password:  user.Password,
		AvatarURL: user.AvatarURL,
	}
	if err := u.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return 0, err
	}
	return dbUser.ID, nil
}

func (u *userDB) GetUserInfo(ctx context.Context, username string) (*model.User, error) {
	var user User
	if err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ParamVerifyErrorCode, "user not exist")
		}
		return nil, err
	}
	return &model.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		AvatarURL: user.AvatarURL,
	}, nil
}

func (u *userDB) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	var user User
	if err := u.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ParamVerifyErrorCode, "user not exist")
		}
		return nil, err
	}
	return &model.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		AvatarURL: user.AvatarURL,
	}, nil
}

func (u *userDB) UpdateUser(ctx context.Context, user *model.User) error {
	updates := map[string]interface{}{}
	if user.AvatarURL != "" {
		updates["avatar_url"] = user.AvatarURL
	}
	if user.Username != "" {
		updates["username"] = user.Username
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	if len(updates) == 0 {
		return nil
	}
	return u.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.ID).
		Updates(updates).Error
}

func (u *userDB) UpdateAvatar(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.ID).
		Update("avatar_url", user.AvatarURL).Error
}

func (u *userDB) ExistByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := u.db.WithContext(ctx).Model(&User{}).
		Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
