package repository

import (
	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/google/uuid"
)

type LobbyRepository interface {
	AddPlayer(p domain.Player) error
	GetPlayer(id uuid.UUID) (domain.Player, error)
	RemovePlayer(id uuid.UUID) error
}
