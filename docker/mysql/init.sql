-- 创建数据库
CREATE DATABASE IF NOT EXISTS videoHub;

-- 创建用户并授权
CREATE USER IF NOT EXISTS 'videoHub'@'%' IDENTIFIED BY 'videoHub';
GRANT ALL PRIVILEGES ON videoHub.* TO 'videoHub'@'%';
FLUSH PRIVILEGES;

USE videoHub;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(32) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 视频主表
CREATE TABLE IF NOT EXISTS videos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    video_url VARCHAR(255) NOT NULL,
    cover_url VARCHAR(255) NOT NULL,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    duration BIGINT NOT NULL DEFAULT 0,
    category VARCHAR(32) NOT NULL,
    visit_count BIGINT NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    comment_count BIGINT NOT NULL DEFAULT 0,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_category (category),
    INDEX idx_created_at (created_at),
    INDEX idx_hot_videos ((visit_count + like_count * 1.5) DESC, id DESC),
    INDEX idx_visit (visit_count DESC),
    INDEX idx_like (like_count DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 视频标签表
CREATE TABLE IF NOT EXISTS video_tags (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    tag VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_video_tag (video_id, tag),
    INDEX idx_tag (tag)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 视频访问记录表
CREATE TABLE IF NOT EXISTS video_visits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    ip VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_video_id (video_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 视频点赞表
CREATE TABLE IF NOT EXISTS video_likes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_video_user (video_id, user_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 视频评论表
CREATE TABLE IF NOT EXISTS video_comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content VARCHAR(512) NOT NULL,
    like_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_video_id (video_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 私信表
CREATE TABLE IF NOT EXISTS private_messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sender_id BIGINT NOT NULL,           -- 发送者ID
    receiver_id BIGINT NOT NULL,         -- 接收者ID
    content TEXT NOT NULL,               -- 消息内容
    is_read BOOLEAN DEFAULT FALSE,       -- 是否已读
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_sender_receiver (sender_id, receiver_id),
    INDEX idx_receiver_sender (receiver_id, sender_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 聊天室表
CREATE TABLE IF NOT EXISTS chat_rooms (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,           -- 聊天室名称
    creator_id BIGINT NOT NULL,          -- 创建者ID
    type TINYINT NOT NULL,              -- 类型：1=私聊,2=群聊
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_creator (creator_id),
    INDEX idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 聊天室成员表
CREATE TABLE IF NOT EXISTS chat_room_members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    room_id BIGINT NOT NULL,             -- 聊天室ID
    user_id BIGINT NOT NULL,             -- 用户ID
    nickname VARCHAR(32),                -- 在群里的昵称
    role TINYINT DEFAULT 0,              -- 角色：0=普通成员,1=管理员,2=群主
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_room_user (room_id, user_id),
    INDEX idx_user_room (user_id, room_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 聊天消息表
CREATE TABLE IF NOT EXISTS chat_messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    room_id BIGINT NOT NULL,             -- 聊天室ID
    sender_id BIGINT NOT NULL,           -- 发送者ID
    content TEXT NOT NULL,               -- 消息内容
    type TINYINT DEFAULT 0,              -- 消息类型：0=文本,1=图片,2=视频,3=文件
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_room_created (room_id, created_at),
    INDEX idx_sender (sender_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 好友关系表
CREATE TABLE IF NOT EXISTS friendships (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,             -- 用户ID
    friend_id BIGINT NOT NULL,           -- 好友ID
    status TINYINT DEFAULT 0,            -- 状态：0=待确认,1=已接受,2=已拒绝,3=已拉黑
    remark VARCHAR(32),                  -- 好友备注
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_user_friend (user_id, friend_id),
    INDEX idx_friend_user (friend_id, user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 好友申请表
CREATE TABLE IF NOT EXISTS friend_requests (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sender_id BIGINT NOT NULL,           -- 发送者ID
    receiver_id BIGINT NOT NULL,         -- 接收者ID
    message VARCHAR(255),                -- 申请消息
    status TINYINT DEFAULT 0,            -- 状态：0=待处理,1=已接受,2=已拒绝
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_sender (sender_id),
    INDEX idx_receiver (receiver_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 消息已读状态表
CREATE TABLE IF NOT EXISTS message_read_status (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    message_id BIGINT NOT NULL,          -- 消息ID
    user_id BIGINT NOT NULL,             -- 用户ID
    is_read BOOLEAN DEFAULT FALSE,       -- 是否已读
    read_at TIMESTAMP NULL,              -- 读取时间
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_message_user (message_id, user_id),
    INDEX idx_user_message (user_id, message_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
 