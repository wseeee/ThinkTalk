// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"ThinkTalk/application/user/rpc/user"
	"context"
	"encoding/json"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return
	}

	id, err := l.svcCtx.UserRPC.FindById(l.ctx,
		&user.FindByIdRequest{UserId: userId})
	if err != nil {
		logx.Errorf("findById error: %v", err)
		return nil, err
	}

	return &types.UserInfoResponse{
		UserId:   id.UserId,
		Username: id.Username,
		Avatar:   id.Avatar,
	}, nil
}
