package logic

import (
	"context"
	"encoding/json"

	"ThinkTalk/application/reply/code"
	"ThinkTalk/application/reply/rpc/internal/svc"
	"ThinkTalk/application/reply/rpc/pb"

	"ThinkTalk/application/reply/rpc/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type CreateReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReplyLogic {
	return &CreateReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateReplyLogic) CreateReply(in *pb.CreateReplyRequest) (*pb.CreateReplyResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.TargetId == 0 {
		return nil, code.TargetIdEmpty
	}
	if in.ReplyUserId == 0 {
		return nil, code.ReplyUserIdEmpty
	}
	if in.Content == "" {
		return nil, code.ContentEmpty
	}
	if len(in.Content) > 5000 {
		return nil, code.ContentTooLong
	}

	msg := &types.ReplyMsg{
		BizId:         in.BizId,
		TargetId:      in.TargetId,
		ReplyUserId:   in.ReplyUserId,
		BeReplyUserId: in.BeReplyUserId,
		ParentId:      in.ParentId,
		Content:       in.Content,
		OpType:        types.OpTypeCreate,
	}

	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Errorf("[CreateReply] marshal msg: %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(l.ctx, string(data))
		if err != nil {
			l.Errorf("[CreateReply] kq push data: %s error: %v", data, err)
		}
	})

	return &pb.CreateReplyResponse{}, nil
}
