package svc

import (
	"ThinkTalk/application/article/rpc/internal/config"
	"ThinkTalk/application/article/rpc/internal/model"
	"ThinkTalk/pkg/es"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config            config.Config
	ArticleMOdel      model.ArticleModel
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
	Es                *es.Es
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:       c,
		ArticleMOdel: model.NewArticleModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:     rds,
		Es: es.MustNewEs(&es.Config{
			Addresses: c.Es.Addresses,
			Username:  c.Es.Username,
			Password:  c.Es.Password,
		}),
	}
}
