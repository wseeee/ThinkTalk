// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"ThinkTalk/pkg/xcode"
	"flag"
	"fmt"

	"ThinkTalk/application/applet/internal/config"
	"ThinkTalk/application/applet/internal/handler"
	"ThinkTalk/application/applet/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/applet-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(xcode.ErrHandler)
	
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
