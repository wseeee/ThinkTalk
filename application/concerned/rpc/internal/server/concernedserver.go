package server

import (
	"context"

	"ThinkTalk/application/concerned/rpc/internal/logic"
	"ThinkTalk/application/concerned/rpc/internal/svc"
	"ThinkTalk/application/concerned/rpc/pb"
)

type ConcernedServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedConcernedServer
}

func NewConcernedServer(svcCtx *svc.ServiceContext) *ConcernedServer {
	return &ConcernedServer{svcCtx: svcCtx}
}

// 收藏
func (s *ConcernedServer) AddConcerned(ctx context.Context, in *pb.AddConcernedRequest) (*pb.AddConcernedResponse, error) {
	l := logic.NewAddConcernedLogic(ctx, s.svcCtx)
	return l.AddConcerned(in)
}

// 取消收藏
func (s *ConcernedServer) CancelConcerned(ctx context.Context, in *pb.CancelConcernedRequest) (*pb.CancelConcernedResponse, error) {
	l := logic.NewCancelConcernedLogic(ctx, s.svcCtx)
	return l.CancelConcerned(in)
}

// 是否已收藏
func (s *ConcernedServer) IsConcerned(ctx context.Context, in *pb.IsConcernedRequest) (*pb.IsConcernedResponse, error) {
	l := logic.NewIsConcernedLogic(ctx, s.svcCtx)
	return l.IsConcerned(in)
}

// 收藏列表
func (s *ConcernedServer) ConcernedList(ctx context.Context, in *pb.ConcernedListRequest) (*pb.ConcernedListResponse, error) {
	l := logic.NewConcernedListLogic(ctx, s.svcCtx)
	return l.ConcernedList(in)
}

// 收藏计数
func (s *ConcernedServer) ConcernedCount(ctx context.Context, in *pb.ConcernedCountRequest) (*pb.ConcernedCountResponse, error) {
	l := logic.NewConcernedCountLogic(ctx, s.svcCtx)
	return l.ConcernedCount(in)
}
