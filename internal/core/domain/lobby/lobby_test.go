package lobby_test

import (
	"testing"

	"github.com/ekefan/ddd-game-engine/internal/core/domain/lobby"
	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createNewLobby(t *testing.T) *lobby.Lobby {
	newLobby, err := lobby.NewLobby()
	require.NoError(t, err)
	require.NotEmpty(t, newLobby)

	return newLobby
}

func TestNewLobby(t *testing.T) {
	createNewLobby(t)
}

func randomPlayer() domain.Player {
	return domain.Player{
		ID:         uuid.New(),
		Connection: &websocket.Conn{},
	}
}
func TestAddPlayer(t *testing.T) {
	type testCase struct {
		test        string
		player      domain.Player
		expectedErr error
	}

	pl := randomPlayer()
	lb := createNewLobby(t)
	err := lb.AddPlayer(pl)
	assert.NoError(t, err)

	testCases := []testCase{
		{
			test:        "valid case, player doesn't already exist",
			player:      randomPlayer(),
			expectedErr: nil,
		}, {
			test:        "player already exist",
			player:      pl,
			expectedErr: lobby.ErrLobbyPlayerExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T){
			err := lb.AddPlayer(tc.player)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetPlayer(t *testing.T) {
	lb := createNewLobby(t)
	pl := randomPlayer()

	err := lb.AddPlayer(pl)
	assert.NoError(t, err)

	type testCase struct {
		test string
		id uuid.UUID
		expectedErr error
	}
	testCases := []testCase {
		{
			test: "player exist in lobby",
			id: pl.ID,
			expectedErr: nil,
		}, {
			test: "player doesn't exist in lobby",
			id: uuid.New(),
			expectedErr: lobby.ErrLobbyPlayerNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T){
			_, err := lb.GetPlayer(tc.id)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRemovePlayer(t *testing.T) {
	lb := createNewLobby(t)
	pl := randomPlayer()

	err := lb.AddPlayer(pl)
	assert.NoError(t, err)

	type testCase struct {
		test string
		id uuid.UUID
		expectedErr error
	}
	testCases := []testCase {
		{
			test: "player exist in lobby",
			id: pl.ID,
			expectedErr: nil,
		}, {
			test: "player doesn't exist in lobby",
			id: uuid.New(),
			expectedErr: lobby.ErrLobbyPlayerNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T){
			err := lb.RemovePlayer(tc.id)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}