namespace go user

// 基础响应结构
struct BaseResp {
    1: i32 code                  // 状态码
    2: string message            // 响应信息
}

// 用户信息
struct User {
    1: i64 id                    // 用户ID
    2: string username           // 用户名
    3: string avatar_url         // 头像URL
    4: i64 created_at            // 创建时间
    5: i64 updated_at            // 更新时间
    6: i64 deleted_at            // 删除时间
}

// 注册请求
struct RegisterRequest {
    1: string username (vt.min_size = "1", vt.max_size = "32")    // 用户名
    2: string password (vt.min_size = "6", vt.max_size = "32")    // 密码
}

// 注册响应
struct RegisterResponse {
    1: BaseResp base             // 基础响应
    2: i64 user_id              // 用户ID
    3: string token             // 访问令牌
    4: string refresh_token     // 刷新令牌
}

// 登录请求
struct LoginRequest {
    1: string username (vt.min_size = "1")    // 用户名
    2: string password (vt.min_size = "6")    // 密码
}

// 登录响应
struct LoginResponse {
    1: BaseResp base            // 基础响应
    2: i64 user_id             // 用户ID
    3: string token            // 访问令牌
    4: string refresh_token    // 刷新令牌
}

// 获取用户信息请求
struct UserInfoRequest {
    1: i64 user_id             // 用户ID
}

// 获取用户信息响应
struct UserInfoResponse {
    1: BaseResp base           // 基础响应
    2: User user               // 用户信息
}

// 上传头像请求
struct UploadAvatarRequest {
    1: binary avatar_data      // 头像二进制数据
    2: string content_type     // 文件类型，如 image/jpeg
}

// 上传头像响应
struct UploadAvatarResponse {
    1: BaseResp base           // 基础响应
    2: string avatar_url       // 头像URL
}

// MFA二维码请求
struct MFAQRCodeRequest {
    1: i64 user_id            // 用户ID
}

// MFA二维码响应
struct MFAQRCodeResponse {
    1: BaseResp base          // 基础响应
    2: string qr_code         // 二维码数据
}

// MFA绑定请求
struct BindMFARequest {
    1: i64 user_id           // 用户ID
    2: string code           // MFA验证码
}

// MFA绑定响应
struct BindMFAResponse {
    1: BaseResp base         // 基础响应
}

// 刷新令牌请求
struct RefreshTokenRequest {
    1: string refresh_token     // 刷新令牌
}

// 刷新令牌响应
struct RefreshTokenResponse {
    1: BaseResp base           // 基础响应
    2: string token           // 新的访问令牌
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
    // 获取MFA二维码
    MFAQRCodeResponse GetMFAQRCode(1: MFAQRCodeRequest req)
    // 绑定MFA
    BindMFAResponse BindMFA(1: BindMFARequest req)
    // 刷新令牌
    RefreshTokenResponse RefreshToken(1: RefreshTokenRequest req)
}