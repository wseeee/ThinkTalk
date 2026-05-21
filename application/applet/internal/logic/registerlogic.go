// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"ThinkTalk/application/applet/internal/code"
	"ThinkTalk/application/user/rpc/user"
	"ThinkTalk/pkg/encrypt"
	"ThinkTalk/pkg/jwt"
	"context"
	"errors"
	"strings"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		return nil, code.RegisterNameEmpty
	}

	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.RegisterMobileEmpty
	}

	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return nil, code.RegisterPasswordEmpty
	} else {
		req.Password = encrypt.MD5Password(req.Password)
	}

	req.VerificationCode = strings.TrimSpace(req.VerificationCode)

	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}

	err = CheckVerificationCode(l.svcCtx.RDB, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}

	mobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("encMobile error: %v", err)
		return nil, err
	}

	Mobile, err := l.svcCtx.UserRPC.
		FindByMobile(l.ctx, &user.FindByMobileRequest{Mobile: mobile})
	if err != nil {
		logx.Errorf("findByMobile error: %v", err)
		return nil, err
	}
	if Mobile != nil && Mobile.UserId > 0 {
		return nil, code.MobileHasRegistered
	}
	userId, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterRequest{
		Mobile:   mobile,
		Username: req.Name,
		Password: req.Password,
	})
	if err != nil {
		logx.Errorf("register error: %v", err)
		return nil, err
	}

	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": userId.UserId,
		},
	})
	if err != nil {
		logx.Errorf("buildTokens error: %v", err)
		return nil, err
	}
	_ = deleteActivationCode(req.Mobile, req.VerificationCode, l.svcCtx.RDB)

	return &types.RegisterResponse{
		UserId: userId.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}

func CheckVerificationCode(rds *redis.Redis, mobile, code string) error {
	activationCode, err := GetActivationCode(mobile, rds)
	if err != nil {
		return err
	}
	if activationCode == "" {
		return errors.New("验证码已过期")
	}
	if activationCode != code {
		return errors.New("验证码错误")
	}
	return nil
}
