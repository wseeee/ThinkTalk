package logic

import (
	"context"

	"ThinkTalk/application/message/code"
	"ThinkTalk/application/message/rpc/internal/svc"
	"ThinkTalk/application/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkAllReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllReadLogic {
	return &MarkAllReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkAllReadLogic) MarkAllRead(in *pb.MarkAllReadRequest) (*pb.MarkAllReadResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	err := l.svcCtx.NotificationModel.UpdateAllRead(l.ctx, in.UserId, in.Type)
	if err != nil {
		l.Errorf("[MarkAllRead] UpdateAllRead err: %v userId: %d type: %d", err, in.UserId, in.Type)
		return nil, err
	}

	return &pb.MarkAllReadResponse{}, nil
}
