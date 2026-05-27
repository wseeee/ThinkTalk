create database thinktalk_qa;
use thinktalk_qa;

CREATE TABLE `question` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `title` varchar(256) NOT NULL DEFAULT '' COMMENT '问题标题',
    `content` text NOT NULL COMMENT '问题描述',
    `author_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '提问者ID',
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:正常 1:删除',
    `answer_num` int(11) NOT NULL DEFAULT '0' COMMENT '回答数',
    `view_num` int(11) NOT NULL DEFAULT '0' COMMENT '浏览数',
    `tag_ids` varchar(512) NOT NULL DEFAULT '' COMMENT '标签ID列表',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    KEY `ix_author_id` (`author_id`),
    KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='问题表';

CREATE TABLE `answer` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `question_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '问题ID',
    `author_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '回答者ID',
    `content` text NOT NULL COMMENT '回答内容',
    `is_accepted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否采纳 0:否 1:是',
    `like_num` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
    `reply_num` int(11) NOT NULL DEFAULT '0' COMMENT '回复数',
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:正常 1:删除',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    KEY `ix_question_id` (`question_id`),
    KEY `ix_author_id` (`author_id`),
    KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='回答表';
