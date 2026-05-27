package code

import "ThinkTalk/pkg/xcode"

var (
	BizIdEmpty    = xcode.New(80001, "业务ID不能为空")
	ObjIdEmpty    = xcode.New(80002, "收藏对象ID不能为空")
	UserIdEmpty   = xcode.New(80003, "用户ID不能为空")
)
