package types

type PublishQuestionRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	TagIds  string `json:"tag_ids"`
}

type PublishQuestionResponse struct {
	QuestionId int64 `json:"question_id"`
}

type AnswerQuestionRequest struct {
	QuestionId int64  `json:"question_id"`
	Content    string `json:"content"`
}

type AnswerQuestionResponse struct {
	AnswerId int64 `json:"answer_id"`
}

type AcceptAnswerRequest struct {
	QuestionId int64 `json:"question_id"`
	AnswerId   int64 `json:"answer_id"`
}

type QuestionDetailRequest struct {
	QuestionId int64 `form:"question_id"`
}

type QuestionListRequest struct {
	Cursor   int64 `form:"cursor"`
	PageSize int64 `form:"page_size"`
	SortType int32 `form:"sort_type,optional"`
}

type QuestionItem struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"author_id"`
	AnswerNum  int64  `json:"answer_num"`
	ViewNum    int64  `json:"view_num"`
	TagIds     string `json:"tag_ids"`
	CreateTime int64  `json:"create_time"`
}

type QuestionListResponse struct {
	Items  []*QuestionItem `json:"items"`
	Cursor int64           `json:"cursor"`
	IsEnd  bool            `json:"is_end"`
}

type QuestionDetailResponse struct {
	Question *QuestionItem `json:"question"`
}

type QuestionDeleteRequest struct {
	QuestionId int64 `json:"question_id"`
}

type AnswerListRequest struct {
	QuestionId int64 `form:"question_id"`
	Cursor     int64 `form:"cursor"`
	PageSize   int64 `form:"page_size"`
}

type AnswerItem struct {
	Id         int64  `json:"id"`
	QuestionId int64  `json:"question_id"`
	AuthorId   int64  `json:"author_id"`
	Content    string `json:"content"`
	IsAccepted bool   `json:"is_accepted"`
	LikeNum    int64  `json:"like_num"`
	ReplyNum   int64  `json:"reply_num"`
	CreateTime int64  `json:"create_time"`
}

type AnswerListResponse struct {
	Items  []*AnswerItem `json:"items"`
	Cursor int64         `json:"cursor"`
	IsEnd  bool          `json:"is_end"`
}

type AnswerDeleteRequest struct {
	AnswerId int64 `json:"answer_id"`
}

type SearchQuestionsRequest struct {
	Keyword  string `form:"keyword"`
	Cursor   int64  `form:"cursor"`
	PageSize int64  `form:"page_size"`
}

type SearchQuestionItem struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"author_id"`
	AnswerNum  int64  `json:"answer_num"`
	TagIds     string `json:"tag_ids"`
	CreateTime int64  `json:"create_time"`
}

type SearchQuestionsResponse struct {
	Items  []*SearchQuestionItem `json:"items"`
	Cursor int64                 `json:"cursor"`
	IsEnd  bool                  `json:"is_end"`
}
