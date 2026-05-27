package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID            int64 `gorm:"primary_key"`
	UserID        int64
	Type          int32
	Title         string
	Content       string
	RefID         int64
	BizID         string
	TriggerUserID int64
	IsRead        int
	CreateTime    time.Time
	UpdateTime    time.Time
}

func (m *Notification) TableName() string {
	return "notification"
}

type NotificationModel struct {
	db *gorm.DB
}

func NewNotificationModel(db *gorm.DB) *NotificationModel {
	return &NotificationModel{db: db}
}

func (m *NotificationModel) FindOne(ctx context.Context, id int64) (*Notification, error) {
	var result Notification
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *NotificationModel) FindByUserId(ctx context.Context, userId int64, notifType int32, cursor, limit int64) ([]*Notification, error) {
	var result []*Notification
	query := m.db.WithContext(ctx).Where("user_id = ?", userId)

	if notifType > 0 {
		query = query.Where("type = ?", notifType)
	}
	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}

	err := query.Order("id desc").Limit(int(limit)).Find(&result).Error
	return result, err
}

func (m *NotificationModel) UpdateRead(ctx context.Context, id int64, userId int64) error {
	return m.db.WithContext(ctx).Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userId).
		Update("is_read", 1).Error
}

func (m *NotificationModel) UpdateAllRead(ctx context.Context, userId int64, notifType int32) error {
	query := m.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? AND is_read = 0", userId)
	if notifType > 0 {
		query = query.Where("type = ?", notifType)
	}
	return query.Update("is_read", 1).Error
}

func (m *NotificationModel) CountUnread(ctx context.Context, userId int64) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? AND is_read = 0", userId).
		Count(&count).Error
	return count, err
}

func (m *NotificationModel) CountUnreadByType(ctx context.Context, userId int64) (map[int32]int64, error) {
	type typeCount struct {
		Type  int32
		Count int64
	}
	var rows []typeCount
	err := m.db.WithContext(ctx).
		Model(&Notification{}).
		Select("type, count(*) as count").
		Where("user_id = ? AND is_read = 0", userId).
		Group("type").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int32]int64)
	for _, row := range rows {
		result[row.Type] = row.Count
	}
	return result, nil
}
