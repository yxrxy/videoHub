package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/yxrrxy/videoHub/app/user/model"
	"github.com/yxrrxy/videoHub/app/user/repository"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/user"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/jwt"
	"github.com/yxrrxy/videoHub/pkg/storage"
)

type UserService struct {
	repo *repository.User
}

func NewUserService(repo *repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	exist, err := s.repo.ExistByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errno.ErrUserAlreadyExist
	}

	passwordHash := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))

	newUser := &model.User{
		Username:  req.Username,
		Password:  passwordHash,
		AvatarURL: config.User.DefaultAvatar,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(newUser.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(newUser.ID)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{
		Base: &user.BaseResp{
			Code:    errno.Success.ErrCode,
			Message: errno.Success.ErrMsg,
		},
		UserId:       newUser.ID,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	userModel, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	passwordHash := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
	if passwordHash != userModel.Password {
		return nil, errno.ErrPasswordWrong
	}

	token, err := jwt.GenerateToken(userModel.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(userModel.ID)
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		Base: &user.BaseResp{
			Code:    errno.Success.ErrCode,
			Message: errno.Success.ErrMsg,
		},
		UserId:       userModel.ID,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	userModel, err := s.repo.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var deletedAt int64
	if userModel.DeletedAt != nil {
		deletedAt = userModel.DeletedAt.Unix()
	}

	return &user.UserInfoResponse{
		Base: &user.BaseResp{
			Code:    errno.Success.ErrCode,
			Message: errno.Success.ErrMsg,
		},
		User: &user.User{
			Id:        userModel.ID,
			Username:  userModel.Username,
			AvatarUrl: userModel.AvatarURL,
			CreatedAt: userModel.CreatedAt.Unix(),
			UpdatedAt: userModel.UpdatedAt.Unix(),
			DeletedAt: deletedAt,
		},
	}, nil
}

func (s *UserService) UploadAvatar(ctx context.Context, req *user.UploadAvatarRequest) (*user.UploadAvatarResponse, error) {
	userID, ok := pkgcontext.GetUserID(ctx)
	if !ok {
		return nil, errno.ErrUnauthorized
	}

	contentType := req.ContentType
	if !isValidImageType(contentType) {
		return nil, errno.ErrInvalidParam
	}

	fileName := fmt.Sprintf("avatar_%d_%d.%s", userID, time.Now().Unix(), getFileExt(contentType))

	storage := storage.NewLocalStorage(config.Upload.Avatar.UploadDir, config.Upload.Avatar.BaseURL)
	avatarURL, err := storage.Save(req.AvatarData, fileName)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateAvatar(ctx, &model.User{
		ID:        userID,
		AvatarURL: avatarURL,
	}); err != nil {
		return nil, err
	}

	return &user.UploadAvatarResponse{
		Base: &user.BaseResp{
			Code:    errno.Success.ErrCode,
			Message: errno.Success.ErrMsg,
		},
		AvatarUrl: avatarURL,
	}, nil
}

// 辅助函数：检查文件类型
func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	return validTypes[contentType]
}

// 辅助函数：获取文件扩展名
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

func (s *UserService) GetMFAQRCode(ctx context.Context, req *user.MFAQRCodeRequest) (*user.MFAQRCodeResponse, error) {
	// TODO: 使用第三方库生成MFA密钥和二维码
	qrCode := "data:image/png;base64,..."

	return &user.MFAQRCodeResponse{
		Base: &user.BaseResp{
			Code:    errno.Success.ErrCode,
			Message: errno.Success.ErrMsg,
		},
		QrCode: qrCode,
	}, nil
}

func (s *UserService) BindMFA(ctx context.Context, req *user.BindMFARequest) (*user.BindMFAResponse, error) {
	// TODO: 验证MFA码并保存密钥
	if err := s.repo.UpdateMFASecret(ctx, req.UserId, "mfa_secret"); err != nil {
		return nil, err
	}

	return &user.BindMFAResponse{
		Base: &user.BaseResp{
			Code:    errno.Success.ErrCode,
			Message: errno.Success.ErrMsg,
		},
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req *user.RefreshTokenRequest) (*user.RefreshTokenResponse, error) {
	claims, err := jwt.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errno.ErrInvalidToken
	}

	newToken, err := jwt.GenerateToken(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &user.RefreshTokenResponse{
		Base: &user.BaseResp{
			Code:    errno.Success.ErrCode,
			Message: errno.Success.ErrMsg,
		},
		Token: newToken,
	}, nil
}
