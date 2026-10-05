package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ThinkTalk/application/article/mq/internal/svc"
	"ThinkTalk/application/article/mq/internal/types"
	"ThinkTalk/application/user/rpc/user"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleLogic) Consume(ctx context.Context, _, val string) error {
	logx.Infof("Consume msg val: %s", val)
	var msg *types.CanalArticleMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", val, err)
		return err
	}

	return l.articleOperate(msg)
}

func (l *ArticleLogic) articleOperate(msg *types.CanalArticleMsg) error {
	if len(msg.Data) == 0 {
		return nil
	}

	var esUpsertData []*types.ArticleEsMsg
	var esDeleteIds []string

	for _, d := range msg.Data {
		status, _ := strconv.Atoi(d.Status)
		likNum, _ := strconv.ParseInt(d.LikeNum, 10, 64)
		commentNum, _ := strconv.ParseInt(d.CommentNum, 10, 64)
		collectNum, _ := strconv.ParseInt(d.CollectNum, 10, 64)
		viewNum, _ := strconv.ParseInt(d.ViewNum, 10, 64)
		shareNum, _ := strconv.ParseInt(d.ShareNum, 10, 64)
		articleId, _ := strconv.ParseInt(d.ID, 10, 64)
		authorId, _ := strconv.ParseInt(d.AuthorId, 10, 64)

		t, err := time.ParseInLocation("2006-01-02 15:04:05", d.PublishTime, time.Local)
		if err != nil {
			t = time.Now()
		}
		publishTimeKey := articlesKey(d.AuthorId, 0)
		likeNumKey := articlesKey(d.AuthorId, 1)

		switch status {
		case types.ArticleStatusVisible:
			b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, publishTimeKey)
			if b {
				_, _ = l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, t.Unix(), d.ID)
			}
			b, _ = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
			if b {
				_, _ = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, likNum, d.ID)
			}

			u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: authorId})
			if err != nil {
				l.Logger.Errorf("FindById userId: %d error: %v", authorId, err)
				return err
			}

			esUpsertData = append(esUpsertData, &types.ArticleEsMsg{
				ArticleId:   articleId,
				AuthorId:    authorId,
				AuthorName:  u.Username,
				Title:       d.Title,
				Content:     d.Content,
				Description: d.Description,
				Cover:       d.Cover,
				Status:      status,
				LikeNum:     likNum,
				CommentNum:  commentNum,
				CollectNum:  collectNum,
				ViewNum:     viewNum,
				ShareNum:    shareNum,
				PublishTime: d.PublishTime,
				CreateTime:  d.CreateTime,
				UpdateTime:  d.UpdateTime,
			})

		case types.ArticleStatusUserDelete:
			_, _ = l.svcCtx.BizRedis.ZremCtx(l.ctx, publishTimeKey, d.ID)
			_, _ = l.svcCtx.BizRedis.ZremCtx(l.ctx, likeNumKey, d.ID)
			esDeleteIds = append(esDeleteIds, d.ID)
		}
	}

	if len(esUpsertData) > 0 {
		if err := l.BatchUpsertToEs(l.ctx, esUpsertData); err != nil {
			l.Logger.Errorf("BatchUpsertToEs data: %v error: %v", esUpsertData, err)
		}
	}
	if len(esDeleteIds) > 0 {
		if err := l.BatchDeleteFromEs(l.ctx, esDeleteIds); err != nil {
			l.Logger.Errorf("BatchDeleteFromEs ids: %v error: %v", esDeleteIds, err)
		}
	}

	return nil
}

func (l *ArticleLogic) BatchUpsertToEs(ctx context.Context, data []*types.ArticleEsMsg) error {
	if len(data) == 0 {
		return nil
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: l.svcCtx.Es.Client,
		Index:  "article-index",
	})
	if err != nil {
		return err
	}

	for _, d := range data {
		v, err := json.Marshal(d)
		if err != nil {
			return err
		}

		payload := fmt.Sprintf(`{"doc":%s,"doc_as_upsert":true}`, string(v))
		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "update",
			DocumentID: fmt.Sprintf("%d", d.ArticleId),
			Body:       strings.NewReader(payload),
		})
		if err != nil {
			return err
		}
	}

	return bi.Close(ctx)
}

func (l *ArticleLogic) BatchDeleteFromEs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: l.svcCtx.Es.Client,
		Index:  "article-index",
	})
	if err != nil {
		return err
	}

	for _, id := range ids {
		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "delete",
			DocumentID: id,
		})
		if err != nil {
			return err
		}
	}

	return bi.Close(ctx)
}

func articlesKey(uid string, sortType int32) string {
	return fmt.Sprintf("biz#articles#%s#%d", uid, sortType)
}
