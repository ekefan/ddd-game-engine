package ports

import (
	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	"github.com/google/uuid"
)

type GameService interface {
	CreateSession(player *domain.Player) (*session.Session, error)
	GetSession(id uuid.UUID) (*session.Session, error)
	PlayGame(id uuid.UUID, player2 *domain.Player) error //play route web socket
}
