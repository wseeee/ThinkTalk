package logic

import (
	"context"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/member/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberLogic {
	return &MemberLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MemberLogic) MemberInfo(userId int64) (*types.MemberInfoResponse, error) {
	resp, err := l.svcCtx.MemberRPC.MemberInfo(l.ctx, &pb.MemberInfoRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[MemberInfo] rpc err: %v", err)
		return nil, err
	}
	return &types.MemberInfoResponse{
		UserId:     resp.UserId,
		Level:      resp.Level,
		LevelName:  resp.LevelName,
		ExpireTime: resp.ExpireTime,
		Status:     resp.Status,
	}, nil
}

func (l *MemberLogic) CheckRight(userId int64, req *types.MemberRightRequest) (*types.MemberRightResponse, error) {
	resp, err := l.svcCtx.MemberRPC.CheckMemberRight(l.ctx, &pb.CheckMemberRightRequest{
		UserId:   userId,
		RightKey: req.RightKey,
	})
	if err != nil {
		l.Errorf("[CheckMemberRight] rpc err: %v", err)
		return nil, err
	}
	return &types.MemberRightResponse{
		HasRight: resp.HasRight,
		Level:    resp.Level,
	}, nil
}

func (l *MemberLogic) UpgradeMember(userId int64, req *types.UpgradeMemberRequest) (*types.UpgradeMemberResponse, error) {
	_, err := l.svcCtx.MemberRPC.UpgradeMember(l.ctx, &pb.UpgradeMemberRequest{
		UserId:        userId,
		Level:         req.Level,
		DurationDays:  req.DurationDays,
		TransactionId: req.TransactionId,
		Amount:        req.Amount,
		PayChannel:    req.PayChannel,
	})
	if err != nil {
		l.Errorf("[UpgradeMember] rpc err: %v", err)
		return nil, err
	}
	return &types.UpgradeMemberResponse{}, nil
}

func (l *MemberLogic) MemberOrderList(userId int64, req *types.MemberOrderListRequest) (*types.MemberOrderListResponse, error) {
	resp, err := l.svcCtx.MemberRPC.MemberOrderList(l.ctx, &pb.MemberOrderListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[MemberOrderList] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.MemberOrderItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.MemberOrderItem{
			Id:           item.Id,
			UserId:       item.UserId,
			Level:        item.Level,
			DurationDays: item.DurationDays,
			Amount:       item.Amount,
			PayChannel:   item.PayChannel,
			Status:       item.Status,
			CreateTime:   item.CreateTime,
		})
	}
	return &types.MemberOrderListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}
