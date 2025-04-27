package rpc

import (
	"context"
	"log"

	"github.com/yxrxy/videoHub/kitex_gen/model"
	"github.com/yxrxy/videoHub/kitex_gen/user"
	"github.com/yxrxy/videoHub/pkg/base/client"
	"github.com/yxrxy/videoHub/pkg/errno"
)

// InitUserRPC 初始化用户服务客户端
func InitUserRPC() {
	c, err := client.InitUserRPC()
	if err != nil {
		log.Fatalf("初始化用户服务客户端失败: %v", err)
	}
	userClient = *c
}

// RegisterRPC 用户注册
func RegisterRPC(ctx context.Context, req *user.RegisterRequest) (int64, error) {
	resp, err := userClient.Register(ctx, req)
	if err != nil {
		log.Printf("注册用户RPC调用失败: %v", err)
		return 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.UserId, nil
}

// LoginRPC 用户登录
func LoginRPC(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	resp, err := userClient.Login(ctx, req)
	if err != nil {
		log.Printf("用户登录RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp, nil
}

// GetUserInfoRPC 获取用户信息
func GetUserInfoRPC(ctx context.Context, req *user.UserInfoRequest) (*model.User, error) {
	resp, err := userClient.GetUserInfo(ctx, req)
	if err != nil {
		log.Printf("获取用户信息RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.User, nil
}

// UploadAvatarRPC 上传头像
func UploadAvatarRPC(ctx context.Context, req *user.UploadAvatarRequest) (string, error) {
	resp, err := userClient.UploadAvatar(ctx, req)
	if err != nil {
		log.Printf("上传头像RPC调用失败: %v", err)
		return "", errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return "", errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.AvatarUrl, nil
}

// RefreshTokenRPC 刷新令牌
func RefreshTokenRPC(ctx context.Context, req *user.RefreshTokenRequest) (string, error) {
	resp, err := userClient.RefreshToken(ctx, req)
	if err != nil {
		log.Printf("刷新令牌RPC调用失败: %v", err)
		return "", errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return "", errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Token, nil
}
