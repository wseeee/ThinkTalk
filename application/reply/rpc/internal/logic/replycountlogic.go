package logic

import (
	"context"

	"ThinkTalk/application/reply/code"
	"ThinkTalk/application/reply/rpc/internal/svc"
	"ThinkTalk/application/reply/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyCountLogic {
	return &ReplyCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyCountLogic) ReplyCount(in *pb.ReplyCountRequest) (*pb.ReplyCountResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.TargetId == 0 {
		return nil, code.TargetIdEmpty
	}

	total, err := l.svcCtx.ReplyModel.CountByBizIDAndTargetID(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		l.Errorf("[ReplyCount] CountByBizIDAndTargetID err: %v req: %+v", err, in)
		return nil, err
	}

	rootTotal, err := l.svcCtx.ReplyModel.CountRootByBizIDAndTargetID(l.ctx, in.BizId, in.TargetId)
	if err != nil {
		l.Errorf("[ReplyCount] CountRootByBizIDAndTargetID err: %v req: %+v", err, in)
		return nil, err
	}

	return &pb.ReplyCountResponse{
		ReplyNum:     total,
		ReplyRootNum: rootTotal,
	}, nil
}
