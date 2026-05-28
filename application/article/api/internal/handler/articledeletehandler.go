package handler

import (
	"net/http"

	"ThinkTalk/application/article/api/internal/logic"
	"ThinkTalk/application/article/api/internal/svc"
	"ThinkTalk/application/article/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ArticleDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticleDeleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewArticleDeleteLogic(r.Context(), svcCtx)
		err := l.ArticleDelete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
