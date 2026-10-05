package logic

import (
	"context"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FollowLogic) Follow(userId int64, req *types.FollowRequest) (*types.FollowResponse, error) {
	_, err := l.svcCtx.FollowRPC.Follow(l.ctx, &pb.FollowRequest{
		UserId:         userId,
		FollowedUserId: req.FollowedUserId,
	})
	if err != nil {
		l.Errorf("[Follow] rpc err: %v", err)
		return nil, err
	}
	return &types.FollowResponse{}, nil
}

func (l *FollowLogic) UnFollow(userId int64, req *types.UnfollowRequest) (*types.UnfollowResponse, error) {
	_, err := l.svcCtx.FollowRPC.UnFollow(l.ctx, &pb.UnFollowRequest{
		UserId:         userId,
		FollowedUserId: req.FollowedUserId,
	})
	if err != nil {
		l.Errorf("[UnFollow] rpc err: %v", err)
		return nil, err
	}
	return &types.UnfollowResponse{}, nil
}

func (l *FollowLogic) FollowList(userId int64, req *types.FollowListRequest) (*types.FollowListResponse, error) {
	resp, err := l.svcCtx.FollowRPC.FollowList(l.ctx, &pb.FollowListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[FollowList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.FollowItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.FollowItem{
			Id:             item.Id,
			FollowedUserId: item.FollowedUserId,
			CreateTime:     item.CreateTime,
			FansCount:      item.FansCount,
		})
	}
	return &types.FollowListResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}

func (l *FollowLogic) FansList(userId int64, req *types.FansListRequest) (*types.FansListResponse, error) {
	resp, err := l.svcCtx.FollowRPC.FansList(l.ctx, &pb.FansListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[FansList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.FollowItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.FollowItem{
			Id:             item.UserId,
			FollowedUserId: item.FansUserId,
			CreateTime:     item.CreateTime,
			FansCount:      item.FansCount,
		})
	}
	return &types.FansListResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
