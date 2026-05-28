package reply

import (
	"context"

	"ThinkTalk/application/reply/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateReplyRequest  = pb.CreateReplyRequest
	CreateReplyResponse = pb.CreateReplyResponse
	DeleteReplyRequest  = pb.DeleteReplyRequest
	DeleteReplyResponse = pb.DeleteReplyResponse
	ReplyDetailRequest  = pb.ReplyDetailRequest
	ReplyDetailResponse = pb.ReplyDetailResponse
	ReplyListRequest    = pb.ReplyListRequest
	ReplyListResponse   = pb.ReplyListResponse
	ReplyCountRequest   = pb.ReplyCountRequest
	ReplyCountResponse  = pb.ReplyCountResponse
	ReplyItem           = pb.ReplyItem

	Reply interface {
		CreateReply(ctx context.Context, in *CreateReplyRequest, opts ...grpc.CallOption) (*CreateReplyResponse, error)
		DeleteReply(ctx context.Context, in *DeleteReplyRequest, opts ...grpc.CallOption) (*DeleteReplyResponse, error)
		ReplyDetail(ctx context.Context, in *ReplyDetailRequest, opts ...grpc.CallOption) (*ReplyDetailResponse, error)
		ReplyList(ctx context.Context, in *ReplyListRequest, opts ...grpc.CallOption) (*ReplyListResponse, error)
		ReplyCount(ctx context.Context, in *ReplyCountRequest, opts ...grpc.CallOption) (*ReplyCountResponse, error)
	}

	defaultReply struct {
		cli zrpc.Client
	}
)

func NewReply(cli zrpc.Client) Reply {
	return &defaultReply{cli: cli}
}

func (m *defaultReply) CreateReply(ctx context.Context, in *CreateReplyRequest, opts ...grpc.CallOption) (*CreateReplyResponse, error) {
	client := pb.NewReplyClient(m.cli.Conn())
	return client.CreateReply(ctx, in, opts...)
}

func (m *defaultReply) DeleteReply(ctx context.Context, in *DeleteReplyRequest, opts ...grpc.CallOption) (*DeleteReplyResponse, error) {
	client := pb.NewReplyClient(m.cli.Conn())
	return client.DeleteReply(ctx, in, opts...)
}

func (m *defaultReply) ReplyDetail(ctx context.Context, in *ReplyDetailRequest, opts ...grpc.CallOption) (*ReplyDetailResponse, error) {
	client := pb.NewReplyClient(m.cli.Conn())
	return client.ReplyDetail(ctx, in, opts...)
}

func (m *defaultReply) ReplyList(ctx context.Context, in *ReplyListRequest, opts ...grpc.CallOption) (*ReplyListResponse, error) {
	client := pb.NewReplyClient(m.cli.Conn())
	return client.ReplyList(ctx, in, opts...)
}

func (m *defaultReply) ReplyCount(ctx context.Context, in *ReplyCountRequest, opts ...grpc.CallOption) (*ReplyCountResponse, error) {
	client := pb.NewReplyClient(m.cli.Conn())
	return client.ReplyCount(ctx, in, opts...)
}
