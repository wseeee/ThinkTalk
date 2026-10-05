package logic

import (
	"context"

	"ThinkTalk/application/chat/code"
	"ThinkTalk/application/chat/rpc/internal/svc"
	"ThinkTalk/application/chat/rpc/internal/types"
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

func (l *ConversationsLogic) Conversations(in *pb.ConversationsRequest) (*pb.ConversationsResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}

	convs, err := l.svcCtx.ConversationModel.FindByUserId(l.ctx, in.UserId, in.Cursor, in.PageSize+1)
	if err != nil {
		l.Errorf("[Conversations] FindByUserId err: %v userId: %d", err, in.UserId)
		return nil, err
	}

	var isEnd bool
	if len(convs) > int(in.PageSize) {
		convs = convs[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(convs) == 0 {
		return &pb.ConversationsResponse{IsEnd: true}, nil
	}

	items := make([]*pb.ConversationItem, 0, len(convs))
	for _, c := range convs {
		items = append(items, &pb.ConversationItem{
			Id:              c.ID,
			TargetUserId:    c.TargetUserID,
			LastMessage:     c.LastMessage,
			LastMessageTime: c.LastMessageTime.Unix(),
			UnreadCount:     int64(c.UnreadCount),
		})
	}

	cursor := convs[len(convs)-1].LastMessageTime.Unix()
	return &pb.ConversationsResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}
