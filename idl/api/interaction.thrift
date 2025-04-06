namespace go api.interaction

include "../interaction.thrift"

// API 服务
service InteractionAPI {
    // 点赞相关接口
    interaction.LikeResponse Like(1: interaction.LikeRequest req) (api.post="/api/v1/comment/like")
    interaction.GetLikesResponse GetLikes(1: interaction.GetLikesRequest req) (api.get="/api/v1/video/likelist")
    
    // 评论相关接口
    interaction.CommentResponse Comment(1: interaction.CommentRequest req) (api.post="/api/v1/video/comment")
    interaction.CommentListResponse GetComments(1: interaction.GetCommentsRequest req) (api.get="/api/v1/video/comments")
    interaction.DeleteCommentResponse DeleteComment(1: interaction.DeleteCommentRequest req) (api.delete="/api/v1/comment/:comment_id")
    interaction.LikeCommentResponse LikeComment(1: interaction.LikeCommentRequest req) (api.post="/api/v1/comment/:comment_id/like")
} 