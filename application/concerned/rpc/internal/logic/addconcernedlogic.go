package logic

import (
	"context"
	"encoding/json"

	"ThinkTalk/application/concerned/code"
	"ThinkTalk/application/concerned/rpc/internal/svc"
	"ThinkTalk/application/concerned/rpc/internal/types"
	"ThinkTalk/application/concerned/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type AddConcernedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddConcernedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddConcernedLogic {
	return &AddConcernedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddConcernedLogic) AddConcerned(in *pb.AddConcernedRequest) (*pb.AddConcernedResponse, error) {
	if in.BizId == "" {
		return nil, code.BizIdEmpty
	}
	if in.ObjId == 0 {
		return nil, code.ObjIdEmpty
	}
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}

	msg := &types.ConcernedMsg{
		BizId:  in.BizId,
		ObjId:  in.ObjId,
		UserId: in.UserId,
		OpType: types.OpTypeAdd,
	}

	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Errorf("[AddConcerned] marshal msg: %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(l.ctx, string(data))
		if err != nil {
			l.Errorf("[AddConcerned] kq push data: %s error: %v", data, err)
		}
	})

	return &pb.AddConcernedResponse{}, nil
}
