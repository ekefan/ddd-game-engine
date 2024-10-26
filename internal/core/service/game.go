package service

import (
	"errors"
	"sync"
	repo "github.com/ekefan/ddd-game-engine/internal/ports/repository"
)

var (
	ErrNoGameServer = errors.New("can not create game server")
)

type GameService struct {
	sessionRepository repo.SessionRepository
	lobbyRespository  repo.LobbyRepository
	sync.Mutex
}



func NewGameService(sessionRepo repo.SessionRepository, lobbyRepo repo.LobbyRepository) *GameService {
	return &GameService{
		sessionRepository: sessionRepo,
		lobbyRespository: lobbyRepo,
	}
}

// TODO: write tests for this service
func (gs *GameService) WaitInLobby() error {
	return nil
}


func (gs *GameService) PlayGame() error {
	return nil
}