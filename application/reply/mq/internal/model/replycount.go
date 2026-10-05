package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ReplyCount struct {
	ID           int64 `gorm:"primary_key"`
	BizID        string
	TargetID     int64
	ReplyNum     int
	ReplyRootNum int
	CreateTime   time.Time
	UpdateTime   time.Time
}

func (m *ReplyCount) TableName() string {
	return "reply_count"
}

type ReplyCountModel struct {
	db *gorm.DB
}

func NewReplyCountModel(db *gorm.DB) *ReplyCountModel {
	return &ReplyCountModel{db: db}
}

func (m *ReplyCountModel) Insert(ctx context.Context, data *ReplyCount) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *ReplyCountModel) FindOne(ctx context.Context, id int64) (*ReplyCount, error) {
	var result ReplyCount
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *ReplyCountModel) IncrReplyNum(ctx context.Context, bizId string, targetId int64, isRoot bool) error {
	rootIncr := 0
	if isRoot {
		rootIncr = 1
	}
	return m.db.WithContext(ctx).
		Exec("INSERT INTO reply_count (biz_id, target_id, reply_num, reply_root_num) VALUES (?, ?, 1, ?) ON DUPLICATE KEY UPDATE reply_num = reply_num + 1, reply_root_num = reply_root_num + ?",
			bizId, targetId, rootIncr, rootIncr).
		Error
}

func (m *ReplyCountModel) DecrReplyNum(ctx context.Context, bizId string, targetId int64, isRoot bool) error {
	rootDecr := 0
	if isRoot {
		rootDecr = 1
	}
	return m.db.WithContext(ctx).
		Exec("UPDATE reply_count SET reply_num = reply_num - 1, reply_root_num = reply_root_num - ? WHERE biz_id = ? AND target_id = ? AND reply_num > 0",
			rootDecr, bizId, targetId).
		Error
}
