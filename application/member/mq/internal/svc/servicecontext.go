package svc

import (
	"ThinkTalk/application/member/mq/internal/config"
	"ThinkTalk/application/member/mq/internal/model"
	"ThinkTalk/pkg/orm"
)

type ServiceContext struct {
	Config           config.Config
	DB               *orm.DB
	MemberModel      *model.MemberModel
	MemberOrderModel *model.MemberOrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN: c.Mysql.DataSource,
	})

	return &ServiceContext{
		Config:           c,
		DB:               db,
		MemberModel:      model.NewMemberModel(db.DB),
		MemberOrderModel: model.NewMemberOrderModel(db.DB),
	}
}
