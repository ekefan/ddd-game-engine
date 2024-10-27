package session

import (
	"errors"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ErrInvalidPlayer       = errors.New("a player must have an active socket connection")
	ErrInvalidMove         = errors.New("invalid move")
	ErrInvalidRoundOutcome = errors.New("invalid round outcome")
	ErrInvalidFlag         = errors.New("flag can either be 0 for Player1Flag and 1 for Player2Flag")
	ErrNoSession           = errors.New("no game session has been created")
	ErrPlayerMissing       = errors.New("a session can start only with two players")
)

const (
	InitRound      = 1
	MaxRound       = 3
	Player1Flag    = 0
	Player2Flag    = 1
	MoveBufferSize = 2
)

// Session represents a single instance of the game
//
// it would have two players playing in a match
//
// An instance is created with it's factory NewSession()
type Session struct {
	match       *domain.Match
	player1     *domain.Player
	player2     *domain.Player
	move        chan domain.PlayerMove
}

func NewSession(player1 *domain.Player) *Session {
	match := &domain.Match{
		ID:    uuid.New(),
		Round: InitRound,
	}
	return &Session{
		match:       match,
		player1:     player1,
		move: make(chan domain.PlayerMove, MoveBufferSize),
	}
}

// GetID retuns the current session id which is the current match id
func (s *Session) GetID() uuid.UUID {
	if s.match == nil {
		return uuid.Nil
	}
	return s.match.ID
}

// TODO: yet to test these
func (s *Session) SetPlayer2(player *domain.Player) error {
	if player.Connection == nil {
		return ErrInvalidPlayer
	}
	s.player2 = player
	return nil
}
func (s *Session) GetRound() int {
	return s.match.Round
}

func (s *Session) GetPlayer1() (player1 *domain.Player) {
	return s.player1
}

func (s *Session) getPlayerPoint() (player1Point, player2Point int) {
	return int(s.player1.Points), int(s.player2.Points)
}
func (s *Session) Ended() bool {
	round := s.GetRound()
	player1Point, player2Point := s.getPlayerPoint()
	if round == MaxRound {
		if player1Point != player2Point {
			return true
		}
	}
	if round > MaxRound {
		if player1Point != player2Point {
			return true
		}
	}
	return false
}

// TODO: refactor
func (s *Session) Write() {
	defer s.player1.Connection.Close()
	defer s.player2.Connection.Close()
	for {
		move1 := <-s.move
		move2 := <-s.move
		if move1.Conn == s.player1.Connection{
			s.player1.Move = move1.Move
		}
		if move2.Conn == s.player1.Connection {
			s.player1.Move = move1.Move
		}
		if move1.Conn == s.player2.Connection{
			s.player2.Move = move2.Move
		}
		if move2.Conn == s.player2.Connection {
			s.player2.Move = move2.Move
		}
		s.player1.Connection.WriteMessage(websocket.TextMessage, []byte("move received"))
		s.player2.Connection.WriteMessage(websocket.TextMessage, []byte("move received"))
	}
}


func (s *Session) Read() {
	defer s.player1.Connection.Close()
	defer s.player2.Connection.Close()
	type WrongMoveResp struct {
		msg string
	}
	for {
		_, msg, err := s.player1.Connection.ReadMessage()
		if err != nil {
			return
		}
		move1, err := parseMove(msg)
		if err != nil {
			resp := WrongMoveResp {
				msg: "wrong move, Move must be rock paper or scissor",
			}
			s.player1.Connection.ReadJSON(resp)
		}

		player1Move := domain.PlayerMove{
			Move: move1,
			Conn: s.player1.Connection,
		}
		_, msg, err = s.player2.Connection.ReadMessage()
		if err != nil {
			return
		}
		move2, err := parseMove(msg)
		if err != nil {
			resp := WrongMoveResp {
				msg: "wrong move, Move must be rock paper or scissor",
			}
			s.player2.Connection.ReadJSON(resp)
		}
		player2Move := domain.PlayerMove{
			Move: move2,
			Conn: s.player2.Connection,
		}

		s.move <- player1Move
		s.move <- player2Move
		
	}
}
