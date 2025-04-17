CREATE TABLE IF NOT EXISTS operation_records (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `action` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '操作',
    `target` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '目标',
    `details` TEXT COMMENT '详情',
    `username` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '用户名',
    `user_id` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '用户ID',
    `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '是否删除',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作记录表';