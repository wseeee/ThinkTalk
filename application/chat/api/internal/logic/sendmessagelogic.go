package logic

import (
	"context"

	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/application/chat/api/internal/types"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SendMessageLogic) SendMessage(userId int64, req *types.SendMessageRequest) (*types.SendMessageResponse, error) {
	_, err := l.svcCtx.Chat.SendMessage(l.ctx, &pb.SendMessageRequest{
		SenderId:   userId,
		ReceiverId: req.ReceiverId,
		Content:    req.Content,
		MsgType:    req.MsgType,
	})
	if err != nil {
		l.Errorf("[SendMessage] rpc err: %v req: %+v", err, req)
		return nil, err
	}

	return &types.SendMessageResponse{}, nil
}
