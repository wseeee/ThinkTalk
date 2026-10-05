// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"ThinkTalk/application/user/rpc/user"
	"ThinkTalk/pkg/util"
	"context"
	"fmt"
	"strconv"
	"time"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	prefixVerificationCount = "biz#verification#count#%s"
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
	prefixActivation        = "biz#activation#%s"
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	count, err := l.GetVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobile: %s error: %v", req.Mobile, err)
	}
	if count > verificationLimitPerDay {
		logx.Errorf("mobile: %s verification count over limit", req.Mobile)
		return nil, err
	}
	code, err := GetActivationCode(req.Mobile, l.svcCtx.RDB)
	if err != nil {
		logx.Errorf("getActivationCode mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	if len(code) == 0 {
		code = util.GenerateCode(6)
	}
	_, err = l.svcCtx.UserRPC.SendSms(l.ctx, &user.SendSmsRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		logx.Errorf("sendSms mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	err = SaveActivationCode(req.Mobile, code, l.svcCtx.RDB)
	if err != nil {
		logx.Errorf("saveActivationCode mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	err = l.IncrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerificationCount mobile: %s error: %v", req.Mobile, err)
	}
	return &types.VerificationResponse{}, nil
}

func (l *VerificationLogic) GetVerificationCount(moblie string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, moblie)
	get, err := l.svcCtx.RDB.Get(key)
	if err != nil {
		return 0, err
	}
	if len(get) == 0 {
		return 0, nil
	}
	return strconv.Atoi(get)
}

func (l *VerificationLogic) IncrVerificationCount(mobile string) error {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	_, err := l.svcCtx.RDB.Incr(key)
	if err != nil {
		return err
	}
	expireTime := int(util.EndOfDay(time.Now()).Unix())
	return l.svcCtx.RDB.Expire(key, expireTime)
}

func GetActivationCode(mobile string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Get(key)
}
func SaveActivationCode(mobile string, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Setex(key, code, expireActivation)
}

func deleteActivationCode(mobile string, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	_, err := rds.Del(key)
	return err
}
