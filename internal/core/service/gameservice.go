package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	repo "github.com/ekefan/ddd-game-engine/internal/ports/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ErrNoGameServer = errors.New("can not create game server")
)

type GameService struct {
	sessRepo   repo.SessionRepository
	endSession chan bool
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
		player2.Connection.WriteMessage(websocket.TextMessage, []byte("session invalid, player must create a session"))
		player2.Connection.Close()
		return err
	}
	player1 := sess.GetPlayer1()
	if player1 == nil {
		gs.sessRepo.DeleteSession(id)
		player2.Connection.WriteMessage(websocket.TextMessage, []byte("to play a game there must be two players"))
		player2.Connection.Close()
		return errors.New("to play a game there must be two players")
	}

	if err := sess.SetPlayer2(player2); err != nil {
		gs.sessRepo.DeleteSession(id)
		return err
	}
	player1.Connection.WriteMessage(websocket.TextMessage, []byte("game started"))
	player2.Connection.WriteMessage(websocket.TextMessage, []byte("game started"))

	// TODO: handle session contexts
	fmt.Println("game started")
	endSession := make(chan bool)
	go sess.Write(endSession)
	sess.Read(endSession)
	<-endSession

	return nil
}
