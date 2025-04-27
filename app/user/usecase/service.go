package usecase

import (
	"context"

	"github.com/yxrxy/videoHub/app/user/domain/model"
	"github.com/yxrxy/videoHub/pkg/errno"
)

func (s *useCase) Register(ctx context.Context, req *model.User) (int64, error) {
	exist, err := s.db.IsUserExist(ctx, req.Username)
	if err != nil {
		return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "database error")
	}
	if exist {
		return 0, errno.NewErrNo(errno.ServiceUserExist, "user already exist")
	}

	return s.svc.Register(ctx, req)
}

func (s *useCase) Login(ctx context.Context, req *model.User) (*model.User, error) {
	exist, err := s.db.IsUserExist(ctx, req.Username)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "database error")
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceUserNotExist, "user not exist")
	}

	return s.svc.Login(ctx, req)
}

func (s *useCase) GetUserInfo(ctx context.Context, uid int64) (*model.User, error) {
	return s.db.GetUserById(ctx, uid)
}

func (s *useCase) UploadAvatar(ctx context.Context, uid int64, avatarData []byte, contentType string) (string, error) {
	exist, err := s.db.IsUserExistById(ctx, uid)
	if err != nil {
		return "", errno.NewErrNo(errno.InternalDatabaseErrorCode, "database error")
	}
	if !exist {
		return "", errno.NewErrNo(errno.ServiceUserNotExist, "user not exist")
	}

	return s.svc.UploadAvatar(ctx, uid, avatarData, contentType)
}

func (s *useCase) RefreshToken(ctx context.Context, uid int64) (string, error) {
	return s.svc.RefreshToken(ctx, uid)
}
