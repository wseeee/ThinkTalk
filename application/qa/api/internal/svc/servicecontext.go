package svc

import (
	"ThinkTalk/application/qa/api/internal/config"
	"ThinkTalk/application/qa/rpc/qa"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	QaRPC  qa.QA
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		QaRPC:  qa.NewQA(zrpc.MustNewClient(c.QaRpc)),
	}
}
