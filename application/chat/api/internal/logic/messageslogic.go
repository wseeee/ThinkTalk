package logic

import (
	"context"

	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/application/chat/api/internal/types"
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

func (l *MessagesLogic) Messages(userId int64, req *types.MessagesRequest) (*types.MessagesResponse, error) {
	if req.PageSize == 0 {
		req.PageSize = types.DefaultPageSize
	}

	resp, err := l.svcCtx.Chat.Messages(l.ctx, &pb.MessagesRequest{
		ConversationId: req.ConversationId,
		Cursor:         req.Cursor,
		PageSize:       req.PageSize,
	})
	if err != nil {
		l.Errorf("[Messages] rpc err: %v convId: %d", err, req.ConversationId)
		return nil, err
	}

	items := make([]*types.MessageItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.MessageItem{
			Id:             item.Id,
			ConversationId: item.ConversationId,
			SenderId:       item.SenderId,
			ReceiverId:     item.ReceiverId,
			Content:        item.Content,
			MsgType:        item.MsgType,
			IsRead:         item.IsRead,
			CreateTime:     item.CreateTime,
		})
	}

	return &types.MessagesResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
