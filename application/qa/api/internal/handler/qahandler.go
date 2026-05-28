package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/qa/api/internal/logic"
	"ThinkTalk/application/qa/api/internal/svc"
	"ThinkTalk/application/qa/api/internal/types"
)

func getUserID(r *http.Request) int64 {
	userId, _ := r.Context().Value("userId").(json.Number)
	uid, _ := userId.Int64()
	return uid
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

func PublishQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.PublishQuestionRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.PublishQuestion(uid, &req)
		writeJSON(w, resp, err)
	}
}

func AnswerQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.AnswerQuestionRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.AnswerQuestion(uid, &req)
		writeJSON(w, resp, err)
	}
}

func AcceptAnswerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.AcceptAnswerRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		err := l.AcceptAnswer(uid, &req)
		writeJSON(w, nil, err)
	}
}

func QuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.QuestionListRequest
		q := r.URL.Query()
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		if v := q.Get("sort_type"); v != "" {
			json.Unmarshal([]byte(v), &req.SortType)
		}
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.Questions(uid, &req)
		writeJSON(w, resp, err)
	}
}

func QuestionDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuestionDetailRequest
		if v := r.URL.Query().Get("question_id"); v != "" {
			json.Unmarshal([]byte(v), &req.QuestionId)
		}
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.QuestionDetail(&req)
		writeJSON(w, resp, err)
	}
}

func QuestionDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.QuestionDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		err := l.QuestionDelete(uid, &req)
		writeJSON(w, nil, err)
	}
}

func AnswerListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AnswerListRequest
		q := r.URL.Query()
		if v := q.Get("question_id"); v != "" {
			json.Unmarshal([]byte(v), &req.QuestionId)
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.AnswerList(&req)
		writeJSON(w, resp, err)
	}
}

func AnswerDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.AnswerDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		err := l.AnswerDelete(uid, &req)
		writeJSON(w, nil, err)
	}
}

func SearchQuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchQuestionsRequest
		q := r.URL.Query()
		if v := q.Get("keyword"); v != "" {
			req.Keyword = v
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.SearchQuestions(&req)
		writeJSON(w, resp, err)
	}
}
