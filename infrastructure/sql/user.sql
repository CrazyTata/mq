CREATE TABLE IF NOT EXISTS `user` (
    `id` INT UNSIGNED AUTO_INCREMENT COMMENT '用户ID',  
    `app_id`         VARCHAR(50) NOT NULL DEFAULT '' COMMENT '归属应用',
    `user_id` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '用户ID', 
    `phone` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '手机号（支持手机号一键登录）',
    `apple_id` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Apple 用户唯一标识符(sub 字段)',
    `email` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '用户邮箱(Apple登录可能提供)',
    `is_private_email` TINYINT NOT NULL DEFAULT 0 COMMENT '是否是Apple隐私邮箱',
    `name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '用户昵称',
    `avatar` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '用户头像URL',
    `password`   VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码',
    `register_time`   DATETIME,
    `status`          TINYINT NOT NULL DEFAULT 1 COMMENT '用户状态:1 正常 2冻结 3 注销',
    `source`          TINYINT NOT NULL DEFAULT 1 COMMENT '用户来源:4 安卓 5 IOS',
    `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '是否删除',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_apple_id` (`apple_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

