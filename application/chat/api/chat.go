package main

import (
		"ThinkTalk/pkg/env"
"flag"
	"fmt"

	"ThinkTalk/application/chat/api/internal/config"
	"ThinkTalk/application/chat/api/internal/handler"
	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/pkg/xcode"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/chat-api.yaml", "the config file")

func main() {
	flag.Parse()

	env.LoadEnv()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(xcode.ErrHandler)

	fmt.Printf("Starting chat API server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
