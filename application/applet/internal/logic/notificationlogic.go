package logic

import (
	"context"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationLogic {
	return &NotificationLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *NotificationLogic) NotificationList(userId int64, req *types.NotificationRequest) (*types.NotificationResponse, error) {
	resp, err := l.svcCtx.MessageRPC.NotificationList(l.ctx, &pb.NotificationListRequest{
		UserId:   userId,
		Type:     req.NotifType,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[NotificationList] rpc err: %v", err)
		return nil, err
	}

	items := make([]*types.NotificationItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.NotificationItem{
			Id:            item.Id,
			Type:          item.Type,
			Title:         item.Title,
			Content:       item.Content,
			IsRead:        item.IsRead,
			TriggerUserId: item.TriggerUserId,
			RefId:         item.RefId,
			CreateTime:    item.CreateTime,
		})
	}
	return &types.NotificationResponse{
		Items:  items,
		Cursor: resp.Cursor,
		IsEnd:  resp.IsEnd,
	}, nil
}

func (l *NotificationLogic) UnreadCount(userId int64) (*types.UnreadCountResponse, error) {
	resp, err := l.svcCtx.MessageRPC.UnreadCount(l.ctx, &pb.UnreadCountRequest{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("[UnreadCount] rpc err: %v", err)
		return nil, err
	}
	return &types.UnreadCountResponse{
		Total:      resp.Total,
		TypeCounts: resp.TypeCounts,
	}, nil
}

func (l *NotificationLogic) MarkRead(userId int64, req *types.MarkReadRequest) (*types.MarkReadResponse, error) {
	_, err := l.svcCtx.MessageRPC.MarkRead(l.ctx, &pb.MarkReadRequest{
		UserId:         userId,
		NotificationId: req.NotificationId,
	})
	if err != nil {
		l.Errorf("[MarkRead] rpc err: %v", err)
		return nil, err
	}
	return &types.MarkReadResponse{}, nil
}

func (l *NotificationLogic) MarkAllRead(userId int64, req *types.MarkAllReadRequest) (*types.MarkReadResponse, error) {
	_, err := l.svcCtx.MessageRPC.MarkAllRead(l.ctx, &pb.MarkAllReadRequest{
		UserId: userId,
		Type:   req.NotifType,
	})
	if err != nil {
		l.Errorf("[MarkAllRead] rpc err: %v", err)
		return nil, err
	}
	return &types.MarkReadResponse{}, nil
}
