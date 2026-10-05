package handler

import (
	"encoding/json"
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

func UpgradeMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.UpgradeMemberRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewMemberLogic(r.Context(), svcCtx)
		resp, err := l.UpgradeMember(uid, &req)
		writeJSON(w, resp, err)
	}
}

func MemberOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.MemberOrderListRequest
		parseQuery(r, &req)
		l := logic.NewMemberLogic(r.Context(), svcCtx)
		resp, err := l.MemberOrderList(uid, &req)
		writeJSON(w, resp, err)
	}
}
