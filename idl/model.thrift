namespace go model

// 基本响应结构
struct BaseResp {
    1: i64 code,    // 错误码，0表示成功
    2: string msg,  // 错误信息
}

// 用户模型
struct User {
    1: required i64 id,              // 用户ID
    2: required string username,     // 用户名
    4: optional string avatar,       // 头像URL
    6: optional i64 followCount,     // 关注数
    7: optional i64 followerCount,   // 粉丝数
    8: optional bool isFollow,       // 是否已关注
    9: optional i64 likeCount,       // 获赞数量
    10: optional i64 videoCount,     // 视频数量
}

// 视频模型
struct Video {
    1: required i64 id,              // 视频ID
    2: required i64 authorId,        // 作者ID
    3: required string title,        // 标题
    4: required string playUrl,      // 播放地址
    5: required string coverUrl,     // 封面地址
    6: optional i64 favoriteCount,   // 点赞数
    7: optional i64 commentCount,    // 评论数
    8: optional bool isFavorite,     // 是否已点赞
    9: optional User author,         // 作者信息
    10: optional string description, // 视频描述
    11: optional i64 createdAt,      // 创建时间
    12: optional i64 updatedAt,      // 更新时间
}

// 评论模型
struct Comment {
    1: required i64 id,              // 评论ID
    2: required i64 userId,          // 用户ID
    3: required i64 videoId,         // 视频ID
    4: required string content,      // 评论内容
    5: optional User user,           // 用户信息
    6: optional i64 createdAt,       // 创建时间
    7: optional i64 updatedAt,       // 更新时间
}

struct LikeInfo {
    1: required i64 id
    2: required i64 user_id
    3: required i64 video_id
    4: required i64 created_at
    5: optional i64 deleted_at
}

struct CommentInfo {
    1: required i64 id
    2: required i64 user_id
    3: required i64 video_id
    4: required string content
    5: optional i64 parent_id
    6: required i64 created_at
    7: optional i64 deleted_at
    8: optional i32 like_count
    9: optional bool is_liked
}

// 私信消息
struct PrivateMessage {
    1: required i64 id                   // 消息ID
    2: required i64 sender_id            // 发送者ID
    3: required i64 receiver_id          // 接收者ID
    4: required string content           // 消息内容
    5: required bool is_read             // 是否已读
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 聊天室
struct ChatRoom {
    1: required i64 id                   // 聊天室ID
    2: required string name              // 聊天室名称
    3: required i64 creator_id           // 创建者ID
    4: required i8 type                  // 类型：1=私聊,2=群聊
    8: optional list<ChatRoomMember> members // 成员列表
}

// 聊天室成员
struct ChatRoomMember {
    1: required i64 id                   // 成员ID
    2: required i64 room_id              // 聊天室ID
    3: required i64 user_id              // 用户ID
    4: optional string nickname          // 在群里的昵称
    5: required i8 role                  // 角色：0=普通成员,1=管理员,2=群主
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 聊天消息
struct ChatMessage {
    1: required i64 id                   // 消息ID
    2: required i64 room_id              // 聊天室ID
    3: required i64 sender_id            // 发送者ID
    4: required string content           // 消息内容
    5: required i8 type                  // 消息类型：0=文本,1=图片,2=视频,3=文件
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 好友关系
struct Friendship {
    1: required i64 id                   // 关系ID
    2: required i64 user_id              // 用户ID
    3: required i64 friend_id            // 好友ID
    4: required i8 status                // 状态：0=待确认,1=已接受,2=已拒绝,3=已拉黑
    5: optional string remark            // 好友备注
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
}

// 好友申请
struct FriendRequest {
    1: required i64 id                   // 申请ID
    2: required i64 sender_id            // 发送者ID
    3: required i64 receiver_id          // 接收者ID
    4: optional string message           // 申请消息
    5: required i8 status                // 状态：0=待处理,1=已接受,2=已拒绝
    6: required i64 created_at           // 创建时间
    7: required i64 updated_at           // 更新时间
    8: optional i64 deleted_at           // 删除时间
} 

// 语义搜索结果项
struct SemanticSearchResultItem {
    1: required i64 video_id            // 视频ID
    2: required string title            // 标题
    3: optional string description      // 描述
    4: required double score            // 相似度分数
    5: required string cover_url        // 封面URL
    6: optional i32 duration_seconds    // 时长(秒)
    7: required i64 publish_time        // 发布时间
    8: optional string category         // 分类
    9: optional list<string> tags       // 标签
    10: required i64 user_id            // 用户ID
    11: required string username        // 用户名
    12: optional i64 visit_count        // 访问次数
    13: optional i64 like_count         // 点赞数
}