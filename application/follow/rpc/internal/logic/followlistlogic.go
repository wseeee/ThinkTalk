package logic

import (
	"context"

	"ThinkTalk/application/follow/rpc/internal/svc"
	"ThinkTalk/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 关注列表
func (l *FollowListLogic) FollowList(in *pb.FollowListRequest) (*pb.FollowListResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.FollowListResponse{}, nil
}
