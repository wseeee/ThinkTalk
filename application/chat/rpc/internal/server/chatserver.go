package server

import (
	"context"

	"ThinkTalk/application/chat/rpc/internal/logic"
	"ThinkTalk/application/chat/rpc/internal/svc"
	"ThinkTalk/application/chat/rpc/pb"
)

type ChatServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedChatServer
}

func NewChatServer(svcCtx *svc.ServiceContext) *ChatServer {
	return &ChatServer{svcCtx: svcCtx}
}

func (s *ChatServer) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	l := logic.NewSendMessageLogic(ctx, s.svcCtx)
	return l.SendMessage(in)
}

func (s *ChatServer) Conversations(ctx context.Context, in *pb.ConversationsRequest) (*pb.ConversationsResponse, error) {
	l := logic.NewConversationsLogic(ctx, s.svcCtx)
	return l.Conversations(in)
}

func (s *ChatServer) Messages(ctx context.Context, in *pb.MessagesRequest) (*pb.MessagesResponse, error) {
	l := logic.NewMessagesLogic(ctx, s.svcCtx)
	return l.Messages(in)
}

func (s *ChatServer) MarkRead(ctx context.Context, in *pb.MarkReadRequest) (*pb.MarkReadResponse, error) {
	l := logic.NewMarkReadLogic(ctx, s.svcCtx)
	return l.MarkRead(in)
}

func (s *ChatServer) UnreadCount(ctx context.Context, in *pb.UnreadCountRequest) (*pb.UnreadCountResponse, error) {
	l := logic.NewUnreadCountLogic(ctx, s.svcCtx)
	return l.UnreadCount(in)
}
