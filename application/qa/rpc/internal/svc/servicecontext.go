package svc

import (
	"ThinkTalk/application/qa/rpc/internal/config"
	"ThinkTalk/application/qa/rpc/internal/model"
	"ThinkTalk/pkg/es"
	"ThinkTalk/pkg/orm"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config            config.Config
	DB                *orm.DB
	QuestionModel     *model.QuestionModel
	AnswerModel       *model.AnswerModel
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
	Es                *es.Es
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
		Config:        c,
		DB:            db,
		QuestionModel: model.NewQuestionModel(db.DB),
		AnswerModel:   model.NewAnswerModel(db.DB),
		BizRedis:      rds,
		Es: es.MustNewEs(&es.Config{
			Addresses: c.Es.Addresses,
			Username:  c.Es.Username,
			Password:  c.Es.Password,
		}),
	}
}
