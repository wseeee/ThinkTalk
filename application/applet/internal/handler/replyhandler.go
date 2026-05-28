package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/applet/internal/logic"
	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
)

func ReplyCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.ReplyCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.CreateReply(uid, &req)
		writeJSON(w, resp, err)
	}
}

func ReplyDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.ReplyDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.DeleteReply(uid, &req)
		writeJSON(w, resp, err)
	}
}

func ReplyDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReplyDetailRequest
		if v := r.URL.Query().Get("reply_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ReplyId)
		}
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.ReplyDetail(&req)
		writeJSON(w, resp, err)
	}
}

func ReplyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReplyListRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("target_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TargetId)
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		if v := q.Get("sort_type"); v != "" {
			json.Unmarshal([]byte(v), &req.SortType)
		}
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.ReplyList(&req)
		writeJSON(w, resp, err)
	}
}

func ReplyCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReplyCountRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("target_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TargetId)
		}
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.ReplyCount(&req)
		writeJSON(w, resp, err)
	}
}
