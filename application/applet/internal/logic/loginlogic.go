// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"ThinkTalk/application/applet/internal/code"
	"ThinkTalk/application/user/rpc/user"
	"ThinkTalk/pkg/encrypt"
	"ThinkTalk/pkg/jwt"
	"ThinkTalk/pkg/xcode"
	"context"
	"strings"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.LoginMobileEmpty
	}
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	req.Password = strings.TrimSpace(req.Password)

	var isPassMode bool
	if len(req.VerificationCode) > 0 {
		isPassMode = false
	} else if len(req.Password) > 0 {
		isPassMode = true
	} else {
		return nil, code.VerificationCodeEmpty
	}

	encMobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("encMobile error: %v", err)
		return nil, err
	}
	mobile, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{Mobile: encMobile})
	if err != nil {
		logx.Errorf("findByMobile error: %v", err)
		return nil, err
	}
	if mobile == nil || mobile.UserId == 0 {
		return nil, xcode.AccessDenied
	}

	if isPassMode {
		if encrypt.MD5Password(req.Password) != mobile.Password {
			return nil, xcode.AccessDenied
		}
	} else {
		if err := CheckVerificationCode(l.svcCtx.RDB, req.Mobile, req.VerificationCode); err != nil {
			logx.Errorf("checkVerificationCode error: %v", err)
			return nil, err
		}
		_ = deleteActivationCode(req.Mobile, req.VerificationCode, l.svcCtx.RDB)
	}

	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": mobile.UserId,
		},
	})
	if err != nil {
		logx.Errorf("buildTokens error: %v", err)
		return nil, err
	}
	return &types.LoginResponse{
		UserId: mobile.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}
