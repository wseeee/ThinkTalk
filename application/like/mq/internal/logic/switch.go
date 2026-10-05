package logic

import (
	"context"
	"time"

	"ThinkTalk/application/like/mq/internal/model"
)

func (l *ThumbupLogic) switchLike(ctx context.Context, record *model.LikeRecord, newLikeType int32) error {
	oldLikeType := record.LikeType
	record.LikeType = int64(newLikeType)
	record.UpdateTime = time.Now()

	if err := l.svcCtx.LikeRecordModel.Update(ctx, record); err != nil {
		l.Errorf("[Thumbup] update like record error: %v", err)
		return err
	}

	count, err := l.getOrCreateCount(ctx, record.BizId, record.ObjId)
	if err != nil {
		return err
	}

	if oldLikeType == 0 {
		count.LikeNum--
		count.DislikeNum++
	} else {
		count.DislikeNum--
		count.LikeNum++
	}
	return l.updateCount(ctx, count)
}
