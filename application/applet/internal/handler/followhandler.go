package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/applet/internal/logic"
	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
)

func getUserID(r *http.Request) int64 {
	userId, _ := r.Context().Value("userId").(json.Number)
	uid, _ := userId.Int64()
	return uid
}

func FollowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.FollowRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewFollowLogic(r.Context(), svcCtx)
		resp, err := l.Follow(uid, &req)
		writeJSON(w, resp, err)
	}
}

func UnFollowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.UnfollowRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewFollowLogic(r.Context(), svcCtx)
		resp, err := l.UnFollow(uid, &req)
		writeJSON(w, resp, err)
	}
}

func FollowListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.FollowListRequest
		parseQuery(r, &req)
		l := logic.NewFollowLogic(r.Context(), svcCtx)
		resp, err := l.FollowList(uid, &req)
		writeJSON(w, resp, err)
	}
}

func FansListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.FansListRequest
		parseQuery(r, &req)
		l := logic.NewFollowLogic(r.Context(), svcCtx)
		resp, err := l.FansList(uid, &req)
		writeJSON(w, resp, err)
	}
}

func parseQuery(r *http.Request, v interface{}) {
	if c := r.URL.Query().Get("cursor"); c != "" {
		json.Unmarshal([]byte(c), v)
	}
	if p := r.URL.Query().Get("page_size"); p != "" {
		json.Unmarshal([]byte(p), v)
	}
}

func writeJSON(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(data)
}
