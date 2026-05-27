package code

import "ThinkTalk/pkg/xcode"

var (
	UserIdEmpty        = xcode.New(90001, "用户ID不能为空")
	LevelInvalid       = xcode.New(90002, "会员等级无效")
	TransactionIdEmpty = xcode.New(90003, "支付流水号不能为空")
	OrderNotFound      = xcode.New(90004, "订单不存在")
	MemberNotFound     = xcode.New(90005, "会员信息不存在")
	DuplicateTransaction = xcode.New(90006, "重复的支付流水号")
)
