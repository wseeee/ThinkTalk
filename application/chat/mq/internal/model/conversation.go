package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	ID              int64 `gorm:"primary_key"`
	UserID          int64
	TargetUserID    int64
	LastMessage     string
	LastMessageTime time.Time
	UnreadCount     int
	CreateTime      time.Time
	UpdateTime      time.Time
}

func (m *Conversation) TableName() string {
	return "conversation"
}

type ConversationModel struct {
	db *gorm.DB
}

func NewConversationModel(db *gorm.DB) *ConversationModel {
	return &ConversationModel{db: db}
}

func (m *ConversationModel) FindByUserIDAndTargetID(ctx context.Context, userId, targetId int64) (*Conversation, error) {
	var result Conversation
	err := m.db.WithContext(ctx).Where("user_id = ? AND target_user_id = ?", userId, targetId).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *ConversationModel) UpsertConversation(ctx context.Context, conv *Conversation) error {
	return m.db.WithContext(ctx).Exec(
		"INSERT INTO conversation (user_id, target_user_id, last_message, last_message_time, unread_count, create_time, update_time) VALUES (?, ?, ?, ?, 1, ?, ?) ON DUPLICATE KEY UPDATE last_message = ?, last_message_time = ?, unread_count = unread_count + 1, update_time = ?",
		conv.UserID, conv.TargetUserID, conv.LastMessage, conv.LastMessageTime, conv.CreateTime, conv.UpdateTime,
		conv.LastMessage, conv.LastMessageTime, conv.UpdateTime,
	).Error
}
