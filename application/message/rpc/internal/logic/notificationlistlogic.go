package logic

import (
	"context"
	"math"

	"ThinkTalk/application/message/code"
	"ThinkTalk/application/message/rpc/internal/svc"
	"ThinkTalk/application/message/rpc/internal/types"
	"ThinkTalk/application/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationListLogic {
	return &NotificationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotificationListLogic) NotificationList(in *pb.NotificationListRequest) (*pb.NotificationListResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.PageSize > types.MaxPageSize {
		in.PageSize = types.MaxPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = math.MaxInt64
	}

	notifs, err := l.svcCtx.NotificationModel.FindByUserId(l.ctx, in.UserId, in.Type, in.Cursor, in.PageSize+1)
	if err != nil {
		l.Errorf("[NotificationList] FindByUserId err: %v userId: %d", err, in.UserId)
		return nil, err
	}

	var isEnd bool
	if len(notifs) > int(in.PageSize) {
		notifs = notifs[:in.PageSize]
	} else {
		isEnd = true
	}
	if len(notifs) == 0 {
		return &pb.NotificationListResponse{IsEnd: true}, nil
	}

	items := make([]*pb.NotificationItem, 0, len(notifs))
	for _, n := range notifs {
		items = append(items, &pb.NotificationItem{
			Id:            n.ID,
			Type:          n.Type,
			Title:         n.Title,
			Content:       n.Content,
			RefId:         n.RefID,
			BizId:         n.BizID,
			TriggerUserId: n.TriggerUserID,
			IsRead:        n.IsRead == 1,
			CreateTime:    n.CreateTime.Unix(),
		})
	}

	cursor := notifs[len(notifs)-1].ID
	return &pb.NotificationListResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}
