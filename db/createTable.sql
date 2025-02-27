---- 微信token信息表
CREATE TABLE IF NOT EXISTS wx_token (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `token` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'access token',
    `expire_at` bigint(20) NOT NULL DEFAULT 0 COMMENT '过期时间',
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `is_locked` TINYINT(1) DEFAULT 0 COMMENT '是否锁住'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='微信access token'
---- 初始化插入默认id 1数据
insert into wx_token (id) value (1);

----- 未注册用户表
CREATE TABLE IF NOT EXISTS unRegister_user_info (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `open_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信open id',
    `union_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信union_id',
    `session_key` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信返回的session_key',
    `create_time`    timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='未注册用户表'

----- 用户信息表
CREATE TABLE IF NOT EXISTS t_user (
    `uid` bigint(20) PRIMARY KEY COMMENT 'uid',
    `nick` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '昵称',
    `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像',
    `age` SMALLINT(4) NOT NUL DEFAULT 0 COMMENT '年龄',
    `sex` VARCHAR(16) NOT NULL DEFAULT '男' COMMENT '性别',
    `open_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信open id',
    `union_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信union_id',
    `phone` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '电话',
    `session_key` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信返回的session_key',
    `create_time`    timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表'