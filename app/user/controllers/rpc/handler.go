package rpc

import (
	"context"

	"github.com/yxrxy/videoHub/app/user/controllers/rpc/pack"
	"github.com/yxrxy/videoHub/app/user/domain/model"
	"github.com/yxrxy/videoHub/app/user/usecase"
	"github.com/yxrxy/videoHub/kitex_gen/user"
	"github.com/yxrxy/videoHub/pkg/base"
	uc "github.com/yxrxy/videoHub/pkg/base/context"
)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (r *user.RegisterResponse, err error) {
	r = new(user.RegisterResponse)
	u := &model.User{
		Username: req.Username,
		Password: req.Password,
	}

	var uid int64
	if uid, err = h.useCase.Register(ctx, u); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	r.UserId = uid
	return
}

func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (r *user.LoginResponse, err error) {
	r = new(user.LoginResponse)
	u := &model.User{
		Username: req.Username,
		Password: req.Password,
	}

	var user *model.User
	if user, err = h.useCase.Login(ctx, u); err != nil {
		return
	}
	r.UserId = user.ID
	r.Token = user.Token
	r.RefreshToken = user.RefreshToken
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *UserHandler) GetUserInfo(ctx context.Context, req *user.UserInfoRequest) (r *user.UserInfoResponse, err error) {
	r = new(user.UserInfoResponse)

	var user *model.User
	if user, err = h.useCase.GetUserInfo(ctx, req.UserId); err != nil {
		return
	}
	r.User = pack.User(user)
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *UserHandler) UploadAvatar(ctx context.Context, req *user.UploadAvatarRequest) (r *user.UploadAvatarResponse, err error) {
	r = new(user.UploadAvatarResponse)
	uid, err := uc.GetUserID(ctx)
	if err != nil {
		return
	}

	var avatarURL string
	if avatarURL, err = h.useCase.UploadAvatar(ctx, uid, req.AvatarData, req.ContentType); err != nil {
		return
	}
	r.AvatarUrl = avatarURL
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *user.RefreshTokenRequest) (r *user.RefreshTokenResponse, err error) {
	r = new(user.RefreshTokenResponse)

	var token string
	if token, err = h.useCase.RefreshToken(ctx, req.UserId); err != nil {
		return
	}
	r.Token = token
	r.Base = base.BuildBaseResp(err)
	return
}
