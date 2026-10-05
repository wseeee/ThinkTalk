package logic

import (
	"context"

	"ThinkTalk/application/concerned/code"
	"ThinkTalk/application/concerned/rpc/internal/svc"
	"ThinkTalk/application/concerned/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedCountLogic {
	return &ConcernedCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConcernedCountLogic) ConcernedCount(in *pb.ConcernedCountRequest) (*pb.ConcernedCountResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.ObjId == 0 {
		return nil, code.ObjIdEmpty
	}

	count, err := l.svcCtx.ConcernedCountModel.FindByBizIDAndObjID(l.ctx, in.BizId, in.ObjId)
	if err != nil {
		l.Errorf("[ConcernedCount] FindByBizIDAndObjID err: %v req: %+v", err, in)
		return nil, err
	}
	if count == nil {
		return &pb.ConcernedCountResponse{ConcernedNum: 0}, nil
	}

	return &pb.ConcernedCountResponse{ConcernedNum: int64(count.ConcernedNum)}, nil
}
