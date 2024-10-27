package repository

import (
	"errors"

	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound       = errors.New("session was not found")
	ErrFailedToCreateSession = errors.New("failed to add session")
	ErrFailedToUpdateSession = errors.New("failed to update session")
	ErrFailedToDeleteSession = errors.New("failed to delete session")
)

type SessionRepository interface {
	GetSession(id uuid.UUID) (*session.Session, error)
	CreateSession(s *session.Session) error
	UpdateSession(s *session.Session) error
	DeleteSession(id uuid.UUID) error
}
