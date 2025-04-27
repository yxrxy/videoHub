namespace go user

include "model.thrift"

// 注册请求
struct RegisterRequest {
    1: required string username   // 用户名
    2: required string password   // 密码
}

// 注册响应
struct RegisterResponse {
    1: required model.BaseResp Base // 基本响应信息
    2: required i64 user_id         // 用户ID
}

// 登录请求
struct LoginRequest {
    1: required string username   // 用户名
    2: required string password   // 密码
}

// 登录响应
struct LoginResponse {
    1: required model.BaseResp Base // 基本响应信息
    2: required i64 user_id         // 用户ID
    3: required string token        // 访问令牌
    4: required string refresh_token // 刷新令牌
}

// 获取用户信息请求
struct UserInfoRequest {
    1: required i64 user_id             // 用户ID
}

// 获取用户信息响应
struct UserInfoResponse {
    1: required model.BaseResp Base // 基本响应信息
    2: required model.User User         // 用户信息
}

// 上传头像请求
struct UploadAvatarRequest {
    1: required binary avatar_data     // 头像二进制数据
    2: required string content_type    // 文件类型
}

// 上传头像响应
struct UploadAvatarResponse {
    1: required model.BaseResp Base // 基本响应信息
    2: required string avatar_url      // 头像URL
}

// 刷新令牌请求
struct RefreshTokenRequest {
    1: required i64 user_id
}

// 刷新令牌响应
struct RefreshTokenResponse {
    1: required model.BaseResp Base // 基本响应信息
    2: required string token         // 新的访问令牌
}

// 用户服务
service UserService {
    // 用户注册
    RegisterResponse Register(1: RegisterRequest req)
    // 用户登录
    LoginResponse Login(1: LoginRequest req)
    // 获取用户信息
    UserInfoResponse GetUserInfo(1: UserInfoRequest req)
    // 上传头像
    UploadAvatarResponse UploadAvatar(1: UploadAvatarRequest req)
    // 刷新令牌
    RefreshTokenResponse RefreshToken(1: RefreshTokenRequest req)
}