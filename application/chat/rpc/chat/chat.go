package chat

import (
	"context"

	"ThinkTalk/application/chat/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	SendMessageRequest      = pb.SendMessageRequest
	SendMessageResponse     = pb.SendMessageResponse
	ConversationsRequest    = pb.ConversationsRequest
	ConversationsResponse   = pb.ConversationsResponse
	MessagesRequest         = pb.MessagesRequest
	MessagesResponse        = pb.MessagesResponse
	MarkReadRequest         = pb.MarkReadRequest
	MarkReadResponse        = pb.MarkReadResponse
	UnreadCountRequest      = pb.UnreadCountRequest
	UnreadCountResponse     = pb.UnreadCountResponse

	Chat interface {
		SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
		Conversations(ctx context.Context, in *ConversationsRequest, opts ...grpc.CallOption) (*ConversationsResponse, error)
		Messages(ctx context.Context, in *MessagesRequest, opts ...grpc.CallOption) (*MessagesResponse, error)
		MarkRead(ctx context.Context, in *MarkReadRequest, opts ...grpc.CallOption) (*MarkReadResponse, error)
		UnreadCount(ctx context.Context, in *UnreadCountRequest, opts ...grpc.CallOption) (*UnreadCountResponse, error)
	}

	defaultChat struct {
		cli zrpc.Client
	}
)

func NewChat(cli zrpc.Client) Chat {
	return &defaultChat{cli: cli}
}

func (m *defaultChat) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.SendMessage(ctx, in, opts...)
}

func (m *defaultChat) Conversations(ctx context.Context, in *ConversationsRequest, opts ...grpc.CallOption) (*ConversationsResponse, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.Conversations(ctx, in, opts...)
}

func (m *defaultChat) Messages(ctx context.Context, in *MessagesRequest, opts ...grpc.CallOption) (*MessagesResponse, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.Messages(ctx, in, opts...)
}

func (m *defaultChat) MarkRead(ctx context.Context, in *MarkReadRequest, opts ...grpc.CallOption) (*MarkReadResponse, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.MarkRead(ctx, in, opts...)
}

func (m *defaultChat) UnreadCount(ctx context.Context, in *UnreadCountRequest, opts ...grpc.CallOption) (*UnreadCountResponse, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.UnreadCount(ctx, in, opts...)
}
