package logic

import (
	"context"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/concerned/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedLogic {
	return &ConcernedLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConcernedLogic) Add(userId int64, req *types.ConcernedAddRequest) (*types.ConcernedAddResponse, error) {
	_, err := l.svcCtx.ConcernedRPC.AddConcerned(l.ctx, &pb.AddConcernedRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[ConcernedAdd] rpc err: %v", err)
		return nil, err
	}
	return &types.ConcernedAddResponse{}, nil
}

func (l *ConcernedLogic) Cancel(userId int64, req *types.ConcernedCancelRequest) (*types.ConcernedCancelResponse, error) {
	_, err := l.svcCtx.ConcernedRPC.CancelConcerned(l.ctx, &pb.CancelConcernedRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[ConcernedCancel] rpc err: %v", err)
		return nil, err
	}
	return &types.ConcernedCancelResponse{}, nil
}

func (l *ConcernedLogic) Check(userId int64, req *types.ConcernedCheckRequest) (*types.ConcernedCheckResponse, error) {
	resp, err := l.svcCtx.ConcernedRPC.IsConcerned(l.ctx, &pb.IsConcernedRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[ConcernedCheck] rpc err: %v", err)
		return nil, err
	}
	return &types.ConcernedCheckResponse{IsConcerned: resp.IsConcerned}, nil
}

func (l *ConcernedLogic) ConcernedCount(req *types.ConcernedCountRequest) (*types.ConcernedCountResponse, error) {
	resp, err := l.svcCtx.ConcernedRPC.ConcernedCount(l.ctx, &pb.ConcernedCountRequest{
		BizId: req.BizId,
		ObjId: req.ObjId,
	})
	if err != nil {
		l.Errorf("[ConcernedCount] rpc err: %v", err)
		return nil, err
	}
	return &types.ConcernedCountResponse{ConcernedNum: resp.ConcernedNum}, nil
}

func (l *ConcernedLogic) List(userId int64, req *types.ConcernedListRequest) (*types.ConcernedListResponse, error) {
	resp, err := l.svcCtx.ConcernedRPC.ConcernedList(l.ctx, &pb.ConcernedListRequest{
		UserId:   userId,
		BizId:    req.BizId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[ConcernedList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.ConcernedItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.ConcernedItem{
			Id:         item.Id,
			BizId:      item.BizId,
			ObjId:      item.ObjId,
			CreateTime: item.CreateTime,
		})
	}
	return &types.ConcernedListResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}
