// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"ThinkTalk/application/article/api/internal/svc"
	"ThinkTalk/application/article/api/internal/types"
	"ThinkTalk/application/article/rpc/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllArticlesLogic {
	return &AllArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllArticlesLogic) AllArticles(req *types.AllArticlesRequest) (resp *types.AllArticlesResponse, err error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	ret, err := l.svcCtx.ArticleRPC.SearchArticles(l.ctx, &article.SearchRequest{
		Keyword:  "", // Empty keyword gets all articles sorted by time
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("AllArticles SearchArticles error: %v", err)
		return nil, err
	}

	items := make([]types.SearchInfo, 0, len(ret.Items))
	for _, item := range ret.Items {
		items = append(items, types.SearchInfo{
			ArticleId:   item.ArticleId,
			Title:       item.Title,
			Description: item.Description,
			Cover:       item.Cover,
			AuthorId:    item.AuthorId,
			AuthorName:  item.AuthorName,
			LikeNum:     item.LikeNum,
			CommentNum:  item.CommentNum,
			PublishTime: item.PublishTime,
		})
	}

	return &types.AllArticlesResponse{
		Articles: items,
		Cursor:   ret.Cursor,
		IsEnd:    ret.IsEnd,
	}, nil
}
