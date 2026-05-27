package logic

import (
	"context"

	"ThinkTalk/application/message/code"
	"ThinkTalk/application/message/rpc/internal/svc"
	"ThinkTalk/application/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkReadLogic) MarkRead(in *pb.MarkReadRequest) (*pb.MarkReadResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.NotificationId == 0 {
		return nil, code.NotificationIdEmpty
	}

	notif, err := l.svcCtx.NotificationModel.FindOne(l.ctx, in.NotificationId)
	if err != nil {
		l.Errorf("[MarkRead] FindOne err: %v notificationId: %d", err, in.NotificationId)
		return nil, err
	}
	if notif == nil {
		return nil, code.NotificationNotFound
	}

	err = l.svcCtx.NotificationModel.UpdateRead(l.ctx, in.NotificationId, in.UserId)
	if err != nil {
		l.Errorf("[MarkRead] UpdateRead err: %v notificationId: %d", err, in.NotificationId)
		return nil, err
	}

	return &pb.MarkReadResponse{}, nil
}
