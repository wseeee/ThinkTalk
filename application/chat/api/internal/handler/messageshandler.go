package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/chat/api/internal/logic"
	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/application/chat/api/internal/types"
)

func MessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(json.Number)
		uid, err := userId.Int64()
		if err != nil || uid == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req types.MessagesRequest
		if v := r.URL.Query().Get("conversation_id"); v != "" {
			_ = json.Unmarshal([]byte(v), &req.ConversationId)
		}
		if v := r.URL.Query().Get("cursor"); v != "" {
			_ = json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := r.URL.Query().Get("page_size"); v != "" {
			_ = json.Unmarshal([]byte(v), &req.PageSize)
		}

		l := logic.NewMessagesLogic(r.Context(), svcCtx)
		resp, err := l.Messages(uid, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
