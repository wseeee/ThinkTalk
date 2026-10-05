package logic

import (
	"context"

	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/application/chat/api/internal/types"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MarkReadLogic) MarkRead(userId int64, req *types.MarkReadRequest) (*types.MarkReadResponse, error) {
	_, err := l.svcCtx.Chat.MarkRead(l.ctx, &pb.MarkReadRequest{
		UserId:         userId,
		ConversationId: req.ConversationId,
	})
	if err != nil {
		l.Errorf("[MarkRead] rpc err: %v convId: %d", err, req.ConversationId)
		return nil, err
	}

	return &types.MarkReadResponse{}, nil
}
