package server

import (
	"context"

	"ThinkTalk/application/member/rpc/internal/logic"
	"ThinkTalk/application/member/rpc/internal/svc"
	"ThinkTalk/application/member/rpc/pb"
)

type MemberServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedMemberServer
}

func NewMemberServer(svcCtx *svc.ServiceContext) *MemberServer {
	return &MemberServer{svcCtx: svcCtx}
}

func (s *MemberServer) MemberInfo(ctx context.Context, in *pb.MemberInfoRequest) (*pb.MemberInfoResponse, error) {
	return logic.NewMemberInfoLogic(ctx, s.svcCtx).MemberInfo(in)
}

func (s *MemberServer) UpgradeMember(ctx context.Context, in *pb.UpgradeMemberRequest) (*pb.UpgradeMemberResponse, error) {
	return logic.NewUpgradeMemberLogic(ctx, s.svcCtx).UpgradeMember(in)
}

func (s *MemberServer) CheckMemberRight(ctx context.Context, in *pb.CheckMemberRightRequest) (*pb.CheckMemberRightResponse, error) {
	return logic.NewCheckMemberRightLogic(ctx, s.svcCtx).CheckMemberRight(in)
}

func (s *MemberServer) MemberOrderList(ctx context.Context, in *pb.MemberOrderListRequest) (*pb.MemberOrderListResponse, error) {
	return logic.NewMemberOrderListLogic(ctx, s.svcCtx).MemberOrderList(in)
}
