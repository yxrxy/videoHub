-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS videohub;

-- 创建用户并授权（允许从任何主机访问）
CREATE USER IF NOT EXISTS 'videohub'@'%' IDENTIFIED BY 'videohub';
GRANT ALL PRIVILEGES ON videohub.* TO 'videohub'@'%';
FLUSH PRIVILEGES;

-- 使用数据库
USE videohub;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` varchar(32) NOT NULL COMMENT '用户名',
    `password` varchar(128) NOT NULL COMMENT '密码',
    `avatar_url` varchar(255) NOT NULL DEFAULT '' COMMENT '头像URL',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';