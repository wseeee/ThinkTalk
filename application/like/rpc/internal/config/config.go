package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	Mysql      sqlx.SqlConf
	CacheRedis cache.CacheConf
}
