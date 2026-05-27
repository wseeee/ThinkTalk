package svc

import (
	"ThinkTalk/application/reply/rpc/internal/config"
	"ThinkTalk/application/reply/rpc/internal/model"
	"ThinkTalk/pkg/orm"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config          config.Config
	DB              *orm.DB
	ReplyModel      *model.ReplyModel
	ReplyCountModel *model.ReplyCountModel
	KqPusherClient  *kq.Pusher
	BizRedis        *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})

	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})

	return &ServiceContext{
		Config:          c,
		DB:              db,
		ReplyModel:      model.NewReplyModel(db.DB),
		ReplyCountModel: model.NewReplyCountModel(db.DB),
		KqPusherClient:  kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		BizRedis:        rds,
	}
}
