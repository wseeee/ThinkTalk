// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"ThinkTalk/application/article/rpc/article"
	"context"

	"ThinkTalk/application/article/api/internal/svc"
	"ThinkTalk/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleListLogic {
	return &ArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleListLogic) ArticleList(req *types.ArticleListRequest) (resp *types.ArticleListResponse, err error) {
	articles, err := l.svcCtx.ArticleRPC.Articles(l.ctx, &article.ArticlesRequest{
		UserId:    req.AuthorId,
		Cursor:    req.Cursor,
		PageSize:  req.PageSize,
		SortType:  req.SortType,
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("get articles req: %v err: %v", req, err)
		return nil, err
	}
	if articles == nil || len(articles.Articles) == 0 {
		return &types.ArticleListResponse{}, nil
	}

	infos := make([]types.ArticleInfo, 0, len(articles.Articles))
	for _, article := range articles.Articles {
		infos = append(infos, types.ArticleInfo{
			ArticleId:   article.Id,
			Cover:       article.Cover,
			Description: article.Description,
			Title:       article.Title,
		})
	}

	return &types.ArticleListResponse{
		Articles: infos,
	}, nil
}
