package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yxrxy/videoHub/app/user/domain/model"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/pkg/errno"
	"github.com/yxrxy/videoHub/pkg/jwt"
	"github.com/yxrxy/videoHub/pkg/storage"
)

func (s *UserService) Register(ctx context.Context, req *model.User) (int64, error) {
	passwordHash := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))

	newUser := &model.User{
		Username:  req.Username,
		Password:  passwordHash,
		AvatarURL: config.User.DefaultAvatar,
	}

	id, err := s.db.CreateUser(ctx, newUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserService) Login(ctx context.Context, u *model.User) (*model.User, error) {
	// 从数据库获取用户信息
	dbUser, err := s.db.GetUserInfo(ctx, u.Username)
	if err != nil {
		return nil, err
	}

	// 计算输入密码的哈希值
	inputPasswordHash := fmt.Sprintf("%x", md5.Sum([]byte(u.Password)))

	// 比较密码哈希值
	if inputPasswordHash != dbUser.Password {
		return nil, errno.NewErrNo(errno.ServiceWrongPassword, "密码错误")
	}

	token, err := jwt.GenerateToken(dbUser.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(dbUser.ID)
	if err != nil {
		return nil, err
	}

	// 将令牌存入缓存
	if err := s.cache.SetUserAccessToken(ctx, fmt.Sprintf("access_token:%d", dbUser.ID), token); err != nil {
		// 这里我们只记录错误但不返回，因为主要流程已经成功
		log.Printf("failed to cache access token: %v", err)
	}
	if err := s.cache.SetUserRefreshToken(ctx, fmt.Sprintf("refresh_token:%d", dbUser.ID), refreshToken); err != nil {
		log.Printf("failed to cache refresh token: %v", err)
	}

	return &model.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		AvatarURL:    dbUser.AvatarURL,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) UploadAvatar(ctx context.Context, uid int64, avatarData []byte, contentType string) (string, error) {
	if !isValidImageType(contentType) {
		return "", errno.NewErrNo(errno.InternalServiceErrorCode, "invalid param")
	}

	fileName := fmt.Sprintf("avatar_%d_%d.%s", uid, time.Now().Unix(), getFileExt(contentType))

	storage := storage.NewLocalStorage(config.Upload.Avatar.UploadDir, config.Upload.Avatar.BaseURL)
	avatarURL, err := storage.Save(avatarData, fileName)
	if err != nil {
		return "", err
	}

	if err := s.db.UpdateAvatar(ctx, &model.User{
		ID:        uid,
		AvatarURL: avatarURL,
	}); err != nil {
		return "", err
	}

	return avatarURL, nil
}

func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	return validTypes[contentType]
}

func getFileExt(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	default:
		return "jpg"
	}
}

func (s *UserService) RefreshToken(ctx context.Context, uid int64) (string, error) {
	refreshToken, err := s.cache.GetToken(ctx, fmt.Sprintf("refresh_token:%d", uid))
	if err != nil {
		if errors.Is(err, redis.Nil) { // token不存在或已过期
			return "", errno.NewErrNo(errno.AuthRefreshExpiredCode, "refresh token expired")
		}
		// Redis连接错误等其他错误
		return "", errno.NewErrNo(errno.InternalRedisErrorCode, "failed to get refresh token")
	}

	claims, err := jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	newToken, err := jwt.GenerateToken(claims.UserID)
	if err != nil {
		return "", err
	}

	return newToken, nil
}
