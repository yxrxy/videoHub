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

-- 视频表
CREATE TABLE IF NOT EXISTS video (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '视频ID',
    user_id BIGINT NOT NULL COMMENT '作者ID',
    video_url VARCHAR(255) NOT NULL COMMENT '视频URL',
    cover_url VARCHAR(255) NOT NULL COMMENT '封面URL',
    title VARCHAR(128) NOT NULL COMMENT '视频标题',
    description VARCHAR(512) DEFAULT '' COMMENT '视频描述',
    duration BIGINT NOT NULL DEFAULT 0 COMMENT '视频时长（秒）',
    category VARCHAR(32) DEFAULT '其他' COMMENT '视频分类',
    tags VARCHAR(255) DEFAULT '' COMMENT '视频标签，以逗号分隔',
    visit_count BIGINT NOT NULL DEFAULT 0 COMMENT '播放量',
    like_count BIGINT NOT NULL DEFAULT 0 COMMENT '点赞数',
    comment_count BIGINT NOT NULL DEFAULT 0 COMMENT '评论数',
    is_private BOOLEAN NOT NULL DEFAULT false COMMENT '是否私有',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_category (category),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频表';