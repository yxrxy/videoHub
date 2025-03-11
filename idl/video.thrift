struct BaseResp {
    1: i32 code                  // 状态码
    2: string message            // 响应信息
}

struct PublishRequest {
    1: i64 user_id,                // 用户ID
    2: binary video_data,          // 视频二进制数据
    3: string content_type,        // 文件类型
    4: string title,               // 视频标题
    5: string description,         // 视频描述
}

struct PublishResponse {
    1: BaseResp base,
    2: string video_url,           // 视频URL
}

service VideoService {
    PublishResponse Publish(1: PublishRequest req),
} 