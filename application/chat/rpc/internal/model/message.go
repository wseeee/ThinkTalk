package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID             int64 `gorm:"primary_key"`
	ConversationID int64
	SenderID       int64
	ReceiverID     int64
	Content        string
	MsgType        int
	IsRead         int
	Status         int
	CreateTime     time.Time
}

func (m *Message) TableName() string {
	return "message"
}

type MessageModel struct {
	db *gorm.DB
}

func NewMessageModel(db *gorm.DB) *MessageModel {
	return &MessageModel{db: db}
}

func (m *MessageModel) Insert(ctx context.Context, data *Message) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *MessageModel) FindByConversationId(ctx context.Context, conversationId int64, cursor, limit int64) ([]*Message, error) {
	var result []*Message
	query := m.db.WithContext(ctx).Where("conversation_id = ? AND status = 0", conversationId)
	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}
	err := query.Order("id desc").Limit(int(limit)).Find(&result).Error
	return result, err
}

func (m *MessageModel) MarkReadByConversation(ctx context.Context, conversationId int64, receiverId int64) error {
	return m.db.WithContext(ctx).Model(&Message{}).
		Where("conversation_id = ? AND receiver_id = ? AND is_read = 0", conversationId, receiverId).
		Update("is_read", 1).Error
}

func (m *MessageModel) CountByConversationId(ctx context.Context, conversationId int64) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&Message{}).
		Where("conversation_id = ? AND status = 0", conversationId).
		Count(&count).Error
	return count, err
}
