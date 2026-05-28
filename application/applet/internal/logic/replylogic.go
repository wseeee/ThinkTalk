package logic

import (
	"context"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/reply/rpc/reply"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyLogic {
	return &ReplyLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyLogic) CreateReply(userId int64, req *types.ReplyCreateRequest) (*types.ReplyCreateResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.CreateReply(l.ctx, &reply.CreateReplyRequest{
		BizId:         req.BizId,
		TargetId:      req.TargetId,
		ReplyUserId:   userId,
		BeReplyUserId: req.BeReplyUserId,
		ParentId:      req.ParentId,
		Content:       req.Content,
	})
	if err != nil {
		l.Errorf("[CreateReply] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyCreateResponse{ReplyId: resp.ReplyId}, nil
}

func (l *ReplyLogic) DeleteReply(userId int64, req *types.ReplyDeleteRequest) (*types.ReplyDeleteResponse, error) {
	_, err := l.svcCtx.ReplyRPC.DeleteReply(l.ctx, &reply.DeleteReplyRequest{
		ReplyId: req.ReplyId,
		UserId:  userId,
	})
	if err != nil {
		l.Errorf("[DeleteReply] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyDeleteResponse{}, nil
}

func (l *ReplyLogic) ReplyDetail(req *types.ReplyDetailRequest) (*types.ReplyDetailResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyDetail(l.ctx, &reply.ReplyDetailRequest{
		ReplyId: req.ReplyId,
	})
	if err != nil {
		l.Errorf("[ReplyDetail] rpc err: %v", err)
		return nil, err
	}
	if resp == nil || resp.Reply == nil {
		return &types.ReplyDetailResponse{}, nil
	}
	return &types.ReplyDetailResponse{Reply: convertReplyItem(resp.Reply)}, nil
}

func (l *ReplyLogic) ReplyList(req *types.ReplyListRequest) (*types.ReplyListResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyList(l.ctx, &reply.ReplyListRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
	})
	if err != nil {
		l.Errorf("[ReplyList] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.ReplyItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, convertReplyItem(item))
	}
	return &types.ReplyListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *ReplyLogic) ReplyCount(req *types.ReplyCountRequest) (*types.ReplyCountResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyCount(l.ctx, &reply.ReplyCountRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
	})
	if err != nil {
		l.Errorf("[ReplyCount] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyCountResponse{
		ReplyNum:     resp.ReplyNum,
		ReplyRootNum: resp.ReplyRootNum,
	}, nil
}

func convertReplyItem(pb *reply.ReplyItem) *types.ReplyItem {
	if pb == nil {
		return nil
	}
	item := &types.ReplyItem{
		ReplyId:       pb.ReplyId,
		BizId:         pb.BizId,
		TargetId:      pb.TargetId,
		ReplyUserId:   pb.ReplyUserId,
		BeReplyUserId: pb.BeReplyUserId,
		ParentId:      pb.ParentId,
		Content:       pb.Content,
		LikeNum:       pb.LikeNum,
		CreateTime:    pb.CreateTime,
	}
	for _, sub := range pb.SubReplies {
		item.SubReplies = append(item.SubReplies, convertReplyItem(sub))
	}
	return item
}
