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

func (m *ConversationModel) FindByUserId(ctx context.Context, userId int64, cursor, limit int64) ([]*Conversation, error) {
	var result []*Conversation
	query := m.db.WithContext(ctx).Where("user_id = ?", userId)
	if cursor > 0 {
		query = query.Where("last_message_time < ?", time.Unix(cursor, 0))
	}
	err := query.Order("last_message_time desc").Limit(int(limit)).Find(&result).Error
	return result, err
}

func (m *ConversationModel) UpsertConversation(ctx context.Context, conv *Conversation) error {
	return m.db.WithContext(ctx).Exec(
		"INSERT INTO conversation (user_id, target_user_id, last_message, last_message_time, unread_count, create_time, update_time) VALUES (?, ?, ?, ?, 1, ?, ?) ON DUPLICATE KEY UPDATE last_message = ?, last_message_time = ?, unread_count = unread_count + 1, update_time = ?",
		conv.UserID, conv.TargetUserID, conv.LastMessage, conv.LastMessageTime, conv.CreateTime, conv.UpdateTime,
		conv.LastMessage, conv.LastMessageTime, conv.UpdateTime,
	).Error
}

func (m *ConversationModel) ClearUnread(ctx context.Context, userId, conversationId int64) error {
	return m.db.WithContext(ctx).Model(&Conversation{}).
		Where("id = ? AND user_id = ?", conversationId, userId).
		Update("unread_count", 0).Error
}

func (m *ConversationModel) CountUnread(ctx context.Context, userId int64) (int64, error) {
	var total int64
	err := m.db.WithContext(ctx).Model(&Conversation{}).
		Where("user_id = ?", userId).
		Select("COALESCE(SUM(unread_count), 0)").
		Scan(&total).Error
	return total, err
}
