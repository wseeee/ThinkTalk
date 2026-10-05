package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ThinkTalk/application/article/rpc/internal/code"
	"ThinkTalk/application/article/rpc/internal/svc"
	"ThinkTalk/application/article/rpc/internal/types"
	"ThinkTalk/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticlesLogic) Articles(in *pb.ArticlesRequest) (*pb.ArticlesResponse, error) {
	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		return nil, code.SortTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}

	sortField := "publish_time"
	if in.SortType == types.SortLikeCount {
		sortField = "like_num"
	}

	query := buildArticlesByUserQuery(in.UserId, sortField, in.PageSize+1, in.Cursor, in.ArticleId)

	res, err := l.svcCtx.Es.Search(
		l.svcCtx.Es.Search.WithContext(l.ctx),
		l.svcCtx.Es.Search.WithIndex("article-index"),
		l.svcCtx.Es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		l.Errorf("[Articles] ES search err: %v userId: %d", err, in.UserId)
		return nil, err
	}
	defer res.Body.Close()

	var result esListResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		l.Errorf("[Articles] decode err: %v", err)
		return nil, err
	}

	hits := result.Hits.Hits
	var isEnd bool
	if len(hits) > int(in.PageSize) {
		hits = hits[:in.PageSize]
	} else {
		isEnd = true
	}

	items := make([]*pb.ArticleItem, 0, len(hits))
	for _, hit := range hits {
		src := hit.Source
		t, err := time.ParseInLocation("2006-01-02 15:04:05", src.PublishTime, time.Local)
		var publishTimeUnix int64
		if err == nil {
			publishTimeUnix = t.Unix()
		}
		items = append(items, &pb.ArticleItem{
			Id:           src.ArticleID,
			Title:        src.Title,
			Content:      src.Content,
			Description:  src.Description,
			Cover:        src.Cover,
			AuthorId:     src.AuthorID,
			LikeCount:    src.LikeNum,
			CommentCount: src.CommentNum,
			PublishTime:  publishTimeUnix,
		})
	}

	var cursor, lastId int64
	if len(hits) > 0 {
		last := hits[len(hits)-1]
		lastId = last.Source.ArticleID
		if in.SortType == types.SortLikeCount {
			cursor = last.Source.LikeNum
		} else {
			t, err := time.ParseInLocation("2006-01-02 15:04:05", last.Source.PublishTime, time.Local)
			if err == nil {
				cursor = t.Unix()
			}
		}
		if cursor < 0 {
			cursor = 0
		}
	}

	return &pb.ArticlesResponse{
		Articles:  items,
		IsEnd:     isEnd,
		Cursor:    cursor,
		ArticleId: lastId,
	}, nil
}

func buildArticlesByUserQuery(userId int64, sortField string, size, cursor, articleId int64) string {
	q := fmt.Sprintf(
		`{"query":{"bool":{"must":[{"term":{"author_id":%d}},{"term":{"status":2}}]}},"size":%d,"sort":[{"%s":"desc"},{"article_id":"asc"}]`,
		userId, size, sortField,
	)
	if cursor > 0 && articleId > 0 {
		if sortField == "publish_time" {
			cursorStr := time.Unix(cursor, 0).Local().Format("2006-01-02 15:04:05")
			q += fmt.Sprintf(`,"search_after":["%s",%d]`, cursorStr, articleId)
		} else {
			q += fmt.Sprintf(`,"search_after":[%d,%d]`, cursor, articleId)
		}
	}
	q += "}"
	return q
}

type esListResult struct {
	Hits struct {
		Hits []struct {
			Source struct {
				ArticleID   int64  `json:"article_id"`
				Title       string `json:"title"`
				Content     string `json:"content"`
				Description string `json:"description"`
				Cover       string `json:"cover"`
				AuthorID    int64  `json:"author_id"`
				LikeNum     int64  `json:"like_num"`
				CommentNum  int64  `json:"comment_num"`
				PublishTime string `json:"publish_time"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func articlesKey(uid int64, sortType int32) string {
	return fmt.Sprintf("biz#articles#%d#%d", uid, sortType)
}
