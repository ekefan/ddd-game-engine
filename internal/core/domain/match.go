package domain

import (
	"github.com/google/uuid"
)

// Match is an entity that represents a rock-paper-scissors match in all sub-domains
type Match struct {
	ID uuid.UUID
	Round int
	RoundOutcome RoundOutcome
}