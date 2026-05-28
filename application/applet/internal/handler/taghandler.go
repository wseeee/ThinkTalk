package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/applet/internal/logic"
	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
)

func TagCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.TagCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.CreateTag(uid, &req)
		writeJSON(w, resp, err)
	}
}

func TagUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.TagUpdateRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.UpdateTag(uid, &req)
		writeJSON(w, resp, err)
	}
}

func TagDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.TagDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.DeleteTag(uid, &req)
		writeJSON(w, resp, err)
	}
}

func TagResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.TagResourceRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.TagResource(uid, &req)
		writeJSON(w, resp, err)
	}
}

func UntagResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.UntagResourceRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.UntagResource(uid, &req)
		writeJSON(w, resp, err)
	}
}

func TagDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagDetailRequest
		if v := r.URL.Query().Get("tag_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TagId)
		}
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.TagDetail(&req)
		writeJSON(w, resp, err)
	}
}

func TagListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagListRequest
		parseQuery(r, &req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.TagList(&req)
		writeJSON(w, resp, err)
	}
}

func HotTagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HotTagsRequest
		if v := r.URL.Query().Get("limit"); v != "" {
			json.Unmarshal([]byte(v), &req.Limit)
		}
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.HotTags(&req)
		writeJSON(w, resp, err)
	}
}

func ResourcesByTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResourcesByTagRequest
		q := r.URL.Query()
		if v := q.Get("tag_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TagId)
		}
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.ResourcesByTag(&req)
		writeJSON(w, resp, err)
	}
}

func TagsByResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagsByResourceRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("target_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TargetId)
		}
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.TagsByResource(&req)
		writeJSON(w, resp, err)
	}
}
