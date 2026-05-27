package logic

import (
	"context"
	"encoding/json"
	"time"

	"ThinkTalk/application/reply/mq/internal/model"
	"ThinkTalk/application/reply/mq/internal/svc"
	"ThinkTalk/application/reply/mq/internal/types"

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

	var msg types.ReplyMsg
	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		l.Errorf("[Consume] unmarshal msg error: %v", err)
		return err
	}

	switch msg.OpType {
	case types.OpTypeCreate:
		return l.createReply(ctx, &msg)
	case types.OpTypeDelete:
		return l.deleteReply(ctx, &msg)
	default:
		l.Errorf("[Consume] unknown opType: %d", msg.OpType)
		return nil
	}
}

func (l *ConsumeLogic) createReply(ctx context.Context, msg *types.ReplyMsg) error {
	isRoot := msg.ParentId == 0

	reply := &model.Reply{
		BizID:         msg.BizId,
		TargetID:      msg.TargetId,
		ReplyUserID:   msg.ReplyUserId,
		BeReplyUserID: msg.BeReplyUserId,
		ParentID:      msg.ParentId,
		Content:       msg.Content,
		Status:        0,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := model.NewReplyModel(tx).Insert(ctx, reply); err != nil {
			return err
		}
		return model.NewReplyCountModel(tx).IncrReplyNum(ctx, msg.BizId, msg.TargetId, isRoot)
	})
	if err != nil {
		l.Errorf("[createReply] transaction err: %v msg: %+v", err, msg)
		return err
	}

	l.Infof("[createReply] success replyId: %d bizId: %s targetId: %d", reply.ID, msg.BizId, msg.TargetId)
	return nil
}

func (l *ConsumeLogic) deleteReply(ctx context.Context, msg *types.ReplyMsg) error {
	reply, err := l.svcCtx.ReplyModel.FindOne(ctx, msg.ReplyId)
	if err != nil {
		l.Errorf("[deleteReply] ReplyModel.FindOne err: %v replyId: %d", err, msg.ReplyId)
		return err
	}
	if reply == nil || reply.Status == 1 {
		l.Infof("[deleteReply] reply already deleted, replyId: %d", msg.ReplyId)
		return nil
	}

	isRoot := reply.ParentID == 0

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := model.NewReplyModel(tx).UpdateFields(ctx, msg.ReplyId, map[string]interface{}{
			"status": 1,
		}); err != nil {
			return err
		}
		return model.NewReplyCountModel(tx).DecrReplyNum(ctx, reply.BizID, reply.TargetID, isRoot)
	})
	if err != nil {
		l.Errorf("[deleteReply] transaction err: %v replyId: %d", err, msg.ReplyId)
		return err
	}

	l.Infof("[deleteReply] success replyId: %d", msg.ReplyId)
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewConsumeLogic(ctx, svcCtx)),
	}
}
