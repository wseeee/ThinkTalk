package logic

import (
	"context"

	"ThinkTalk/application/article/api/internal/svc"
	"ThinkTalk/application/article/api/internal/types"
	"ThinkTalk/application/article/rpc/article"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	ret, err := l.svcCtx.ArticleRPC.SearchArticles(l.ctx, &article.SearchRequest{
		Keyword:  req.Keyword,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("SearchArticles req: %v error: %v", req, err)
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

	return &types.SearchResponse{
		Articles: items,
		Cursor:   ret.Cursor,
		IsEnd:    ret.IsEnd,
	}, nil
}
