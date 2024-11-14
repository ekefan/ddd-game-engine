package service

import (
	"context"
	"errors"

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


// TODO: refactor code
func (gs *GameService) PlayGame(id uuid.UUID, player2 *domain.Player) error {

	// gets the session id for the game to be played
	sess, err := gs.GetSession(id)
	if err != nil {
		player2.Connection.WriteMessage(websocket.TextMessage, []byte("session invalid, player must create a session"))
		player2.Connection.Close()
		return err
	}

	// get player one for the game to be played
	player1 := sess.GetPlayer1()
	if player1 == nil {
		gs.sessRepo.DeleteSession(id)
		player2.Connection.WriteMessage(websocket.TextMessage, []byte("to play a game there must be two players"))
		player2.Connection.Close()
		return errors.New("to play a game there must be two players")
	}

	// set player two for the game to be played
	if err := sess.SetPlayer2(player2); err != nil {
		gs.sessRepo.DeleteSession(id)
		return err
	}

	// send game started message to both players for the game to be played
	player1.Connection.WriteMessage(websocket.TextMessage, []byte("game started"))
	player2.Connection.WriteMessage(websocket.TextMessage, []byte("game started"))

	// TODO: handle session contexts
	ctx, endSession := context.WithCancel(context.Background())

	// start game on a different go routine
	go func(c context.Context, s *session.Session) {
		<- c.Done()
		s.SendWinMessage()
		s.GetPlayer1().Connection.Close()
		s.GetPlayer2().Connection.Close()
		gs.sessRepo.DeleteSession(s.GetID())
	}(ctx, sess)
	go sess.WriteRoundOutcome(endSession)
	sess.ReadPlayerMoves(endSession)
	return nil
}
