package memory

import (
	"sync"

	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	session_repo "github.com/ekefan/ddd-game-engine/internal/ports/repository"
	"github.com/google/uuid"
)

type MemoryRepository struct {
	sessions map[uuid.UUID]session.Session
	sync.Mutex
}

func NewSessionRepository() *MemoryRepository {
	return &MemoryRepository{
		sessions: make(map[uuid.UUID]session.Session),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (session.Session, error) {
	if session, ok := mr.sessions[id]; ok {
		return session, nil
	}
	return session.Session{}, session_repo.ErrSessionNotFound
}

func (mr *MemoryRepository) Add(s session.Session) error {
	if mr.sessions == nil {
		mr.Lock()
		mr.sessions = make(map[uuid.UUID]session.Session)
		mr.Unlock()
	}
	if _, ok := mr.sessions[s.GetID()]; ok {
		return session_repo.ErrFailedToAddSession
	}
	mr.Lock()
	defer mr.Unlock()
	mr.sessions[s.GetID()] = s
	return nil
}

func (mr *MemoryRepository) Update(sess session.Session) error {
	id := sess.GetID()
	_, ok := mr.sessions[id]
	if !ok {
		return session_repo.ErrFailedToUpdateSession
	}
	mr.Lock()
	defer mr.Unlock()
	mr.sessions[id] = sess
	return nil
}

func (mr *MemoryRepository) Delete(id uuid.UUID) error {
	if _, ok := mr.sessions[id]; !ok {
		return session_repo.ErrFailedToDeleteSession
	}
	mr.Lock()
	defer mr.Unlock()
	delete(mr.sessions, id)
	return nil
}