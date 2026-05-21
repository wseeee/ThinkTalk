// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"ThinkTalk/application/article/api/internal/config"
	"ThinkTalk/application/article/rpc/article"
	"ThinkTalk/application/user/rpc/user"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRPC    user.User
	ArticleRPC article.Article
	MinIO      *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {

	minioClient, err := minio.New(c.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MinIO.AccessKeyID, c.MinIO.AccessKeySecret, ""),
		Secure: c.MinIO.UseSSL,
	})
	if err != nil {
		panic("minio 连接失败: " + err.Error())
	}

	return &ServiceContext{
		Config:     c,
		UserRPC:    user.NewUser(zrpc.MustNewClient(c.UserRPC)),
		ArticleRPC: article.NewArticle(zrpc.MustNewClient(c.ArticleRPC)),
		MinIO:      minioClient,
	}
}
