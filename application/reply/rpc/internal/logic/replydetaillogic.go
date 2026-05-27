package logic

import (
	"context"

	"ThinkTalk/application/reply/code"
	"ThinkTalk/application/reply/rpc/internal/svc"
	"ThinkTalk/application/reply/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyDetailLogic {
	return &ReplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyDetailLogic) ReplyDetail(in *pb.ReplyDetailRequest) (*pb.ReplyDetailResponse, error) {
	if in.ReplyId == 0 {
		return nil, code.ReplyNotFound
	}

	reply, err := l.svcCtx.ReplyModel.FindOne(l.ctx, in.ReplyId)
	if err != nil {
		l.Errorf("[ReplyDetail] ReplyModel.FindOne err: %v replyId: %d", err, in.ReplyId)
		return nil, err
	}
	if reply == nil || reply.Status == 1 {
		return nil, code.ReplyNotFound
	}

	item := &pb.ReplyItem{
		ReplyId:       reply.ID,
		BizId:         reply.BizID,
		TargetId:      reply.TargetID,
		ReplyUserId:   reply.ReplyUserID,
		BeReplyUserId: reply.BeReplyUserID,
		ParentId:      reply.ParentID,
		Content:       reply.Content,
		LikeNum:       int64(reply.LikeNum),
		CreateTime:    reply.CreateTime.Unix(),
	}

	return &pb.ReplyDetailResponse{Reply: item}, nil
}
