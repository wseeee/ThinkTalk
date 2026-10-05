package logic

import (
	"context"
	"math"

	"ThinkTalk/application/concerned/code"
	"ThinkTalk/application/concerned/rpc/internal/svc"
	"ThinkTalk/application/concerned/rpc/internal/types"
	"ThinkTalk/application/concerned/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConcernedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConcernedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConcernedListLogic {
	return &ConcernedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConcernedListLogic) ConcernedList(in *pb.ConcernedListRequest) (*pb.ConcernedListResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.PageSize > types.MaxPageSize {
		in.PageSize = types.MaxPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = math.MaxInt64
	}

	records, err := l.svcCtx.ConcernedRecordModel.FindByUserId(l.ctx, in.UserId, in.BizId, in.Cursor, in.PageSize+1)
	if err != nil {
		l.Errorf("[ConcernedList] FindByUserId err: %v userId: %d", err, in.UserId)
		return nil, err
	}

	var isEnd bool
	if len(records) > int(in.PageSize) {
		records = records[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(records) == 0 {
		return &pb.ConcernedListResponse{IsEnd: true}, nil
	}

	items := make([]*pb.ConcernedItem, 0, len(records))
	for _, r := range records {
		items = append(items, &pb.ConcernedItem{
			Id:         r.ID,
			BizId:      r.BizID,
			ObjId:      r.ObjID,
			CreateTime: r.CreateTime.Unix(),
		})
	}

	cursor := records[len(records)-1].ID
	return &pb.ConcernedListResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}
