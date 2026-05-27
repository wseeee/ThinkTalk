package message

import (
	"context"

	"ThinkTalk/application/message/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	NotificationListRequest  = pb.NotificationListRequest
	NotificationItem         = pb.NotificationItem
	NotificationListResponse = pb.NotificationListResponse
	UnreadCountRequest       = pb.UnreadCountRequest
	UnreadCountResponse      = pb.UnreadCountResponse
	MarkReadRequest          = pb.MarkReadRequest
	MarkReadResponse         = pb.MarkReadResponse
	MarkAllReadRequest       = pb.MarkAllReadRequest
	MarkAllReadResponse      = pb.MarkAllReadResponse

	Message interface {
		NotificationList(ctx context.Context, in *NotificationListRequest, opts ...grpc.CallOption) (*NotificationListResponse, error)
		UnreadCount(ctx context.Context, in *UnreadCountRequest, opts ...grpc.CallOption) (*UnreadCountResponse, error)
		MarkRead(ctx context.Context, in *MarkReadRequest, opts ...grpc.CallOption) (*MarkReadResponse, error)
		MarkAllRead(ctx context.Context, in *MarkAllReadRequest, opts ...grpc.CallOption) (*MarkAllReadResponse, error)
	}

	defaultMessage struct {
		cli zrpc.Client
	}
)

func NewMessage(cli zrpc.Client) Message {
	return &defaultMessage{cli: cli}
}

func (m *defaultMessage) NotificationList(ctx context.Context, in *NotificationListRequest, opts ...grpc.CallOption) (*NotificationListResponse, error) {
	client := pb.NewMessageClient(m.cli.Conn())
	return client.NotificationList(ctx, in, opts...)
}

func (m *defaultMessage) UnreadCount(ctx context.Context, in *UnreadCountRequest, opts ...grpc.CallOption) (*UnreadCountResponse, error) {
	client := pb.NewMessageClient(m.cli.Conn())
	return client.UnreadCount(ctx, in, opts...)
}

func (m *defaultMessage) MarkRead(ctx context.Context, in *MarkReadRequest, opts ...grpc.CallOption) (*MarkReadResponse, error) {
	client := pb.NewMessageClient(m.cli.Conn())
	return client.MarkRead(ctx, in, opts...)
}

func (m *defaultMessage) MarkAllRead(ctx context.Context, in *MarkAllReadRequest, opts ...grpc.CallOption) (*MarkAllReadResponse, error) {
	client := pb.NewMessageClient(m.cli.Conn())
	return client.MarkAllRead(ctx, in, opts...)
}
