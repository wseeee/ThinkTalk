package logic

import (
	"context"

	"ThinkTalk/application/like/mq/internal/model"
)

func (l *ThumbupLogic) cancelLike(ctx context.Context, record *model.LikeRecord) error {
	if err := l.svcCtx.LikeRecordModel.Delete(ctx, record.Id); err != nil {
		l.Errorf("[Thumbup] delete like record error: %v", err)
		return err
	}

	count, err := l.getOrCreateCount(ctx, record.BizId, record.ObjId)
	if err != nil {
		return err
	}

	if record.LikeType == 0 {
		count.LikeNum--
	} else {
		count.DislikeNum--
	}
	return l.updateCount(ctx, count)
}
