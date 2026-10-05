package concerned

import (
	"context"

	"ThinkTalk/application/concerned/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AddConcernedRequest    = pb.AddConcernedRequest
	AddConcernedResponse   = pb.AddConcernedResponse
	CancelConcernedRequest = pb.CancelConcernedRequest
	CancelConcernedResponse = pb.CancelConcernedResponse
	IsConcernedRequest     = pb.IsConcernedRequest
	IsConcernedResponse    = pb.IsConcernedResponse
	ConcernedListRequest   = pb.ConcernedListRequest
	ConcernedListResponse  = pb.ConcernedListResponse
	ConcernedCountRequest  = pb.ConcernedCountRequest
	ConcernedCountResponse = pb.ConcernedCountResponse

	Concerned interface {
		AddConcerned(ctx context.Context, in *AddConcernedRequest, opts ...grpc.CallOption) (*AddConcernedResponse, error)
		CancelConcerned(ctx context.Context, in *CancelConcernedRequest, opts ...grpc.CallOption) (*CancelConcernedResponse, error)
		IsConcerned(ctx context.Context, in *IsConcernedRequest, opts ...grpc.CallOption) (*IsConcernedResponse, error)
		ConcernedList(ctx context.Context, in *ConcernedListRequest, opts ...grpc.CallOption) (*ConcernedListResponse, error)
		ConcernedCount(ctx context.Context, in *ConcernedCountRequest, opts ...grpc.CallOption) (*ConcernedCountResponse, error)
	}

	defaultConcerned struct {
		cli zrpc.Client
	}
)

func NewConcerned(cli zrpc.Client) Concerned {
	return &defaultConcerned{cli: cli}
}

func (m *defaultConcerned) AddConcerned(ctx context.Context, in *AddConcernedRequest, opts ...grpc.CallOption) (*AddConcernedResponse, error) {
	client := pb.NewConcernedClient(m.cli.Conn())
	return client.AddConcerned(ctx, in, opts...)
}

func (m *defaultConcerned) CancelConcerned(ctx context.Context, in *CancelConcernedRequest, opts ...grpc.CallOption) (*CancelConcernedResponse, error) {
	client := pb.NewConcernedClient(m.cli.Conn())
	return client.CancelConcerned(ctx, in, opts...)
}

func (m *defaultConcerned) IsConcerned(ctx context.Context, in *IsConcernedRequest, opts ...grpc.CallOption) (*IsConcernedResponse, error) {
	client := pb.NewConcernedClient(m.cli.Conn())
	return client.IsConcerned(ctx, in, opts...)
}

func (m *defaultConcerned) ConcernedList(ctx context.Context, in *ConcernedListRequest, opts ...grpc.CallOption) (*ConcernedListResponse, error) {
	client := pb.NewConcernedClient(m.cli.Conn())
	return client.ConcernedList(ctx, in, opts...)
}

func (m *defaultConcerned) ConcernedCount(ctx context.Context, in *ConcernedCountRequest, opts ...grpc.CallOption) (*ConcernedCountResponse, error) {
	client := pb.NewConcernedClient(m.cli.Conn())
	return client.ConcernedCount(ctx, in, opts...)
}
