package server

import (
	"context"

	"ThinkTalk/application/reply/rpc/internal/logic"
	"ThinkTalk/application/reply/rpc/internal/svc"
	"ThinkTalk/application/reply/rpc/pb"
)

type ReplyServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedReplyServer
}

func NewReplyServer(svcCtx *svc.ServiceContext) *ReplyServer {
	return &ReplyServer{svcCtx: svcCtx}
}

// 发表评论
func (s *ReplyServer) CreateReply(ctx context.Context, in *pb.CreateReplyRequest) (*pb.CreateReplyResponse, error) {
	l := logic.NewCreateReplyLogic(ctx, s.svcCtx)
	return l.CreateReply(in)
}

// 删除评论
func (s *ReplyServer) DeleteReply(ctx context.Context, in *pb.DeleteReplyRequest) (*pb.DeleteReplyResponse, error) {
	l := logic.NewDeleteReplyLogic(ctx, s.svcCtx)
	return l.DeleteReply(in)
}

// 评论详情
func (s *ReplyServer) ReplyDetail(ctx context.Context, in *pb.ReplyDetailRequest) (*pb.ReplyDetailResponse, error) {
	l := logic.NewReplyDetailLogic(ctx, s.svcCtx)
	return l.ReplyDetail(in)
}

// 评论列表
func (s *ReplyServer) ReplyList(ctx context.Context, in *pb.ReplyListRequest) (*pb.ReplyListResponse, error) {
	l := logic.NewReplyListLogic(ctx, s.svcCtx)
	return l.ReplyList(in)
}

// 评论计数
func (s *ReplyServer) ReplyCount(ctx context.Context, in *pb.ReplyCountRequest) (*pb.ReplyCountResponse, error) {
	l := logic.NewReplyCountLogic(ctx, s.svcCtx)
	return l.ReplyCount(in)
}
