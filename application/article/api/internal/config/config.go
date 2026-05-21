// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	ArticleRPC zrpc.RpcClientConf
	UserRPC    zrpc.RpcClientConf
	MinIO      struct {
		Endpoint        string
		AccessKeyID     string
		AccessKeySecret string
		BucketName      string
		Location        string
		UseSSL          bool
	}
}
