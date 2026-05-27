package handler

import (
	"net/http"

	"ThinkTalk/application/applet/internal/logic"
	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
)

func MemberInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		l := logic.NewMemberLogic(r.Context(), svcCtx)
		resp, err := l.MemberInfo(uid)
		writeJSON(w, resp, err)
	}
}

func MemberRightHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.MemberRightRequest
		if k := r.URL.Query().Get("right_key"); k != "" {
			req.RightKey = k
		}
		l := logic.NewMemberLogic(r.Context(), svcCtx)
		resp, err := l.CheckRight(uid, &req)
		writeJSON(w, resp, err)
	}
}
