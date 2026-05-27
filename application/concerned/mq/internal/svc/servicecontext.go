package svc

import (
	"ThinkTalk/application/concerned/mq/internal/config"
	"ThinkTalk/application/concerned/mq/internal/model"
	"ThinkTalk/pkg/orm"
)

type ServiceContext struct {
	Config               config.Config
	DB                   *orm.DB
	ConcernedRecordModel *model.ConcernedRecordModel
	ConcernedCountModel  *model.ConcernedCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN: c.Mysql.DataSource,
	})

	return &ServiceContext{
		Config:               c,
		DB:                   db,
		ConcernedRecordModel: model.NewConcernedRecordModel(db.DB),
		ConcernedCountModel:  model.NewConcernedCountModel(db.DB),
	}
}
