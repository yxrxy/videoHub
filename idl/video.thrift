namespace go video

include "model.thrift"

struct VideoListRequest {
    1: required i64 user_id              // 查询的用户ID
    2: required i64 page                 // 第几页
    3: required i32 size                 // 每页数量
    4: optional string category          // 按分类筛选
}

// 视频列表响应
struct VideoListResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required list<model.Video> VideoList  // 视频列表
    3: required i64 total                // 总数
}

// 视频详情请求
struct DetailRequest {
    1: required i64 video_id             // 视频ID
}

// 视频详情响应
struct DetailResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required model.Video Video        // 视频信息
}

// 发布视频请求
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

// 发布视频响应
struct PublishResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required string video_url         // 视频URL
    3: required string cover_url         // 封面URL
}

// 热门视频请求
struct HotVideoRequest {
    1: optional i32 limit                // 获取数量限制，默认10
    2: optional string category          // 可选的分类筛选
    3: optional i64 last_visit         
    4: optional i64 last_like
    5: optional i64 last_id
}

// 热门视频响应
struct HotVideoResponse {
    1: required model.BaseResp Base      // 基本响应信息
    2: required list<model.Video> videos // 热门视频列表
    3: required i64 total                // 总数
    4: optional i64 next_visit
    5: optional i64 next_like
    6: optional i64 next_id
}

// 删除视频请求
struct DeleteRequest {
    1: required i64 video_id             // 视频ID
}

// 删除视频响应
struct DeleteResponse {
    1: required model.BaseResp Base      // 基本响应信息
}

// 增加访问量请求
struct IncrementVisitCountRequest {
    1: required i64 video_id              // 视频ID
}

// 增加访问量响应
struct IncrementVisitCountResponse {
    1: required model.BaseResp Base       // 基本响应信息
}

// 增加点赞数请求
struct IncrementLikeCountRequest {
    1: required i64 video_id              // 视频ID
}

// 增加点赞数响应
struct IncrementLikeCountResponse {
    1: required model.BaseResp Base       // 基本响应信息
}

// 搜索视频请求
struct SearchRequest {
    1: required string keywords          // 搜索关键词
    2: required i32 page_size           // 每页数量
    3: required i32 page_num            // 页码
    4: optional i64 from_date           // 开始时间
    5: optional i64 to_date             // 结束时间
    6: optional string username         // 按用户名筛选
}

// 搜索视频响应
struct SearchResponse {
    1: required model.BaseResp Base     // 基本响应信息
    2: required list<model.Video> videos // 视频列表
    3: required i64 total               // 总数
}

// 语义搜索视频请求
struct SemanticSearchRequest {
    1: required string query            // 搜索查询文本
    2: required i32 page_size           // 每页数量
    3: required i32 page_num            // 页码
    4: optional double threshold        // 相似度阈值(0-1)
}

// 语义搜索视频响应
struct SemanticSearchResponse {
    1: required model.BaseResp Base     // 基本响应信息
    2: required list<SemanticSearchResultItem> results // 语义搜索结果
    3: required i64 total               // 总数
    4: optional list<string> related_queries // 相关查询建议
    5: optional string summary          // 搜索结果摘要
}

service VideoService {
    PublishResponse Publish(1: PublishRequest req)
    VideoListResponse List(1: VideoListRequest req)
    DetailResponse Detail(1: DetailRequest req)
    HotVideoResponse GetHotVideos(1: HotVideoRequest req)
    DeleteResponse Delete(1: DeleteRequest req)
    IncrementVisitCountResponse IncrementVisitCount(1: IncrementVisitCountRequest req)
    IncrementLikeCountResponse IncrementLikeCount(1: IncrementLikeCountRequest req)
    SearchResponse Search(1: SearchRequest req)
    SemanticSearchResponse SemanticSearch(1: SemanticSearchRequest req)
}
