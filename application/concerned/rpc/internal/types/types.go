package types

type ConcernedMsg struct {
	BizId  string `json:"bizId"`  // 业务ID
	ObjId  int64  `json:"objId"`  // 收藏对象id
	UserId int64  `json:"userId"` // 用户ID
	OpType int32  `json:"opType"` // 操作类型 0:收藏 1:取消
}

const (
	OpTypeAdd    = 0
	OpTypeCancel = 1
)
