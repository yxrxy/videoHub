namespace go user

// 用户信息
struct User {
    1: required i64 id                    // 用户ID
    2: required string username           // 用户名
    3: required string avatar_url         // 头像URL
    4: required i64 created_at            // 创建时间
    5: required i64 updated_at            // 更新时间
    6: optional i64 deleted_at            // 删除时间
}

// 注册请求
struct RegisterRequest {
    1: required string username (vt.min_size = "1", vt.max_size = "32")    // 用户名
    2: required string password (vt.min_size = "6", vt.max_size = "32")    // 密码
}

// 注册响应
struct RegisterResponse {
    1: required i64 user_id              // 用户ID
    2: required string token             // 访问令牌
    3: required string refresh_token     // 刷新令牌
}

// 登录请求
struct LoginRequest {
    1: required string username (vt.min_size = "1")    // 用户名
    2: required string password (vt.min_size = "6")    // 密码
}

// 登录响应
struct LoginResponse {
    1: required i64 user_id             // 用户ID
    2: required string token            // 访问令牌
    3: required string refresh_token    // 刷新令牌
}

// 获取用户信息请求
struct UserInfoRequest {
    1: required i64 user_id             // 用户ID
}

// 获取用户信息响应
struct UserInfoResponse {
    1: required User user               // 用户信息
}

// 上传头像请求
struct UploadAvatarRequest {
    1: required i64 user_id            // 用户ID
    2: required binary avatar_data     // 头像二进制数据
    3: required string content_type    // 文件类型
}

// 上传头像响应
struct UploadAvatarResponse {
    1: required string avatar_url      // 头像URL
}

// 刷新令牌请求
struct RefreshTokenRequest {
    1: required string refresh_token  // 刷新令牌
}

// 刷新令牌响应
struct RefreshTokenResponse {
    1: required string token         // 新的访问令牌
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