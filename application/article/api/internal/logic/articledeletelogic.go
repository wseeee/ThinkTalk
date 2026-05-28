package logic

import (
	"context"
	"encoding/json"

	"ThinkTalk/application/article/api/internal/svc"
	"ThinkTalk/application/article/api/internal/types"
	"ThinkTalk/application/article/rpc/article"
	"ThinkTalk/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
	return &ArticleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDeleteLogic) ArticleDelete(req *types.ArticleDeleteRequest) error {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value error: %v", err)
		return xcode.NoLogin
	}
	_, err = l.svcCtx.ArticleRPC.ArticleDelete(l.ctx, &article.ArticleDeleteRequest{
		UserId:    userId,
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("ArticleDelete req: %v userId: %d error: %v", req, userId, err)
		return err
	}
	return nil
}
