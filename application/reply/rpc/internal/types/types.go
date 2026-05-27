package types

type ReplyMsg struct {
	BizId         string `json:"bizId"`
	TargetId      int64  `json:"targetId"`
	ReplyUserId   int64  `json:"replyUserId"`
	BeReplyUserId int64  `json:"beReplyUserId"`
	ParentId      int64  `json:"parentId"`
	Content       string `json:"content"`
	OpType        int32  `json:"opType"`
	ReplyId       int64  `json:"replyId"`
	UserId        int64  `json:"userId"`
}

const (
	OpTypeCreate = 0
	OpTypeDelete = 1
)
