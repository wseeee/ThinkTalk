package svc

import (
	"ThinkTalk/application/reply/mq/internal/config"
	"ThinkTalk/application/reply/mq/internal/model"
	"ThinkTalk/pkg/orm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *orm.DB
	ReplyModel      *model.ReplyModel
	ReplyCountModel *model.ReplyCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN: c.Mysql.DataSource,
	})

	return &ServiceContext{
		Config:          c,
		DB:              db,
		ReplyModel:      model.NewReplyModel(db.DB),
		ReplyCountModel: model.NewReplyCountModel(db.DB),
	}
}
