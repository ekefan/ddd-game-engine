package domain

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Player is an entity that represents a player in all sub-domains
type Player struct {
	ID         uuid.UUID
	Name       string
	Connection *websocket.Conn
	Move       Move
	Points     int16
}
