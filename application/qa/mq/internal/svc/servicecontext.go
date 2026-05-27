package svc

import (
	"ThinkTalk/application/qa/mq/internal/config"
	"ThinkTalk/application/qa/mq/internal/model"
	"ThinkTalk/pkg/es"
	"ThinkTalk/pkg/orm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *orm.DB
	QuestionModel *model.QuestionModel
	AnswerModel   *model.AnswerModel
	Es            *es.Es
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN: c.Mysql.DataSource,
	})

	return &ServiceContext{
		Config:        c,
		DB:            db,
		QuestionModel: model.NewQuestionModel(db.DB),
		AnswerModel:   model.NewAnswerModel(db.DB),
		Es: es.MustNewEs(&es.Config{
			Addresses: c.Es.Addresses,
			Username:  c.Es.Username,
			Password:  c.Es.Password,
		}),
	}
}
