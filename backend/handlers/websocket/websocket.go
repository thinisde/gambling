package websocket

import (
	"fmt"
	"net/http"
	"time"

	"backend/handlers/auth"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		fmt.Println("Origin:", r.Header.Get("Origin"))
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000"
	},
}

type joinRoomRequest struct {
	RoomID   int64  `json:"room_id"`
	Password string `json:"password,omitempty"`
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value(auth.UserIDKey)
	userID, ok := v.(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return // can't http.Error reliably here either; upgrade failed
	}
	defer conn.Close()

	// deadlines + pong handler
	conn.SetReadLimit(1 << 20) // 1MB
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	// Send rooms list
	availableRooms := currentAvaiableRooms() // make sure rooms access is thread-safe
	if err := conn.WriteJSON(availableRooms); err != nil {
		return
	}

	// Join room
	var joinReq joinRoomRequest
	if err := conn.ReadJSON(&joinReq); err != nil {
		_ = conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseUnsupportedData, "bad join request"))
		return
	}

	room, ok := getRoom(joinReq.RoomID) // recommended: wrap map access with mutex
	if !ok {
		_ = conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, "room not found"))
		return
	}

	if room.private && room.password != joinReq.Password {
		_ = conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, "invalid room password"))
		return
	}

	room.AddMember(userID, conn)
	defer room.RemoveMember(userID)

	// Read loop
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
