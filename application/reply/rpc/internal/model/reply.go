package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Reply struct {
	ID            int64 `gorm:"primary_key"`
	BizID         string
	TargetID      int64
	ReplyUserID   int64
	BeReplyUserID int64
	ParentID      int64
	Content       string
	Status        int
	LikeNum       int
	CreateTime    time.Time
	UpdateTime    time.Time
}

func (m *Reply) TableName() string {
	return "reply"
}

type ReplyModel struct {
	db *gorm.DB
}

func NewReplyModel(db *gorm.DB) *ReplyModel {
	return &ReplyModel{db: db}
}

func (m *ReplyModel) Insert(ctx context.Context, data *Reply) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *ReplyModel) FindOne(ctx context.Context, id int64) (*Reply, error) {
	var result Reply
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &result, err
}

func (m *ReplyModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&Reply{}).Where("id = ?", id).Updates(values).Error
}

func (m *ReplyModel) FindRootReplies(ctx context.Context, bizId string, targetId int64, sortType int, cursor, limit int64) ([]*Reply, error) {
	var result []*Reply
	query := m.db.WithContext(ctx).
		Where("biz_id = ? AND target_id = ? AND parent_id = 0 AND status = ?", bizId, targetId, 0)

	if cursor > 0 {
		if sortType == 1 {
			query = query.Where("(like_num < ? OR (like_num = ? AND id < ?))", cursor, cursor, cursor)
		} else {
			query = query.Where("id < ?", cursor)
		}
	}

	if sortType == 1 {
		query = query.Order("like_num desc, id desc")
	} else {
		query = query.Order("id desc")
	}

	err := query.Limit(int(limit)).Find(&result).Error
	return result, err
}

func (m *ReplyModel) FindByParentIDs(ctx context.Context, parentIds []int64) ([]*Reply, error) {
	var result []*Reply
	err := m.db.WithContext(ctx).
		Where("parent_id IN ? AND status = ?", parentIds, 0).
		Order("id asc").
		Find(&result).Error
	return result, err
}

func (m *ReplyModel) FindByBizIDAndTargetID(ctx context.Context, bizId string, targetId int64, cursor, limit int64) ([]*Reply, error) {
	var result []*Reply
	err := m.db.WithContext(ctx).
		Where("biz_id = ? AND target_id = ? AND status = ? AND id < ?", bizId, targetId, 0, cursor).
		Order("id desc").
		Limit(int(limit)).
		Find(&result).Error
	return result, err
}

func (m *ReplyModel) CountByBizIDAndTargetID(ctx context.Context, bizId string, targetId int64) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&Reply{}).
		Where("biz_id = ? AND target_id = ? AND status = ?", bizId, targetId, 0).
		Count(&count).Error
	return count, err
}

func (m *ReplyModel) CountRootByBizIDAndTargetID(ctx context.Context, bizId string, targetId int64) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&Reply{}).
		Where("biz_id = ? AND target_id = ? AND parent_id = 0 AND status = ?", bizId, targetId, 0).
		Count(&count).Error
	return count, err
}
