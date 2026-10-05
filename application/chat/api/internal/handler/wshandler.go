package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"ThinkTalk/application/chat/api/internal/svc"
	"ThinkTalk/application/chat/api/internal/types"
	"ThinkTalk/application/chat/rpc/pb"

	"github.com/golang-jwt/jwt/v4"
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
		tokenStr := r.URL.Query().Get("token")
		if tokenStr == "" {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenStr = authHeader[7:]
			}
		}

		if tokenStr == "" {
			http.Error(w, "unauthorized: token missing", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(svcCtx.Config.Auth.AccessSecret), nil
		})
		if err != nil || !token.Valid {
			logx.Errorf("[WS] token parse err: %v", err)
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "unauthorized: invalid claims", http.StatusUnauthorized)
			return
		}

		var uid int64
		if userIdVal, exists := claims["userId"]; exists {
			switch val := userIdVal.(type) {
			case float64:
				uid = int64(val)
			case json.Number:
				uid, _ = val.Int64()
			case string:
				// just in case
			}
		}

		if uid == 0 {
			http.Error(w, "unauthorized: invalid user id", http.StatusUnauthorized)
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
