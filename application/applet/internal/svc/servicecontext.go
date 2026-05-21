package svc

import (
	"ThinkTalk/application/applet/internal/config"
	like "ThinkTalk/application/like/rpc/like"
	"ThinkTalk/application/user/rpc/user"
	"ThinkTalk/pkg/interceptors"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	RDB     *redis.Redis
	UserRPC user.User
	LikeRPC like.Like
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb, _ := redis.NewRedis(c.Redis)
	userRPC := zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	likeRPC := zrpc.MustNewClient(c.LikeRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:  c,
		RDB:     rdb,
		UserRPC: user.NewUser(userRPC),
		LikeRPC: like.NewLike(likeRPC),
	}
}
