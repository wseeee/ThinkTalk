package logic

import (
	"context"

	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/application/chat/api/internal/types"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type WsMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWsMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WsMessageLogic {
	return &WsMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *WsMessageLogic) HandleMessage(userId int64, msg *types.WsInMessage) (*types.WsOutMessage, error) {
	_, err := l.svcCtx.Chat.SendMessage(l.ctx, &pb.SendMessageRequest{
		SenderId:   userId,
		ReceiverId: msg.ReceiverId,
		Content:    msg.Content,
		MsgType:    msg.MsgType,
	})
	if err != nil {
		l.Errorf("[WsMessage] rpc err: %v userId: %d", err, userId)
		return &types.WsOutMessage{Type: "error", Message: err.Error()}, err
	}

	return &types.WsOutMessage{Type: "ack", Message: "sent"}, nil
}
