package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
)

var (
	ErrSessionNotFound       = errors.New("session was not found")
	ErrFailedToAddSession    = errors.New("failed to add session")
	ErrFailedToUpdateSession = errors.New("failed to update session")
	ErrFailedToDeleteSession = errors.New("failed to delete session")
)

type SessionRepository interface {
	Get(id uuid.UUID) (session.Session, error)
	Add(s session.Session) error
	Update(s session.Session) error
	Delete(id uuid.UUID) error
}
