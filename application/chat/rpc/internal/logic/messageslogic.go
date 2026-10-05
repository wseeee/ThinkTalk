package logic

import (
	"context"

	"ThinkTalk/application/chat/rpc/internal/svc"
	"ThinkTalk/application/chat/rpc/internal/types"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessagesLogic {
	return &MessagesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MessagesLogic) Messages(in *pb.MessagesRequest) (*pb.MessagesResponse, error) {
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.PageSize > types.MaxPageSize {
		in.PageSize = types.MaxPageSize
	}

	msgs, err := l.svcCtx.MessageModel.FindByConversationId(l.ctx, in.ConversationId, in.Cursor, in.PageSize+1)
	if err != nil {
		l.Errorf("[Messages] FindByConversationId err: %v convId: %d", err, in.ConversationId)
		return nil, err
	}

	var isEnd bool
	if len(msgs) > int(in.PageSize) {
		msgs = msgs[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(msgs) == 0 {
		return &pb.MessagesResponse{IsEnd: true}, nil
	}

	items := make([]*pb.MessageItem, 0, len(msgs))
	for _, m := range msgs {
		items = append(items, &pb.MessageItem{
			Id:             m.ID,
			ConversationId: m.ConversationID,
			SenderId:       m.SenderID,
			ReceiverId:     m.ReceiverID,
			Content:        m.Content,
			MsgType:        int32(m.MsgType),
			IsRead:         m.IsRead == 1,
			CreateTime:     m.CreateTime.Unix(),
		})
	}

	cursor := msgs[len(msgs)-1].ID
	return &pb.MessagesResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}
