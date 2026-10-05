package logic

import (
	"context"

	"ThinkTalk/application/chat/code"
	"ThinkTalk/application/chat/rpc/internal/svc"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadCountLogic {
	return &UnreadCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnreadCountLogic) UnreadCount(in *pb.UnreadCountRequest) (*pb.UnreadCountResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	total, err := l.svcCtx.ConversationModel.CountUnread(l.ctx, in.UserId)
	if err != nil {
		l.Errorf("[UnreadCount] CountUnread err: %v userId: %d", err, in.UserId)
		return nil, err
	}

	return &pb.UnreadCountResponse{Total: total}, nil
}
