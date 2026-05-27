package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/applet/internal/logic"
	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
)

func ConcernedAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.ConcernedAddRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewConcernedLogic(r.Context(), svcCtx)
		resp, err := l.Add(uid, &req)
		writeJSON(w, resp, err)
	}
}

func ConcernedCancelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.ConcernedCancelRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewConcernedLogic(r.Context(), svcCtx)
		resp, err := l.Cancel(uid, &req)
		writeJSON(w, resp, err)
	}
}

func ConcernedCheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.ConcernedCheckRequest
		if b := r.URL.Query().Get("biz_id"); b != "" {
			req.BizId = b
		}
		if o := r.URL.Query().Get("obj_id"); o != "" {
			json.Unmarshal([]byte(o), &req.ObjId)
		}
		l := logic.NewConcernedLogic(r.Context(), svcCtx)
		resp, err := l.Check(uid, &req)
		writeJSON(w, resp, err)
	}
}

func ConcernedListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.ConcernedListRequest
		if b := r.URL.Query().Get("biz_id"); b != "" {
			req.BizId = b
		}
		parseQuery(r, &req)
		l := logic.NewConcernedLogic(r.Context(), svcCtx)
		resp, err := l.List(uid, &req)
		writeJSON(w, resp, err)
	}
}
