package svc

import (
	"ThinkTalk/application/like/rpc/internal/config"
	"ThinkTalk/application/like/rpc/internal/model"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	KqPusherClient  *kq.Pusher
	LikeRecordModel model.LikeRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:          c,
		KqPusherClient:  kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis),
	}
}
