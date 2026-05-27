package types

const (
	ReplyStatusNormal  = 0 // 正常
	ReplyStatusDeleted = 1 // 删除

	SortByTime  = 0 // 按时间倒序
	SortByLike  = 1 // 按点赞数倒序

	DefaultPageSize      = 20
	MaxSubReplyCount     = 3  // 每条根评论最多展示的子回复数
	CacheMaxReplyCount   = 500
)
