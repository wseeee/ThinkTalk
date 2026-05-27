package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ThinkTalk/application/article/rpc/internal/svc"
	"ThinkTalk/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchArticlesLogic {
	return &SearchArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchArticlesLogic) SearchArticles(in *pb.SearchRequest) (*pb.SearchResponse, error) {
	if in.PageSize == 0 {
		in.PageSize = 20
	}

	query := buildSearchQuery(in.Keyword, in.PageSize+1, in.Cursor)

	res, err := l.svcCtx.Es.Search(
		l.svcCtx.Es.Search.WithContext(l.ctx),
		l.svcCtx.Es.Search.WithIndex("article-index"),
		l.svcCtx.Es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		l.Errorf("[SearchArticles] ES search err: %v keyword: %s", err, in.Keyword)
		return nil, err
	}
	defer res.Body.Close()

	var result esSearchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		l.Errorf("[SearchArticles] decode err: %v", err)
		return nil, err
	}

	hits := result.Hits.Hits
	var isEnd bool
	if len(hits) > int(in.PageSize) {
		hits = hits[:in.PageSize]
	} else {
		isEnd = true
	}

	items := make([]*pb.SearchItem, 0, len(hits))
	for _, hit := range hits {
		src := hit.Source
		items = append(items, &pb.SearchItem{
			ArticleId:   src.ArticleID,
			Title:       src.Title,
			Description: src.Description,
			Cover:       src.Cover,
			AuthorId:    src.AuthorID,
			AuthorName:  src.AuthorName,
			LikeNum:     src.LikeNum,
			CommentNum:  src.CommentNum,
			PublishTime: src.PublishTime,
		})
	}

	var cursor int64
	if len(hits) > 0 && !isEnd {
		cursor = hits[len(hits)-1].Sort[0]
	}

	return &pb.SearchResponse{
		Items:  items,
		Cursor: cursor,
		IsEnd:  isEnd,
	}, nil
}

func buildSearchQuery(keyword string, size int64, cursor int64) string {
	q := fmt.Sprintf(
		`{"query":{"bool":{"must":[{"multi_match":{"query":"%s","fields":["title^3","content","description"]}},{"term":{"status":2}}]}},"size":%d`,
		keyword, size,
	)
	if cursor > 0 {
		q += fmt.Sprintf(`,"search_after":[%d]`, cursor)
	}
	q += `,"sort":[{"_score":"desc"},{"article_id":"asc"}]}`
	return q
}

type esSearchResult struct {
	Hits struct {
		Hits []struct {
			Source struct {
				ArticleID   int64  `json:"article_id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Cover       string `json:"cover"`
				AuthorID    int64  `json:"author_id"`
				AuthorName  string `json:"author_name"`
				LikeNum     int64  `json:"like_num"`
				CommentNum  int64  `json:"comment_num"`
				PublishTime string `json:"publish_time"`
			} `json:"_source"`
			Sort []int64 `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}
