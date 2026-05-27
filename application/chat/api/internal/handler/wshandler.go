package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/application/chat/api/internal/types"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	hub := svcCtx.Hub

	hub.OnMessage = func(userId int64, msg *types.WsInMessage) {
		switch msg.Type {
		case "message":
			_, err := svcCtx.Chat.SendMessage(context.Background(), &pb.SendMessageRequest{
				SenderId:   userId,
				ReceiverId: msg.ReceiverId,
				Content:    msg.Content,
				MsgType:    msg.MsgType,
			})
			if err != nil {
				logx.Errorf("[WS] send message err: %v userId: %d", err, userId)
				hub.SendToUser(userId, &types.WsOutMessage{
					Type:    "error",
					Message: err.Error(),
				})
				return
			}
			hub.SendToUser(userId, &types.WsOutMessage{
				Type:    "ack",
				Message: "sent",
			})
		default:
			hub.SendToUser(userId, &types.WsOutMessage{
				Type:    "error",
				Message: "unknown message type",
			})
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(json.Number)
		uid, err := userId.Int64()
		if err != nil || uid == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Errorf("[WS] upgrade err: %v", err)
			return
		}

		conn := hub.Register(uid, ws)
		logx.Infof("[WS] user connected userId: %d", uid)

		go hub.WritePump(conn)
		go hub.ReadPump(conn)
	}
}
