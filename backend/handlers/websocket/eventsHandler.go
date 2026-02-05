package websocket

// Create an enum of events
type SystemEvent int

const (
	EventTypeMessage SystemEvent = iota
	EventTypeJoin
	EventTypeLeave
	EventTypeGameStart
	EventTypeGameEnd
)

func CreateEventMessage(eventType SystemEvent, data any) map[string]any {
	return map[string]any{
		"event_type": eventType,
		"data":       data,
	}
}

// events handler
func HandleSystemEvents(message []byte, userID int64, roomID int64) {
	roomsMu.RLock()
	defer roomsMu.RUnlock()

	room, ok := rooms[roomID]
	if !ok {
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()
	// Get message event type and data
	eventType := SystemEvent(message[0])
	data := message[1:]

	switch eventType {
	case EventTypeMessage:
		room.Broadcast(EventTypeMessage, data)
	}
}
