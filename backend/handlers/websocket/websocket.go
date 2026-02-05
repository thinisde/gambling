package websocket

import (
	"fmt"
	"net/http"
	"time"

	"backend/cache"
	"backend/handlers/auth"

	"github.com/gorilla/mux"
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
	RoomID     int64  `json:"room_id,omitempty"`
	Password   string `json:"password,omitempty"`
	CreateRoom bool   `json:"create_room,omitempty"`
	RoomName   string `json:"room_name,omitempty"`
	Private    bool   `json:"private,omitempty"`
	MaxMembers int    `json:"max_members,omitempty"`
	GameType   string `json:"game_type,omitempty"`
}

func CreateWebsocketHandler(r *mux.Router) {
	r.HandleFunc("/ws", WebSocketHandler).Methods("GET")
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

	// get user roomID from redis
	roomIDStr, err := cache.CacheClient.HGet(fmt.Sprintf("user:%d:ws", userID), "room_id")
	if err != nil {
		http.Error(w, "Failed to get user room", http.StatusInternalServerError)
		return
	}

	var roomID int64
	_, err = fmt.Sscanf(roomIDStr, "%d", &roomID)

	for {
		if _, p, err := conn.ReadMessage(); err != nil {
			HandleSystemEvents(p, userID, roomID)
		}
	}
}
