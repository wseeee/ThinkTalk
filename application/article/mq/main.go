package main

import (
		"ThinkTalk/pkg/env"
"context"
	"flag"

	"ThinkTalk/application/article/mq/internal/config"
	"ThinkTalk/application/article/mq/internal/logic"
	"ThinkTalk/application/article/mq/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/article.yaml", "the config file")

func main() {
	flag.Parse()

	env.LoadEnv()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	err := c.ServiceConf.SetUp()
	if err != nil {
		panic(err)
	}

	logx.DisableStat()
	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range logic.Consumers(ctx, svcCtx) {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}
