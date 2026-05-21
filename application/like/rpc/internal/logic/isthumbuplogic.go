package logic

import (
	"context"

	"ThinkTalk/application/like/rpc/internal/model"
	"ThinkTalk/application/like/rpc/internal/svc"
	"ThinkTalk/application/like/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
	return &IsThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsThumbupLogic) IsThumbup(in *service.IsThumbupRequest) (*service.IsThumbupResponse, error) {
	// 查询点赞记录
	record, err := l.svcCtx.LikeRecordModel.FindOneByBizIdObjIdUserId(l.ctx, in.BizId, in.TargetId, in.UserId)
	if err != nil && err != model.ErrNotFound {
		l.Errorf("[IsThumbup] find like record error: %v", err)
		return nil, err
	}

	// 构造响应
	resp := &service.IsThumbupResponse{
		UserThumbups: make(map[int64]*service.UserThumbup),
	}

	if record != nil {
		resp.UserThumbups[in.TargetId] = &service.UserThumbup{
			UserId:      record.UserId,
			ThumbupTime: record.CreateTime.UnixMilli(),
			LikeType:    int32(record.LikeType),
		}
	}

	return resp, nil
}
