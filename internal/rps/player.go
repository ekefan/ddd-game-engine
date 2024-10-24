package rps

import (
	"github.com/google/uuid"
)

// Player is an entity that represents a player in all sub-domains
type Player struct {
	ID uuid.UUID
	Name string
	Move Move
	Points int16
}