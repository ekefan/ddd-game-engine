package domain

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)
const (
	DefaultPlayer1Name = "player1"
	DefaultPlayer2Name = "player2"
)

// Player is an entity that represents a player in all sub-domains
type Player struct {
	ID         uuid.UUID
	Name       string
	Move       Move
	Connection *websocket.Conn
	Points     int16
}
