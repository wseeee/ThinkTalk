package svc

import (
	"ThinkTalk/application/chat/mq/internal/config"
	"ThinkTalk/application/chat/mq/internal/model"
	"ThinkTalk/pkg/orm"
)

type ServiceContext struct {
	Config            config.Config
	DB                *orm.DB
	ConversationModel *model.ConversationModel
	MessageModel      *model.MessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN: c.Mysql.DataSource,
	})

	return &ServiceContext{
		Config:            c,
		DB:                db,
		ConversationModel: model.NewConversationModel(db.DB),
		MessageModel:      model.NewMessageModel(db.DB),
	}
}
