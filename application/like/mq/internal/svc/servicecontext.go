// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"ThinkTalk/application/like/mq/internal/config"
	"ThinkTalk/application/like/mq/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	LikeRecordModel model.LikeRecordModel
	LikeCountModel  model.LikeCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:          c,
		LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis),
		LikeCountModel:  model.NewLikeCountModel(conn, c.CacheRedis),
	}
}
