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
