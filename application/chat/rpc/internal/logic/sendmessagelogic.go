package logic

import (
	"context"
	"encoding/json"

	"ThinkTalk/application/chat/code"
	"ThinkTalk/application/chat/rpc/internal/svc"
	"ThinkTalk/application/chat/rpc/internal/types"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SendMessageLogic) SendMessage(in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if in.SenderId == 0 {
		return nil, code.SenderIdEmpty
	}
	if in.ReceiverId == 0 {
		return nil, code.ReceiverIdEmpty
	}
	if in.Content == "" {
		return nil, code.ContentEmpty
	}
	if in.SenderId == in.ReceiverId {
		return nil, code.CannotSelfChat
	}

	msg := &types.ChatMsg{
		SenderId:   in.SenderId,
		ReceiverId: in.ReceiverId,
		Content:    in.Content,
		MsgType:    in.MsgType,
	}

	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Errorf("[SendMessage] marshal err: %v msg: %+v", err, msg)
			return
		}
		if err := l.svcCtx.KqPusherClient.Push(l.ctx, string(data)); err != nil {
			l.Errorf("[SendMessage] kq push err: %v data: %s", err, data)
		}
	})

	return &pb.SendMessageResponse{}, nil
}
