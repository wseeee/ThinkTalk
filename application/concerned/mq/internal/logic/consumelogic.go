package logic

import (
	"context"
	"encoding/json"
	"time"

	"ThinkTalk/application/concerned/mq/internal/model"
	"ThinkTalk/application/concerned/mq/internal/svc"
	"ThinkTalk/application/concerned/mq/internal/types"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"gorm.io/gorm"
)

type ConsumeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConsumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumeLogic {
	return &ConsumeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConsumeLogic) Consume(ctx context.Context, key, val string) error {
	l.Infof("[Consume] consume key: %s val: %s", key, val)

	var msg types.ConcernedMsg
	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		l.Errorf("[Consume] unmarshal msg error: %v", err)
		return err
	}

	switch msg.OpType {
	case types.OpTypeAdd:
		return l.addConcerned(ctx, &msg)
	case types.OpTypeCancel:
		return l.cancelConcerned(ctx, &msg)
	default:
		l.Errorf("[Consume] unknown opType: %d", msg.OpType)
		return nil
	}
}

func (l *ConsumeLogic) addConcerned(ctx context.Context, msg *types.ConcernedMsg) error {
	existing, err := l.svcCtx.ConcernedRecordModel.FindByBizIDObjIDUserID(ctx, msg.BizId, msg.ObjId, msg.UserId)
	if err != nil {
		l.Errorf("[addConcerned] find existing err: %v msg: %+v", err, msg)
		return err
	}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if existing != nil {
			if existing.Status == 0 {
				return nil
			}
			return model.NewConcernedRecordModel(tx).UpdateFields(ctx, existing.ID, map[string]interface{}{
				"status": 0,
			})
		}
		err := model.NewConcernedRecordModel(tx).Insert(ctx, &model.ConcernedRecord{
			BizID:      msg.BizId,
			ObjID:      msg.ObjId,
			UserID:     msg.UserId,
			Status:     0,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		if err != nil {
			return err
		}
		return model.NewConcernedCountModel(tx).IncrConcernedNum(ctx, msg.BizId, msg.ObjId)
	})
	if err != nil {
		l.Errorf("[addConcerned] transaction err: %v msg: %+v", err, msg)
		return err
	}

	l.Infof("[addConcerned] success bizId: %s objId: %d userId: %d", msg.BizId, msg.ObjId, msg.UserId)
	return nil
}

func (l *ConsumeLogic) cancelConcerned(ctx context.Context, msg *types.ConcernedMsg) error {
	existing, err := l.svcCtx.ConcernedRecordModel.FindByBizIDObjIDUserID(ctx, msg.BizId, msg.ObjId, msg.UserId)
	if err != nil {
		l.Errorf("[cancelConcerned] find existing err: %v msg: %+v", err, msg)
		return err
	}
	if existing == nil || existing.Status == 1 {
		return nil
	}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := model.NewConcernedRecordModel(tx).UpdateFields(ctx, existing.ID, map[string]interface{}{
			"status": 1,
		}); err != nil {
			return err
		}
		return model.NewConcernedCountModel(tx).DecrConcernedNum(ctx, msg.BizId, msg.ObjId)
	})
	if err != nil {
		l.Errorf("[cancelConcerned] transaction err: %v msg: %+v", err, msg)
		return err
	}

	l.Infof("[cancelConcerned] success bizId: %s objId: %d userId: %d", msg.BizId, msg.ObjId, msg.UserId)
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewConsumeLogic(ctx, svcCtx)),
	}
}
