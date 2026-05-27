package code

import "ThinkTalk/pkg/xcode"

var (
	SenderIdEmpty   = xcode.New(10001, "发送者ID不能为空")
	ReceiverIdEmpty = xcode.New(10002, "接收者ID不能为空")
	ContentEmpty    = xcode.New(10003, "消息内容不能为空")
	UserIdEmpty     = xcode.New(10004, "用户ID不能为空")
	CannotSelfChat  = xcode.New(10005, "不能给自己发消息")
)
