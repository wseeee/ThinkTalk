package svc

import (
	"ThinkTalk/application/chat/api/internal/config"
	"ThinkTalk/application/chat/api/internal/hub"
	"ThinkTalk/application/chat/rpc/chat"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Chat   chat.Chat
	Redis  *redis.Redis
	Hub    *hub.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})

	return &ServiceContext{
		Config: c,
		Chat:   chat.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		Redis:  rds,
		Hub:    hub.NewHub(),
	}
}
