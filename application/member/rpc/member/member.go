package member

import (
	"context"

	"ThinkTalk/application/member/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	MemberInfoRequest       = pb.MemberInfoRequest
	MemberInfoResponse      = pb.MemberInfoResponse
	UpgradeMemberRequest    = pb.UpgradeMemberRequest
	UpgradeMemberResponse   = pb.UpgradeMemberResponse
	CheckMemberRightRequest = pb.CheckMemberRightRequest
	CheckMemberRightResponse = pb.CheckMemberRightResponse
	MemberOrderListRequest  = pb.MemberOrderListRequest
	MemberOrderListResponse = pb.MemberOrderListResponse

	Member interface {
		MemberInfo(ctx context.Context, in *MemberInfoRequest, opts ...grpc.CallOption) (*MemberInfoResponse, error)
		UpgradeMember(ctx context.Context, in *UpgradeMemberRequest, opts ...grpc.CallOption) (*UpgradeMemberResponse, error)
		CheckMemberRight(ctx context.Context, in *CheckMemberRightRequest, opts ...grpc.CallOption) (*CheckMemberRightResponse, error)
		MemberOrderList(ctx context.Context, in *MemberOrderListRequest, opts ...grpc.CallOption) (*MemberOrderListResponse, error)
	}

	defaultMember struct {
		cli zrpc.Client
	}
)

func NewMember(cli zrpc.Client) Member {
	return &defaultMember{cli: cli}
}

func (m *defaultMember) MemberInfo(ctx context.Context, in *MemberInfoRequest, opts ...grpc.CallOption) (*MemberInfoResponse, error) {
	client := pb.NewMemberClient(m.cli.Conn())
	return client.MemberInfo(ctx, in, opts...)
}

func (m *defaultMember) UpgradeMember(ctx context.Context, in *UpgradeMemberRequest, opts ...grpc.CallOption) (*UpgradeMemberResponse, error) {
	client := pb.NewMemberClient(m.cli.Conn())
	return client.UpgradeMember(ctx, in, opts...)
}

func (m *defaultMember) CheckMemberRight(ctx context.Context, in *CheckMemberRightRequest, opts ...grpc.CallOption) (*CheckMemberRightResponse, error) {
	client := pb.NewMemberClient(m.cli.Conn())
	return client.CheckMemberRight(ctx, in, opts...)
}

func (m *defaultMember) MemberOrderList(ctx context.Context, in *MemberOrderListRequest, opts ...grpc.CallOption) (*MemberOrderListResponse, error) {
	client := pb.NewMemberClient(m.cli.Conn())
	return client.MemberOrderList(ctx, in, opts...)
}
