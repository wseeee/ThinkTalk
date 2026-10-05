package logic

import (
	"context"

	"ThinkTalk/application/like/mq/internal/model"
	"ThinkTalk/application/like/mq/internal/types"
)

func (l *ThumbupLogic) addLike(ctx context.Context, msg types.ThumbupMsg) error {
	_, err := l.svcCtx.LikeRecordModel.Insert(ctx, &model.LikeRecord{
		BizId:    msg.BizId,
		ObjId:    msg.ObjId,
		UserId:   msg.UserId,
		LikeType: int64(msg.LikeType),
	})
	if err != nil {
		l.Errorf("[Thumbup] insert like record error: %v", err)
		return err
	}

	count, err := l.getOrCreateCount(ctx, msg.BizId, msg.ObjId)
	if err != nil {
		return err
	}

	if msg.LikeType == 0 {
		count.LikeNum++
	} else {
		count.DislikeNum++
	}
	if count.Id == 0 {
		_, err = l.svcCtx.LikeCountModel.Insert(ctx, count)
	} else {
		err = l.svcCtx.LikeCountModel.Update(ctx, count)
	}
	if err != nil {
		l.Errorf("[Thumbup] save like count error: %v", err)
		return err
	}
	return nil
}
