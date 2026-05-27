package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/applet/internal/logic"
	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
)

func NotificationListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.NotificationRequest
		if t := r.URL.Query().Get("type"); t != "" {
			json.Unmarshal([]byte(t), &req.NotifType)
		}
		if c := r.URL.Query().Get("cursor"); c != "" {
			json.Unmarshal([]byte(c), &req.Cursor)
		}
		if p := r.URL.Query().Get("page_size"); p != "" {
			json.Unmarshal([]byte(p), &req.PageSize)
		}
		l := logic.NewNotificationLogic(r.Context(), svcCtx)
		resp, err := l.NotificationList(uid, &req)
		writeJSON(w, resp, err)
	}
}

func UnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		l := logic.NewNotificationLogic(r.Context(), svcCtx)
		resp, err := l.UnreadCount(uid)
		writeJSON(w, resp, err)
	}
}

func MarkReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.MarkReadRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewNotificationLogic(r.Context(), svcCtx)
		resp, err := l.MarkRead(uid, &req)
		writeJSON(w, resp, err)
	}
}

func MarkAllReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.MarkAllReadRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewNotificationLogic(r.Context(), svcCtx)
		resp, err := l.MarkAllRead(uid, &req)
		writeJSON(w, resp, err)
	}
}
