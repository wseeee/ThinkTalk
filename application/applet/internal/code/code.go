package code

import "ThinkTalk/pkg/xcode"

var (
	RegisterNameEmpty     = xcode.New(100004, "名字不能为空")
	RegisterMobileEmpty   = xcode.New(10001, "注册手机号不能为空")
	VerificationCodeEmpty = xcode.New(100002, "验证码不能为空")
	MobileHasRegistered   = xcode.New(100003, "手机号已经注册")
	LoginMobileEmpty      = xcode.New(100003, "手机号不能为空")
	RegisterPasswordEmpty = xcode.New(100004, "密码不能为空")
)
