package logic

import (
	"context"
	"encoding/json"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/like/rpc/like"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThumbupLogic) Thumbup(req *types.ThumbupRequest) (*types.ThumbupResponse, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.LikeRPC.Thumbup(l.ctx, &like.ThumbupRequest{
		BizId:    req.BizId,
		ObjId:    req.ObjId,
		UserId:   userId,
		LikeType: req.LikeType,
	})
	if err != nil {
		return nil, err
	}

	return &types.ThumbupResponse{
		BizId:      resp.BizId,
		ObjId:      resp.ObjId,
		LikeNum:    resp.LikeNum,
		DislikeNum: resp.DislikeNum,
	}, nil
}
