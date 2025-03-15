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
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='未注册用户表'

----- 用户信息表
CREATE TABLE IF NOT EXISTS t_user (
    `uid` bigint(20) PRIMARY KEY COMMENT 'uid',
    `nick` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '昵称',
    `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像',
    `age` SMALLINT(4) NOT NULL DEFAULT 0 COMMENT '年龄',
    `sex` VARCHAR(16) NOT NULL DEFAULT '男' COMMENT '性别',
    `education` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '学历',
    `mbti_str` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT 'MBTI 信息',
    `tag_info_str`text NOT NULL COMMENT '标签信息',
    `expand`text NOT NULL COMMENT '扩展信息',
    `open_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信open id',
    `union_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信union_id',
    `phone` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '电话',
    `session_key` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信返回的session_key',
    `is_visible` tinyint(1) NOT NULL DEFAULT 1 COMMENT '他人是否可见 -- 0表示不可见',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表'

----- 聊天统计表
CREATE TABLE IF NOT EXISTS t_chat_statistic (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `digital_uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '机器人对应的uid',
    `chat_uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '发起聊天对应的uid',
    `is_anonymity`  tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是匿名 -- 0不是匿名,1匿名',
    `chat_num` bigint(20)  NOT NULL DEFAULT 0 COMMENT '聊天次数',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uid` (`digital_uid`,`chat_uid`,`is_anonymity`),
    KEY `idx_num` (`digital_uid`,`chat_num`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天统计表'

----- 关注表
CREATE TABLE IF NOT EXISTS t_follow (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '被follow的uid',
    `follower` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注者',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `uid` (`uid`,`follower`),
    KEY `idx_follower` (`follower`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='关注表'

----- 点赞表
CREATE TABLE IF NOT EXISTS t_favor (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '被点赞的uid',
    `liker` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注者',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `uid` (`uid`,`liker`),
    KEY `idx_liker` (`liker`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表'

----- 浏览记录表
CREATE TABLE IF NOT EXISTS t_view_record (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '被浏览的uid',
    `viewer` bigint(20)  NOT NULL DEFAULT 0 COMMENT '浏览者',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_uid` (`uid`) USING BTREE,
    KEY `idx_viewer` (`viewer`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='浏览记录表'

----- 社交信息统计表 -- 关注、点赞、浏览统计
CREATE TABLE IF NOT EXISTS t_social_statistic (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '的uid',
    `follow_num` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注次数',
    `favor_num` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注次数',
    `view_num` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注次数',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社交信息统计表--关注、点赞、浏览统计'

------ 数字人信息
CREATE TABLE IF NOT EXISTS t_digital_info (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '数字人uid',
    `icon` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像',
    `digital_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '分身名',
    `can_anonymity` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否可以匿名',
    `prologue` VARCHAR(4096) NOT NULL DEFAULT '' COMMENT '开场白',
    `clone_voice` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '克隆语音条',
    `charge_conf` text NOT NULL COMMENT '收费设置',
    `digital_all_kinds` text NOT NULL COMMENT '百变分身设置',
    `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0 未创建，1 已创建',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数字人信息'

------ 当前会话配置信息
CREATE TABLE IF NOT EXISTS t_conversation_conf (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天人uid',
    `digital_uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '机器人uid',
    `is_anonymity` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否匿名,0不匿名,1匿名',
    `chat_conf` text NOT NULL COMMENT '聊天配置',
    `last_msg` text NOT NULL COMMENT '最后一条聊天记录',
    `chat_total_num` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天总次数',
    `chat_use_num` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天已用次数',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`,`digital_uid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='当前会话配置信息'

------ 会话列表
CREATE TABLE IF NOT EXISTS t_conversation (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天人uid',
    `digital_uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '机器人uid',
    `is_anonymity` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否匿名,0不匿名,1匿名',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`,`digital_uid`,`is_anonymity`),
    KEY `idx_digital_uid` (`digital_uid`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话列表'

------ 消息记录表
CREATE TABLE IF NOT EXISTS t_message (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天人uid',
    `digital_uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '机器人uid',
    `conversation_id`  bigint(20) NOT NULL DEFAULT 0 COMMENT '所属会话的ID',
    `u_message` TEXT NOT NULL COMMENT '用户说的',
    `d_message` TEXT NOT NULL COMMENT '机器人说的',
    `chat_conf` text NOT NULL COMMENT '聊天配置',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_uid` (`uid`,`digital_uid`) USING BTREE,
    KEY `idx_cid` (`conversation_id`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息记录表'

------ 异步结果表
CREATE TABLE IF NOT EXISTS t_async_result (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天人uid',
    `result` TEXT NOT NULL COMMENT '结果',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='异步结果表'