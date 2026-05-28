package logic

import (
	"context"

	"ThinkTalk/application/qa/api/internal/svc"
	"ThinkTalk/application/qa/api/internal/types"
	"ThinkTalk/application/qa/rpc/qa"

	"github.com/zeromicro/go-zero/core/logx"
)

type QALogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQALogic(ctx context.Context, svcCtx *svc.ServiceContext) *QALogic {
	return &QALogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QALogic) PublishQuestion(userId int64, req *types.PublishQuestionRequest) (*types.PublishQuestionResponse, error) {
	resp, err := l.svcCtx.QaRPC.PublishQuestion(l.ctx, &qa.PublishQuestionRequest{
		UserId:  userId,
		Title:   req.Title,
		Content: req.Content,
		TagIds:  req.TagIds,
	})
	if err != nil {
		l.Errorf("[PublishQuestion] rpc err: %v", err)
		return nil, err
	}
	return &types.PublishQuestionResponse{QuestionId: resp.QuestionId}, nil
}

func (l *QALogic) AnswerQuestion(userId int64, req *types.AnswerQuestionRequest) (*types.AnswerQuestionResponse, error) {
	resp, err := l.svcCtx.QaRPC.AnswerQuestion(l.ctx, &qa.AnswerQuestionRequest{
		QuestionId: req.QuestionId,
		UserId:     userId,
		Content:    req.Content,
	})
	if err != nil {
		l.Errorf("[AnswerQuestion] rpc err: %v", err)
		return nil, err
	}
	return &types.AnswerQuestionResponse{AnswerId: resp.AnswerId}, nil
}

func (l *QALogic) AcceptAnswer(userId int64, req *types.AcceptAnswerRequest) error {
	_, err := l.svcCtx.QaRPC.AcceptAnswer(l.ctx, &qa.AcceptAnswerRequest{
		QuestionId: req.QuestionId,
		AnswerId:   req.AnswerId,
		UserId:     userId,
	})
	if err != nil {
		l.Errorf("[AcceptAnswer] rpc err: %v", err)
		return err
	}
	return nil
}

func (l *QALogic) Questions(userId int64, req *types.QuestionListRequest) (*types.QuestionListResponse, error) {
	resp, err := l.svcCtx.QaRPC.Questions(l.ctx, &qa.QuestionsRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
	})
	if err != nil {
		l.Errorf("[Questions] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.QuestionItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.QuestionItem{
			Id:         item.Id,
			Title:      item.Title,
			Content:    item.Content,
			AuthorId:   item.AuthorId,
			AnswerNum:  item.AnswerNum,
			ViewNum:    item.ViewNum,
			TagIds:     item.TagIds,
			CreateTime: item.CreateTime,
		})
	}
	return &types.QuestionListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *QALogic) QuestionDetail(req *types.QuestionDetailRequest) (*types.QuestionDetailResponse, error) {
	resp, err := l.svcCtx.QaRPC.QuestionDetail(l.ctx, &qa.QuestionDetailRequest{
		QuestionId: req.QuestionId,
	})
	if err != nil {
		l.Errorf("[QuestionDetail] rpc err: %v", err)
		return nil, err
	}
	if resp == nil || resp.Question == nil {
		return &types.QuestionDetailResponse{}, nil
	}
	q := resp.Question
	return &types.QuestionDetailResponse{
		Question: &types.QuestionItem{
			Id:         q.Id,
			Title:      q.Title,
			Content:    q.Content,
			AuthorId:   q.AuthorId,
			AnswerNum:  q.AnswerNum,
			ViewNum:    q.ViewNum,
			TagIds:     q.TagIds,
			CreateTime: q.CreateTime,
		},
	}, nil
}

func (l *QALogic) QuestionDelete(userId int64, req *types.QuestionDeleteRequest) error {
	_, err := l.svcCtx.QaRPC.QuestionDelete(l.ctx, &qa.QuestionDeleteRequest{
		UserId:     userId,
		QuestionId: req.QuestionId,
	})
	if err != nil {
		l.Errorf("[QuestionDelete] rpc err: %v", err)
		return err
	}
	return nil
}

func (l *QALogic) AnswerList(req *types.AnswerListRequest) (*types.AnswerListResponse, error) {
	resp, err := l.svcCtx.QaRPC.AnswerList(l.ctx, &qa.AnswerListRequest{
		QuestionId: req.QuestionId,
		Cursor:     req.Cursor,
		PageSize:   req.PageSize,
	})
	if err != nil {
		l.Errorf("[AnswerList] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.AnswerItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.AnswerItem{
			Id:         item.Id,
			QuestionId: item.QuestionId,
			AuthorId:   item.AuthorId,
			Content:    item.Content,
			IsAccepted: item.IsAccepted,
			LikeNum:    item.LikeNum,
			ReplyNum:   item.ReplyNum,
			CreateTime: item.CreateTime,
		})
	}
	return &types.AnswerListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *QALogic) AnswerDelete(userId int64, req *types.AnswerDeleteRequest) error {
	_, err := l.svcCtx.QaRPC.AnswerDelete(l.ctx, &qa.AnswerDeleteRequest{
		UserId:   userId,
		AnswerId: req.AnswerId,
	})
	if err != nil {
		l.Errorf("[AnswerDelete] rpc err: %v", err)
		return err
	}
	return nil
}

func (l *QALogic) SearchQuestions(req *types.SearchQuestionsRequest) (*types.SearchQuestionsResponse, error) {
	resp, err := l.svcCtx.QaRPC.SearchQuestions(l.ctx, &qa.SearchQuestionsRequest{
		Keyword:  req.Keyword,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[SearchQuestions] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.SearchQuestionItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.SearchQuestionItem{
			Id:         item.Id,
			Title:      item.Title,
			Content:    item.Content,
			AuthorId:   item.AuthorId,
			AnswerNum:  item.AnswerNum,
			TagIds:     item.TagIds,
			CreateTime: item.CreateTime,
		})
	}
	return &types.SearchQuestionsResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}
