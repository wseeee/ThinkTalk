package logic

import (
	"context"
	"encoding/json"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/like/rpc/like"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
	return &IsThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsThumbupLogic) IsThumbup(req *types.IsThumbupRequest) (*types.IsThumbupResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.LikeRPC.IsThumbup(l.ctx, &like.IsThumbupRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		UserId:   userId,
	})
	if err != nil {
		return nil, err
	}

	result := &types.IsThumbupResponse{}
	if thumbup, ok := resp.UserThumbups[req.TargetId]; ok {
		result.LikeType = thumbup.LikeType
		result.ThumbupTime = thumbup.ThumbupTime
	}
	return result, nil
}
