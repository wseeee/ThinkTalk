package code

import "ThinkTalk/pkg/xcode"

var (
	BizIdEmpty        = xcode.New(60001, "业务ID不能为空")
	TargetIdEmpty     = xcode.New(60002, "评论目标ID不能为空")
	ReplyUserIdEmpty  = xcode.New(60003, "评论用户ID不能为空")
	ContentEmpty      = xcode.New(60004, "评论内容不能为空")
	ContentTooLong    = xcode.New(60005, "评论内容过长")
	ReplyNotFound     = xcode.New(60006, "评论不存在")
	CannotDeleteReply = xcode.New(60007, "无权删除此评论")
)
