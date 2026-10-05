CREATE TABLE `member` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `level` tinyint NOT NULL DEFAULT 0 COMMENT '会员等级 0=普通 1=黄金 2=钻石',
  `expire_time` datetime NOT NULL COMMENT '过期时间',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态 0=已过期 1=有效',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  KEY `ix_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='会员表';

CREATE TABLE `member_order` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `level` tinyint NOT NULL COMMENT '会员等级',
  `duration_days` int NOT NULL COMMENT '开通天数',
  `amount` bigint NOT NULL COMMENT '金额(分)',
  `pay_channel` varchar(32) NOT NULL COMMENT '支付渠道 wxpay/alipay',
  `transaction_id` varchar(128) NOT NULL COMMENT '支付流水号',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态 0=待支付 1=已支付 2=已退款',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_id` (`transaction_id`),
  KEY `ix_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='会员订单表';
