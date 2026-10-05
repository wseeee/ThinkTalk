package logic

import (
	"context"

	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/application/chat/api/internal/types"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationsLogic {
	return &ConversationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConversationsLogic) Conversations(userId int64, req *types.ConversationsRequest) (*types.ConversationsResponse, error) {
	if req.PageSize == 0 {
		req.PageSize = types.DefaultPageSize
	}

	resp, err := l.svcCtx.Chat.Conversations(l.ctx, &pb.ConversationsRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[Conversations] rpc err: %v userId: %d", err, userId)
		return nil, err
	}

	items := make([]*types.ConversationItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.ConversationItem{
			Id:              item.Id,
			TargetUserId:    item.TargetUserId,
			LastMessage:     item.LastMessage,
			LastMessageTime: item.LastMessageTime,
			UnreadCount:     item.UnreadCount,
		})
	}

	return &types.ConversationsResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
