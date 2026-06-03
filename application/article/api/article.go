// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
		"ThinkTalk/pkg/env"
"flag"
	"fmt"

	"ThinkTalk/application/article/api/internal/config"
	"ThinkTalk/application/article/api/internal/handler"
	"ThinkTalk/application/article/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/article-api.yaml", "the config file")

func main() {
	flag.Parse()

	env.LoadEnv()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
