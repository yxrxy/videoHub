-- 点赞表
CREATE TABLE IF NOT EXISTS likes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '点赞ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    video_id BIGINT NOT NULL COMMENT '视频ID',
    deleted_at BIGINT NULL COMMENT '删除时间',
    UNIQUE KEY `idx_user_video` (`user_id`, `video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表';

-- 评论表
CREATE TABLE IF NOT EXISTS comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评论ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    video_id BIGINT NOT NULL COMMENT '视频ID',
    content TEXT NOT NULL COMMENT '评论内容',
    parent_id BIGINT NULL COMMENT '父评论ID',
    like_count INT NOT NULL DEFAULT 0 COMMENT '点赞数',
    deleted_at BIGINT NULL COMMENT '删除时间',
    KEY `idx_video` (`video_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 评论点赞表
CREATE TABLE IF NOT EXISTS comment_likes (
    user_id BIGINT NOT NULL COMMENT '用户ID',
    comment_id BIGINT NOT NULL COMMENT '评论ID',
    UNIQUE KEY `idx_user_comment` (`user_id`, `comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论点赞表';