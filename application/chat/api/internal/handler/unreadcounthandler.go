package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/chat/api/internal/logic"
	"ThinkTalk/application/chat/api/internal/svc"
)

func UnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(json.Number)
		uid, err := userId.Int64()
		if err != nil || uid == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		l := logic.NewUnreadCountLogic(r.Context(), svcCtx)
		resp, err := l.UnreadCount(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
