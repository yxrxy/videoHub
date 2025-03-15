struct Video {
    1: required i64 id                   // 视频ID
    2: required i64 user_id              // 作者ID
    3: required string video_url         // 视频URL
    4: required string cover_url         // 封面URL
    5: required string title             // 视频标题
    6: optional string description       // 视频描述
    7: required i64 duration             // 视频时长（秒）
    8: required string category          // 视频分类
    9: required list<string> tags        // 视频标签
    10: required i64 visit_count         // 播放量
    11: required i64 like_count          // 点赞数
    12: required i64 comment_count       // 评论数
    13: required bool is_private         // 是否私有
    14: required i64 created_at          // 创建时间
    15: required i64 updated_at          // 更新时间
    16: optional i64 deleted_at          // 删除时间
}

struct VideoListRequest {
    1: required i64 user_id              // 查询的用户ID
    2: optional i64 page                 // 第几页
    3: optional i32 size                 // 每页数量
    4: optional string category          // 按分类筛选
}

struct VideoListResponse {
    1: required list<Video> videos
    2: required i64 total
}

struct PublishRequest {
    1: required i64 user_id              // 用户ID
    2: required binary video_data        // 视频二进制数据
    3: required string content_type      // 文件类型
    4: required string title             // 视频标题
    5: optional string description       // 视频描述
    6: optional string category          // 视频分类
    7: optional list<string> tags        // 视频标签
    8: optional bool is_private          // 是否私有
}

struct PublishResponse {
    1: required string video_url         // 视频URL
    2: required string cover_url         // 封面URL
}

struct HotVideoRequest {
    1: optional i32 limit                // 获取数量限制，默认10
    2: optional string category          // 可选的分类筛选
    3: optional i64 last_visit               // 分页游标适合抖音不断滑动推荐视频
    4: optional i64 last_like
    5: optional i64 last_id
}

struct HotVideoResponse {
    1: required list<Video> videos       // 热门视频列表
    2: required i64 total               // 总数
    3: optional i64 next_visit               // 分页游标适合抖音不断滑动推荐视频
    4: optional i64 next_like
    5: optional i64 next_id
}

struct IncrementVisitCountRequest {
    1: required i64 video_id              // 视频ID
}

struct IncrementVisitCountResponse {
    1: required bool success
}

struct IncrementLikeCountRequest {
    1: required i64 video_id              // 视频ID
}

struct IncrementLikeCountResponse {
    1: required bool success
}

service VideoService {
    PublishResponse Publish(1: PublishRequest req)
    VideoListResponse GetVideoList(1: VideoListRequest req)
    HotVideoResponse GetHotVideos(1: HotVideoRequest req)
    IncrementVisitCountResponse IncrementVisitCount(1: IncrementVisitCountRequest req)   // 增加播放量
    IncrementLikeCountResponse IncrementLikeCount(1: IncrementLikeCountRequest req)     // 增加点赞数
}
