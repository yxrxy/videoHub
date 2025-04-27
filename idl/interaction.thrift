namespace go interaction

include "model.thrift"

// 评论列表响应
struct CommentListResponse {
    1: required model.BaseResp Base           // 基本响应信息
    2: required list<model.Comment> comments  // 评论列表
    3: required i64 total                     // 总数
}

// 点赞请求
struct LikeRequest {
    1: required i64 video_id         // 视频ID
}

// 点赞响应
struct LikeResponse {
    1: required model.BaseResp Base  // 基本响应信息
}

// 获取点赞列表请求
struct GetLikesRequest {
    1: required i64 video_id         // 视频ID
    2: required i32 page             // 页码
    3: required i32 size             // 每页大小
}

// 获取点赞列表响应
struct GetLikesResponse {
    1: required model.BaseResp Base              // 基本响应信息
    2: required list<model.LikeInfo> LikeList    // 点赞列表
    3: required i64 Total                        // 总数
}

// 发表评论请求
struct CommentRequest {
    1: required i64 video_id         // 视频ID
    2: required string content       // 评论内容
    3: optional i64 parent_id        // 父评论ID
}

// 发表评论响应
struct CommentResponse {
    1: required model.BaseResp Base  // 基本响应信息
    2: required i64 comment_id       // 评论ID
}

// 获取评论列表请求
struct GetCommentsRequest {
    1: required i64 video_id         // 视频ID
    2: required i32 page             // 页码
    3: required i32 size             // 每页大小
}

// 获取评论列表响应
struct GetCommentsResponse {
    1: required model.BaseResp Base              // 基本响应信息
    2: required list<model.Comment> CommentList  // 评论列表
    3: required i64 Total                        // 总数
}

// 删除评论请求
struct DeleteCommentRequest {
    1: required i64 comment_id       // 评论ID
}

// 删除评论响应
struct DeleteCommentResponse {
    1: required model.BaseResp Base  // 基本响应信息
}

// 点赞评论请求
struct LikeCommentRequest {
    1: required i64 comment_id       // 评论ID
}

// 点赞评论响应
struct LikeCommentResponse {
    1: required model.BaseResp Base  // 基本响应信息
}

// 互动服务
service InteractionService {
    // 点赞操作
    LikeResponse Like(1: LikeRequest req)
    
    // 获取点赞列表
    GetLikesResponse GetLikes(1: GetLikesRequest req)
    
    // 发表评论
    CommentResponse Comment(1: CommentRequest req)
    
    // 获取评论列表
    GetCommentsResponse GetComments(1: GetCommentsRequest req)
    
    // 删除评论
    DeleteCommentResponse DeleteComment(1: DeleteCommentRequest req)

    // 点赞评论
    LikeCommentResponse LikeComment(1: LikeCommentRequest req)
}