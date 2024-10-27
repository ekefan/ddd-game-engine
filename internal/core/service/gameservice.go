package service

import (
	"errors"
	"sync"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	repo "github.com/ekefan/ddd-game-engine/internal/ports/repository"
	"github.com/google/uuid"
)

var (
	ErrNoGameServer = errors.New("can not create game server")
)

type GameService struct {
	sessRepo repo.SessionRepository
	sync.Mutex
}

func NewGameService(sessionRepo repo.SessionRepository) *GameService {
	return &GameService{
		sessRepo: sessionRepo,
	}
}

// TODO: write tests for this service
func (gs *GameService) CreateSession(player *domain.Player) (*session.Session, error) {
	sess := session.NewSession(player)
	err := gs.sessRepo.CreateSession(sess)
	if err != nil {
		return &session.Session{}, err
	}
	return sess, nil
}

func (gs *GameService) GetSession(id uuid.UUID) (*session.Session, error) {
	return gs.sessRepo.GetSession(id)
}


func (gs *GameService) PlayGame(id uuid.UUID, player2 *domain.Player) error {
	sess, err := gs.GetSession(id)
	if err != nil {
		return err
	}
	player1 := sess.GetPlayer1()
	if player1 == nil {
		gs.sessRepo.DeleteSession(id)
		return errors.New("to play a game there must be two players")
	}

	if err := sess.SetPlayer2(player2); err != nil {
		gs.sessRepo.DeleteSession(id)
		return err
	}
	
	go sess.Write()
	sess.Read()
	return nil
}
