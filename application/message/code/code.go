package code

import "ThinkTalk/pkg/xcode"

var (
	UserIdEmpty       = xcode.New(70001, "用户ID不能为空")
	NotificationIdEmpty = xcode.New(70002, "通知ID不能为空")
	NotificationNotFound = xcode.New(70003, "通知不存在")
)
