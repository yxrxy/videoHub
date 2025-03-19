namespace go videoInteractions

// 点赞信息
struct LikeInfo {
    1: required i64 id         // 点赞ID
    2: required i64 user_id    // 用户ID
    3: required i64 video_id   // 视频ID
    4: required i64 created_at // 创建时间
    5: optional i64 deleted_at // 删除时间
}

// 评论信息
struct CommentInfo {
    1: required i64 id         // 评论ID
    2: required i64 user_id    // 用户ID
    3: required i64 video_id   // 视频ID
    4: required string content // 评论内容
    5: optional i64 parent_id  // 父评论ID
    6: required i64 created_at // 创建时间
    7: optional i64 deleted_at // 删除时间
    8: optional i32 like_count // 点赞数
    9: optional bool is_liked  // 当前用户是否点赞
}

// 评论列表响应
struct CommentListResponse {
    1: required list<CommentInfo> comments // 评论列表
    2: required i64 total                  // 总数
}

// 互动服务
service InteractionService {
    // 点赞操作
    bool Like(1: i64 user_id, 2: i64 video_id)
    
    // 获取点赞列表
    list<LikeInfo> GetLikes(1: i64 video_id, 2: i32 page, 3: i32 size)
    
    // 发表评论
    bool Comment(1: i64 user_id, 2: i64 video_id, 3: string content, 4: optional i64 parent_id)
    
    // 获取评论列表
    CommentListResponse GetComments(1: i64 video_id, 2: i32 page, 3: i32 size)
    
    // 删除评论
    bool DeleteComment(1: i64 user_id, 2: i64 comment_id)

    // 点赞评论
    bool LikeComment(1: i64 user_id, 2: i64 comment_id)
}