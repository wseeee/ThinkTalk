package logic

import (
	"context"

	"ThinkTalk/application/concerned/code"
	"ThinkTalk/application/concerned/rpc/internal/svc"
	"ThinkTalk/application/concerned/rpc/internal/types"
	"ThinkTalk/application/concerned/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsConcernedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsConcernedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsConcernedLogic {
	return &IsConcernedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsConcernedLogic) IsConcerned(in *pb.IsConcernedRequest) (*pb.IsConcernedResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.ObjId == 0 {
		return nil, code.ObjIdEmpty
	}
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	record, err := l.svcCtx.ConcernedRecordModel.FindByBizIDObjIDUserID(l.ctx, in.BizId, in.ObjId, in.UserId)
	if err != nil {
		l.Errorf("[IsConcerned] FindByBizIDObjIDUserID err: %v req: %+v", err, in)
		return nil, err
	}

	isConcerned := record != nil && record.Status == types.StatusConcerned
	return &pb.IsConcernedResponse{IsConcerned: isConcerned}, nil
}
