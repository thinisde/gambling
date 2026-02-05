package websocket

import (
	"fmt"
	"sync"

	ws "github.com/gorilla/websocket"
)

var (
	rooms   = make(map[int64]*WebSocketRoom)
	roomsMu sync.RWMutex
)

type memberConn struct {
	conn *ws.Conn
	mu   sync.Mutex // serialize writes for this connection
}

type WebSocketRoom struct {
	mu         sync.RWMutex
	members    map[int64]*memberConn
	id         int64
	maxMembers int
	gameType   string
	private    bool
	password   string
}

func NewWebSocketRoom(id int64, maxMembers int, gameType string, private bool, password string) *WebSocketRoom {
	return &WebSocketRoom{
		members:    make(map[int64]*memberConn),
		id:         id,
		maxMembers: maxMembers,
		gameType:   gameType,
		private:    private,
		password:   password,
	}
}

// Rooms management helpers (lock the global rooms map)
func AddRoom(room *WebSocketRoom) {
	roomsMu.Lock()
	defer roomsMu.Unlock()
	rooms[room.id] = room
}

func getRoom(roomID int64) (*WebSocketRoom, bool) {
	roomsMu.RLock()
	defer roomsMu.RUnlock()
	room, ok := rooms[roomID]
	return room, ok
}

func currentAvaiableRooms() []WebSocketRoom {
	roomsMu.RLock()
	defer roomsMu.RUnlock()

	res := make([]WebSocketRoom, 0, len(rooms))
	for _, room := range rooms {
		// Copy only “public” fields (don’t leak password)
		res = append(res, WebSocketRoom{
			id:         room.id,
			maxMembers: room.maxMembers,
			gameType:   room.gameType,
			private:    room.private,
			// members/password omitted intentionally
		})
	}
	return res
}

// Room member ops (lock members map)
func (room *WebSocketRoom) AddMember(userID int64, conn *ws.Conn) {
	room.mu.Lock()
	defer room.mu.Unlock()
	room.members[userID] = &memberConn{conn: conn}
}

func (room *WebSocketRoom) RemoveMember(userID int64) {
	room.mu.Lock()
	defer room.mu.Unlock()
	delete(room.members, userID)
}

// Broadcast json message to all members
func (room *WebSocketRoom) Broadcast(event SystemEvent, message any) {
	// snapshot connections under read lock, then write without holding the room lock
	room.mu.RLock()
	conns := make([]*memberConn, 0, len(room.members))
	for _, mc := range room.members {
		conns = append(conns, mc)
	}
	room.mu.RUnlock()

	for _, mc := range conns {
		mc.mu.Lock()
		err := mc.conn.WriteJSON(map[string]any{
			"event_type": event,
			"data":       message,
		})
		if err != nil {
			// handle error (e.g., log it, remove member, etc.)
			fmt.Println("Error broadcasting to member:", err)
		}
		mc.mu.Unlock()
	}
}

func (room *WebSocketRoom) ID() int64 { return room.id }
