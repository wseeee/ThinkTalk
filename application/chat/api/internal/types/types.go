package types

const DefaultPageSize = 20

type SendMessageRequest struct {
	ReceiverId int64  `json:"receiver_id"`
	Content    string `json:"content"`
	MsgType    int32  `json:"msg_type"`
}

type SendMessageResponse struct{}

type ConversationsRequest struct {
	Cursor   int64 `form:"cursor"`
	PageSize int64 `form:"page_size"`
}

type ConversationItem struct {
	Id              int64  `json:"id"`
	TargetUserId    int64  `json:"target_user_id"`
	LastMessage     string `json:"last_message"`
	LastMessageTime int64  `json:"last_message_time"`
	UnreadCount     int64  `json:"unread_count"`
}

type ConversationsResponse struct {
	Items  []*ConversationItem `json:"items"`
	Cursor int64               `json:"cursor"`
	IsEnd  bool                `json:"is_end"`
}

type MessagesRequest struct {
	ConversationId int64 `form:"conversation_id"`
	Cursor         int64 `form:"cursor"`
	PageSize       int64 `form:"page_size"`
}

type MessageItem struct {
	Id             int64  `json:"id"`
	ConversationId int64  `json:"conversation_id"`
	SenderId       int64  `json:"sender_id"`
	ReceiverId     int64  `json:"receiver_id"`
	Content        string `json:"content"`
	MsgType        int32  `json:"msg_type"`
	IsRead         bool   `json:"is_read"`
	CreateTime     int64  `json:"create_time"`
}

type MessagesResponse struct {
	Items  []*MessageItem `json:"items"`
	Cursor int64          `json:"cursor"`
	IsEnd  bool           `json:"is_end"`
}

type MarkReadRequest struct {
	ConversationId int64 `json:"conversation_id"`
}

type MarkReadResponse struct{}

type UnreadCountResponse struct {
	Total int64 `json:"total"`
}

type MessageResponse struct {
	Id             int64  `json:"id"`
	ConversationId int64  `json:"conversation_id"`
	SenderId       int64  `json:"sender_id"`
	ReceiverId     int64  `json:"receiver_id"`
	Content        string `json:"content"`
	MsgType        int32  `json:"msg_type"`
	CreateTime     int64  `json:"create_time"`
}

type WsInMessage struct {
	Type       string `json:"type"`
	ReceiverId int64  `json:"receiver_id,omitempty"`
	Content    string `json:"content,omitempty"`
	MsgType    int32  `json:"msg_type,omitempty"`
}

type WsOutMessage struct {
	Type       string `json:"type"`
	SenderId   int64  `json:"sender_id,omitempty"`
	Content    string `json:"content,omitempty"`
	MsgType    int32  `json:"msg_type,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
	Message    string `json:"message,omitempty"`
}
