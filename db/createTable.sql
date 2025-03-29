---- 微信token信息表
CREATE TABLE IF NOT EXISTS wx_token (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `token` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'access token',
    `expire_at` bigint(20) NOT NULL DEFAULT 0 COMMENT '过期时间',
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `is_locked` TINYINT(1) DEFAULT 0 COMMENT '是否锁住'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='微信access token';
---- 初始化插入默认id 1数据
insert into wx_token (id) value (1);

----- 未注册用户表
CREATE TABLE IF NOT EXISTS unRegister_user_info (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `open_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信open id',
    `union_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信union_id',
    `session_key` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信返回的session_key',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='未注册用户表';

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
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';

----- 聊天统计表
CREATE TABLE IF NOT EXISTS t_chat_statistic (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `digital_uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '机器人对应的uid',
    `chat_uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '发起聊天对应的uid',
    `is_anonymity`  tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是匿名 -- 0不是匿名,1匿名',
    `chat_num` bigint(20)  NOT NULL DEFAULT 0 COMMENT '聊天次数',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uid` (`digital_uid`,`chat_uid`,`is_anonymity`),
    KEY `idx_num` (`digital_uid`,`chat_num`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天统计表';

----- 关注表
CREATE TABLE IF NOT EXISTS t_follow (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '被follow的uid',
    `follower` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注者',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `uid` (`uid`,`follower`),
    KEY `idx_follower` (`follower`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='关注表';

----- 点赞表
CREATE TABLE IF NOT EXISTS t_favor (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '被点赞的uid',
    `liker` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注者',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `uid` (`uid`,`liker`),
    KEY `idx_liker` (`liker`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';

----- 浏览记录表
CREATE TABLE IF NOT EXISTS t_view_record (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '被浏览的uid',
    `viewer` bigint(20)  NOT NULL DEFAULT 0 COMMENT '浏览者',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_uid` (`uid`) USING BTREE,
    KEY `idx_viewer` (`viewer`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='浏览记录表';

----- 社交信息统计表 -- 关注、点赞、浏览统计
CREATE TABLE IF NOT EXISTS t_social_statistic (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20)  NOT NULL DEFAULT 0 COMMENT '的uid',
    `follow_num` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注次数',
    `favor_num` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注次数',
    `view_num` bigint(20)  NOT NULL DEFAULT 0 COMMENT '关注次数',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社交信息统计表--关注、点赞、浏览统计';

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
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数字人信息';

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
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`,`digital_uid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='当前会话配置信息';

------ 会话列表
CREATE TABLE IF NOT EXISTS t_conversation (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天人uid',
    `digital_uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '机器人uid',
    `is_anonymity` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否匿名,0不匿名,1匿名',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`,`digital_uid`,`is_anonymity`),
    KEY `idx_digital_uid` (`digital_uid`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话列表';

------ 消息记录表
CREATE TABLE IF NOT EXISTS t_message (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天人uid',
    `digital_uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '机器人uid',
    `conversation_id`  bigint(20) NOT NULL DEFAULT 0 COMMENT '所属会话的ID',
    `u_message` TEXT NOT NULL COMMENT '用户说的',
    `d_message` TEXT NOT NULL COMMENT '机器人说的',
    `chat_conf` text NOT NULL COMMENT '聊天配置',
    `voice_url` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '语音条',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_uid` (`uid`,`digital_uid`) USING BTREE,
    KEY `idx_cid` (`conversation_id`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息记录表';

------ 异步结果表
CREATE TABLE IF NOT EXISTS t_async_result (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT '聊天人uid',
    `result` TEXT NOT NULL COMMENT '结果',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='异步结果表';


------ 主页配置表
CREATE TABLE IF NOT EXISTS t_homepage_conf (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT 'uid',
    `conf_str` TEXT NOT NULL COMMENT '配置信息',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid` (`uid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主页配置表';

------ 游戏配置表
CREATE TABLE IF NOT EXISTS t_game (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT 'uid',
    `name` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '游戏名',
    `prologue` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '开场白',
    `re_total_amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '红包总金额 -- 单位分',
    `re_total_num` int(11) NOT NULL DEFAULT 0 COMMENT '红包总数量',
    `re_claim_num` int(11) NOT NULL DEFAULT 0 COMMENT '红包被领取的数量',
    `answer_list_str` TEXT NOT NULL COMMENT '答案信息',
    `state` int(11) NOT NULL DEFAULT 0 COMMENT '游戏状态',
    `version` int(11) NOT NULL DEFAULT 0 COMMENT '版本号',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_uid` (`uid`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='游戏配置表';

------ 支付中心表
CREATE TABLE IF NOT EXISTS t_pay_center (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT 'uid',
    `open_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'open id',
    `trade_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '业务订单号',
    `ch_trade_no`     varchar(64)  NOT NULL DEFAULT '' COMMENT '渠道订单id',
    `amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '支付金额',
    `payer_total`     bigint(20) NOT NULL DEFAULT '0' COMMENT '用户支付的总金额,分',
    `prod_type`       int(10) NOT NULL DEFAULT '1' COMMENT '商品类型',
    `prod_desc`       varchar(256) NOT NULL DEFAULT '' COMMENT '订单描述',
    `prod_attach`     varchar(256) NOT NULL DEFAULT '' COMMENT '订单Attach',
    `order_ctime`     varchar(19)  NOT NULL DEFAULT '' COMMENT '支付请求时间',
    `order_life_time` int(10) NOT NULL DEFAULT '0' COMMENT '订单有效时间, 秒',
    `order_etime`     varchar(19)  NOT NULL DEFAULT '' COMMENT '支付完成时间',
    `trade_state`     int(10)  NOT NULL DEFAULT '0' COMMENT '交易状态 PayState',
    `expand_str`      text         NOT NULL COMMENT '具体业务产生订单时附带的信息',
    `status_msg`      varchar(128) NOT NULL DEFAULT '' COMMENT '状态描述',
    `create_time`     datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`     datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_uid` (`uid`),
    KEY `idx_trade_no` (`trade_no`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付中心表';

CREATE TABLE `sys_config` (
     `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键ID',
     `config_key` char(128) NOT NULL COMMENT '配置KEY唯一',
     `service_name` varchar(30) NOT NULL DEFAULT 'mkg' COMMENT '服务名',
     `config_value` text COMMENT '配置值',
     `description` varchar(1024) DEFAULT NULL COMMENT '说明',
     `operator` varchar(50) DEFAULT NULL COMMENT '操作者',
     `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
     PRIMARY KEY (`id`),
     UNIQUE KEY `uk_sc_configKey` (`config_key`)
) ENGINE=InnoDB AUTO_INCREMENT=185 DEFAULT CHARSET=utf8mb4 COMMENT='系统配置';

------ 游戏轮次记录表
CREATE TABLE IF NOT EXISTS t_game_play_record (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT 'uid',
    `game_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '游戏id',
    `input` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '输入 -- 用户的回答',
    `output` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '输出结果',
    `result_state` int(11) NOT NULL DEFAULT 0 COMMENT '游戏结果1:错误,2:正确',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_uid_game_id` (`uid`,`game_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='游戏轮次记录表';

------ 资金池数据
CREATE TABLE IF NOT EXISTS t_funding_pool (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT 'uid',
    `pool_type` int(11) NOT NULL DEFAULT 1 COMMENT '资金池类型:1:数字人,2:红包',
    `all_total_amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '一直以来的总收入 -- 单位分',
    `all_withdraw_amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '一直以来被提现的金额 -- 单位分',
    `cur_total_amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '当前资金 -- 单位分',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid_type` (`uid`,`pool_type`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资金池数据';

------ 资金池记详情
CREATE TABLE IF NOT EXISTS t_funding_pool_detail (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `uid` bigint(20) NOT NULL DEFAULT 0 COMMENT 'uid',
    `pool_type` int(11) NOT NULL DEFAULT 1 COMMENT '资金池类型:1:数字人,2:红包',
    `trade_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '业务订单号--虚拟的',
    `trade_type` int(11) NOT NULL DEFAULT 1 COMMENT '交易类型,1:进账,2:出账',
    `real_trade_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '业务订单号--真实的,用户看不到',
    `amount` bigint(20) NOT NULL DEFAULT 0 COMMENT '金额 -- 单位分',
    `create_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uq_uid_type_trade_no` (`uid`,`pool_type`,`trade_no`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资金池记详情';
