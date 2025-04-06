namespace go api.user

include "../user.thrift"

// API 服务
service UserAPI {
    // 用户账户相关接口
    user.RegisterResponse Register(1: user.RegisterRequest request) (api.post="/api/v1/user/register")
    user.LoginResponse Login(1: user.LoginRequest request) (api.post="/api/v1/user/login")
    user.RefreshTokenResponse RefreshToken(1: user.RefreshTokenRequest request) (api.post="/api/v1/user/token/refresh")
    
    // 用户信息相关接口
    user.UserInfoResponse GetUserInfo(1: user.UserInfoRequest request) (api.get="/api/v1/user/info")
    user.UploadAvatarResponse UploadAvatar(1: user.UploadAvatarRequest request) (api.post="/api/v1/user/avatar")
} 