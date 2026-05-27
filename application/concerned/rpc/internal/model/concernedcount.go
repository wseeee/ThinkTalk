package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ConcernedCount struct {
	ID           int64 `gorm:"primary_key"`
	BizID        string
	ObjID        int64
	ConcernedNum int
	CreateTime   time.Time
	UpdateTime   time.Time
}

func (m *ConcernedCount) TableName() string {
	return "concerned_count"
}

type ConcernedCountModel struct {
	db *gorm.DB
}

func NewConcernedCountModel(db *gorm.DB) *ConcernedCountModel {
	return &ConcernedCountModel{db: db}
}

func (m *ConcernedCountModel) FindByBizIDAndObjID(ctx context.Context, bizId string, objId int64) (*ConcernedCount, error) {
	var result ConcernedCount
	err := m.db.WithContext(ctx).Where("biz_id = ? AND obj_id = ?", bizId, objId).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *ConcernedCountModel) IncrConcernedNum(ctx context.Context, bizId string, objId int64) error {
	return m.db.WithContext(ctx).
		Exec("INSERT INTO concerned_count (biz_id, obj_id, concerned_num) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE concerned_num = concerned_num + 1", bizId, objId).
		Error
}

func (m *ConcernedCountModel) DecrConcernedNum(ctx context.Context, bizId string, objId int64) error {
	return m.db.WithContext(ctx).
		Exec("UPDATE concerned_count SET concerned_num = concerned_num - 1 WHERE biz_id = ? AND obj_id = ? AND concerned_num > 0", bizId, objId).
		Error
}
