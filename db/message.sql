create database thinktalk_message;
use thinktalk_message;

CREATE TABLE `notification` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '接收用户ID',
    `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '通知类型 1:点赞 2:关注 3:评论 4:系统',
    `title` varchar(128) NOT NULL DEFAULT '' COMMENT '通知标题',
    `content` varchar(512) NOT NULL DEFAULT '' COMMENT '通知内容',
    `ref_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '关联资源ID',
    `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
    `trigger_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '触发用户ID',
    `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读 0:未读 1:已读',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    KEY `ix_user_id` (`user_id`),
    KEY `ix_user_read` (`user_id`, `is_read`),
    KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='通知表';
