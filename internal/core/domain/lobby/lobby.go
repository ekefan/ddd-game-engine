package lobby

import (
	"errors"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrLobbyInvalidPlayer = errors.New("player must have both an id and be connected to the server")
	ErrLobbyPlayerNotFound = errors.New("player not found in lobby")
	ErrLobbyPlayerExists = errors.New("player already exists in lobby")
)

// Lobby represents the waiting room for players before they can
// begin a game session.
//
// Players are added to Lobby when they want to play in a game session
// an array is used to hold players for extensibility (if the game supports random game sessions)
type Lobby struct {
	players map[uuid.UUID]domain.Player
}

func NewLobby() (*Lobby, error) {
	return &Lobby{
		make(map[uuid.UUID]domain.Player),
	}, nil
}

func (lobby *Lobby) AddPlayer(p domain.Player) error {
	if p.Connection == nil || p.ID == uuid.Nil {
		return ErrLobbyInvalidPlayer
	}
	if _, exist := lobby.players[p.ID]; exist {
		return ErrLobbyPlayerExists
	}
	lobby.players[p.ID] = p
	return nil
}
func (lobby *Lobby) GetPlayer(id uuid.UUID) (domain.Player, error) {
	if _, exist := lobby.players[id]; !exist {
		return domain.Player{}, ErrLobbyPlayerNotFound
	}
	return lobby.players[id], nil
}

func (lobby *Lobby) RemovePlayer(id uuid.UUID) error {
	if _, exist := lobby.players[id]; !exist {
		return ErrLobbyPlayerNotFound
	}

	delete(lobby.players, id)
	return nil
}