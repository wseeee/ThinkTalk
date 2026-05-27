package server

import (
	"context"

	"ThinkTalk/application/message/rpc/internal/logic"
	"ThinkTalk/application/message/rpc/internal/svc"
	"ThinkTalk/application/message/rpc/pb"
)

type MessageServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedMessageServer
}

func NewMessageServer(svcCtx *svc.ServiceContext) *MessageServer {
	return &MessageServer{svcCtx: svcCtx}
}

// 通知列表
func (s *MessageServer) NotificationList(ctx context.Context, in *pb.NotificationListRequest) (*pb.NotificationListResponse, error) {
	l := logic.NewNotificationListLogic(ctx, s.svcCtx)
	return l.NotificationList(in)
}

// 未读计数
func (s *MessageServer) UnreadCount(ctx context.Context, in *pb.UnreadCountRequest) (*pb.UnreadCountResponse, error) {
	l := logic.NewUnreadCountLogic(ctx, s.svcCtx)
	return l.UnreadCount(in)
}

// 标记已读
func (s *MessageServer) MarkRead(ctx context.Context, in *pb.MarkReadRequest) (*pb.MarkReadResponse, error) {
	l := logic.NewMarkReadLogic(ctx, s.svcCtx)
	return l.MarkRead(in)
}

// 全部已读
func (s *MessageServer) MarkAllRead(ctx context.Context, in *pb.MarkAllReadRequest) (*pb.MarkAllReadResponse, error) {
	l := logic.NewMarkAllReadLogic(ctx, s.svcCtx)
	return l.MarkAllRead(in)
}
