create database thinktalk_chat;
use thinktalk_chat;

CREATE TABLE `conversation` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '用户ID',
    `target_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '对方用户ID',
    `last_message` varchar(512) NOT NULL DEFAULT '' COMMENT '最后一条消息',
    `last_message_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后消息时间',
    `unread_count` int(11) NOT NULL DEFAULT '0' COMMENT '未读数',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_target` (`user_id`, `target_user_id`),
    KEY `ix_user_id` (`user_id`),
    KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='会话表';

CREATE TABLE `message` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `conversation_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '会话ID',
    `sender_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '发送者ID',
    `receiver_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '接收者ID',
    `content` varchar(2048) NOT NULL DEFAULT '' COMMENT '消息内容',
    `msg_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '消息类型 0:文本 1:图片',
    `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读',
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:正常 1:撤回',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `ix_conversation_id` (`conversation_id`),
    KEY `ix_sender_receiver` (`sender_id`, `receiver_id`),
    KEY `ix_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='消息表';
