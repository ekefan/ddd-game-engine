package gameserver

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ekefan/ddd-game-engine/internal/adapters/memory"
	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/ekefan/ddd-game-engine/internal/core/domain/lobby"
	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	repo "github.com/ekefan/ddd-game-engine/internal/ports/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ErrNoGameServer = errors.New("can not create game server")
)

type GameServer struct {
	sessionRepository repo.SessionRepository
	lobbyRespository  repo.LobbyRepository
	player1Move       chan domain.Move
	player2Move       chan domain.Move
	sync.Mutex
}

type GameServerConfiguration func(gs *GameServer) error

func NewGameServer(cfgs ...GameServerConfiguration) error {
	gs := &GameServer{
		player1Move: make(chan domain.Move),
		player2Move: make(chan domain.Move),
	}
	for _, cfg := range cfgs {
		cfg(gs)
	}
	if gs.sessionRepository == nil || gs.lobbyRespository == nil {
		return ErrNoGameServer
	}
	return nil
}

func WithSessionRepository(s session.Session) GameServerConfiguration {
	sr := memory.NewSessionRepository()
	err := sr.Add(s)
	return func(gs *GameServer) error {
		gs.sessionRepository = sr
		if err != nil {
			return err
		}
		return nil
	}
}
func WithLobbyRepository() GameServerConfiguration {
	lr, err := lobby.NewLobby()
	return func(gs *GameServer) error {
		gs.lobbyRespository = lr
		if err != nil {
			return err
		}
		return nil
	}
}

// TODO: write tests for this service
func (gs *GameServer) EnterLobby() error {
	return nil
}

type PlayGameReq struct {
	SessionID uuid.UUID
}
func (gs *GameServer) PlayGame(req PlayGameReq) error {
	return nil
}

type MakeMoveReq struct {
	PlayerConn *websocket.Conn
	PlayerMove domain.Move
}

func (gs *GameServer) MakeMove(req MakeMoveReq) error {
	gs.player1Move <- req.PlayerMove
	return nil
}

func (gs *GameServer) EndGameSession() error {
	return nil
}

func (gs *GameServer) bothPlayersAreConnected(id uuid.UUID) (bool, error) {
	sess, err := gs.sessionRepository.Get(id)
	if err != nil {
		return false, err
	}
	player1Conn, player2Conn := sess.GetPlayersConn()
	if isWebSocketOpen(player1Conn) || isWebSocketOpen(player2Conn) {
		return true, nil
	}
	return false, nil
}

func isWebSocketOpen(conn *websocket.Conn) bool {
	// Set a short read deadline to test connection
	if err := conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return false
	}

	// Try to read the next control message (like Ping/Pong)
	_, _, err := conn.ReadMessage()

	if err != nil {
		
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Println("Connection closed:", err)
			return false
		}
	}

	// Reset the read deadline after the check
	conn.SetReadDeadline(time.Time{}) // Disable the deadline
	return true
}
