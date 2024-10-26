package memory

import (
	"testing"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	session_repo "github.com/ekefan/ddd-game-engine/internal/ports/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSession(t *testing.T) session.Session {
	ns := session.NewSession(&domain.Player{}, &domain.Player{})

	require.NotEmpty(t, ns)
	return *ns
}

func TestGet(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		expectedErr error
	}

	ns := newSession(t)
	id := ns.GetID()
	repo := NewSessionRepository()
	repo.sessions[id] = ns
	testCases := []testCase{
		{
			test:        "session exists",
			id:          id,
			expectedErr: nil,
		}, {
			test:        "session doesn't exist",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: session_repo.ErrSessionNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := repo.Get(tc.id)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestAdd(t *testing.T) {
	type testCase struct {
		test        string
		sess        session.Session
		expectedErr error
	}
	ns := newSession(t)
	sessions := NewSessionRepository()
	sessions.sessions[ns.GetID()] = ns
	testCases := []testCase{
		{
			test:        "session already exists",
			sess:        ns,
			expectedErr: session_repo.ErrFailedToAddSession,
		}, {
			test:        "session doesn't exist",
			sess:        newSession(t),
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := sessions.Add(tc.sess)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
func TestUpdate(t *testing.T) {
	type testCase struct {
		test        string
		sess        session.Session
		expectedErr error
	}
	ns := newSession(t)
	sessions := NewSessionRepository()
	sessions.sessions[ns.GetID()] = ns
	testCases := []testCase{
		{
			test:        "session doesn't exist",
			sess:        newSession(t),
			expectedErr: session_repo.ErrFailedToUpdateSession,
		}, {
			test:        "session exist",
			sess:        ns,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := sessions.Update(tc.sess)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
func TestDelete(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		expectedErr error
	}

	ns := newSession(t)
	id := ns.GetID()
	repo := NewSessionRepository()
	repo.sessions[id] = ns
	testCases := []testCase{
		{
			test:        "session exists",
			id:          id,
			expectedErr: nil,
		}, {
			test:        "session doesn't exist",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: session_repo.ErrFailedToDeleteSession,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := repo.Delete(tc.id)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
