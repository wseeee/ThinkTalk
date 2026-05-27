package logic

import (
	"context"
	"encoding/json"
	"time"

	"ThinkTalk/application/chat/mq/internal/model"
	"ThinkTalk/application/chat/mq/internal/svc"
	"ThinkTalk/application/chat/mq/internal/types"

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
	return &ConsumeLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConsumeLogic) Consume(ctx context.Context, key, val string) error {
	l.Infof("[Consume] key: %s val: %s", key, val)

	var msg types.ChatMsg
	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		l.Errorf("[Consume] unmarshal err: %v", err)
		return err
	}

	now := time.Now()
	trimmed := truncateMsg(msg.Content)

	var conversationId int64
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 收发双方的会话分别维护
		senderConv := &model.Conversation{
			UserID:          msg.SenderId,
			TargetUserID:    msg.ReceiverId,
			LastMessage:     trimmed,
			LastMessageTime: now,
			CreateTime:      now,
			UpdateTime:      now,
		}
		if err := model.NewConversationModel(tx).UpsertConversation(ctx, senderConv); err != nil {
			return err
		}

		receiverConv := &model.Conversation{
			UserID:          msg.ReceiverId,
			TargetUserID:    msg.SenderId,
			LastMessage:     trimmed,
			LastMessageTime: now,
			CreateTime:      now,
			UpdateTime:      now,
		}
		if err := model.NewConversationModel(tx).UpsertConversation(ctx, receiverConv); err != nil {
			return err
		}

		// 2. 获取发送者会话ID
		conv, err := model.NewConversationModel(tx).FindByUserIDAndTargetID(ctx, msg.SenderId, msg.ReceiverId)
		if err != nil || conv == nil {
			if err != nil {
				return err
			}
			return gorm.ErrRecordNotFound
		}
		conversationId = conv.ID

		// 3. 写入消息
		msgData := &model.Message{
			ConversationID: conversationId,
			SenderID:       msg.SenderId,
			ReceiverID:     msg.ReceiverId,
			Content:        msg.Content,
			MsgType:        int(msg.MsgType),
			IsRead:         0,
			Status:         0,
			CreateTime:     now,
		}
		return model.NewMessageModel(tx).Insert(ctx, msgData)
	})
	if err != nil {
		l.Errorf("[Consume] transaction err: %v msg: %+v", err, msg)
		return err
	}

	l.Infof("[Consume] message persisted convId: %d sender: %d receiver: %d", conversationId, msg.SenderId, msg.ReceiverId)
	return nil
}

func truncateMsg(s string) string {
	if len(s) > 100 {
		return s[:100] + "..."
	}
	return s
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewConsumeLogic(ctx, svcCtx)),
	}
}
