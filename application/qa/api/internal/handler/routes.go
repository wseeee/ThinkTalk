package handler

import (
	"net/http"

	"ThinkTalk/application/qa/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{Method: http.MethodPost, Path: "/publish", Handler: PublishQuestionHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/answer", Handler: AnswerQuestionHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/accept", Handler: AcceptAnswerHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/list", Handler: QuestionsHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/detail", Handler: QuestionDetailHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/delete", Handler: QuestionDeleteHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/answers", Handler: AnswerListHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/answer/delete", Handler: AnswerDeleteHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/search", Handler: SearchQuestionsHandler(serverCtx)},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/v1/qa"),
	)
}
