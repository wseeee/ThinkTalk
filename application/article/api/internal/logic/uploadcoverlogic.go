// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"ThinkTalk/application/article/api/internal/code"
	"context"
	"fmt"
	"net/http"
	"time"

	"ThinkTalk/application/article/api/internal/svc"
	"ThinkTalk/application/article/api/internal/types"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

const maxFileSize = 10 << 20

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(req *http.Request) (resp *types.UploadCoverResponse, err error) {
	err1 := req.ParseMultipartForm(maxFileSize)
	if err1 != nil {
		return nil, code.ParseFormErr
	}

	file, header, err := req.FormFile("cover")
	if err != nil {
		logx.Errorf("get form file failed, err: %v", err)
		return nil, err
	}

	defer file.Close()

	objectKey := genFilename(header.Filename)
	_, err = l.svcCtx.MinIO.PutObject(
		l.ctx,
		l.svcCtx.Config.MinIO.BucketName,
		objectKey,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		logx.Errorf("put object to minio failed, err: %v", err)
		return nil, code.PutBucketErr
	}
	return &types.UploadCoverResponse{
		CoverUrl: genFileURL(
			l.svcCtx.Config.MinIO.Endpoint,
			l.svcCtx.Config.MinIO.BucketName,
			objectKey)}, nil

}
func genFilename(filename string) string {
	return fmt.Sprintf("%d_%s", time.Now().UnixMilli(), filename)
}

func genFileURL(endpoint, bucketName, objectKey string) string {
	return fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, objectKey)
}
