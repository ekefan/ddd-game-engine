package session

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ErrInvalidPlayer       = errors.New("a player must have an active socket connection and a valid ID")
	ErrInvalidMove         = errors.New("invalid move")
	ErrInvalidRoundOutcome = errors.New("invalid round outcome")
	ErrInvalidFlag         = errors.New("flag can either be 0 for Player1Flag and 1 for Player2Flag")
	ErrNoSession           = errors.New("no game session has been created")
	ErrPlayerMissing       = errors.New("a session can start only with two players")
)

const (
	InitRound   = 0
	MaxRound    = 3
	Player1Flag = 0
	Player2Flag = 1
)

// Session represents a single instance of the game
//
// it would have two players playing in a match
//
// An instance is created with it's factory NewSession()
type Session struct {
	match        *domain.Match
	player1      *domain.Player
	player2      *domain.Player
	player1move  chan domain.Move
	player2move  chan domain.Move
	response     *domain.Response
	sessionEnded bool
	mu           sync.Mutex
}

func NewSession(player1 *domain.Player) *Session {
	match := &domain.Match{
		ID:    uuid.New(),
		Round: InitRound,
	}
	return &Session{
		match:       match,
		player1:     player1,
		player1move: make(chan domain.Move),
		player2move: make(chan domain.Move),
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
	if player.Connection == nil || player.ID == uuid.Nil {
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

func (s *Session) GetPlayer2() (player2 *domain.Player) {
	return s.player2
}

func (s *Session) GetResponse() *domain.Response {
	s.generateRoundResponse()
	return s.response
}

func (s *Session) SendWinMessage() {
	fmt.Println("game ended")
	err := s.player1.Connection.WriteJSON(s.getWinResponse())
	if err != nil {
		fmt.Println(err, "from sending player1 win message")
	}
	err = s.player2.Connection.WriteJSON(s.getWinResponse())
	if err != nil {
		fmt.Println(err, "from sending player2 win message")
	}
}

// TODO: refactor
func (s *Session) WriteRoundOutcome(endSession context.CancelFunc) {
	defer close(s.player1move)
	defer close(s.player2move)
	defer endSession()
	for !s.isSessionEnded() {
		// receive moves from player connection
		move1 := <-s.player1move
		move2 := <-s.player2move

		s.player1.Move = move1
		s.player2.Move = move2

		// determine RoundOutcome
		s.determineRoundOutcome()

		// send round outcome to players
		err := s.player1.Connection.WriteJSON(s.response)
		if err != nil {
			return
		}
		err = s.player2.Connection.WriteJSON(s.response)
		if err != nil {
			return
		}

		if s.sessionEnded {
			return
		}
	}
}

// ReadPlayerMoves reads from player connections to get player moves
//
// for tests
// receive domain.PlayerMove through session channels
// receive message from player channel one when two is disconnected and vice versa
// recevive message from player one when wrong move is sent on the channel perform for channel two too
func (s *Session) ReadPlayerMoves(endSession context.CancelFunc) {
	type WrongMoveResp struct {
		Msg string `json:"msg"`
	}
	for {
		var player1Move domain.Move
		var player2Move domain.Move
		validPlayer1Move := false
		validPlayer2Move := false

		// Attempt to read moves from both players
		for (!validPlayer1Move || !validPlayer2Move) && !s.sessionEnded {
			// Read Player 1's move
			if !validPlayer1Move {
				_, msg1, err := s.player1.Connection.ReadMessage()
				if err != nil {
					// TODO: there should be a way to send the reason for ending the session
					s.player2.Connection.WriteMessage(websocket.TextMessage, []byte("player1 disconnected"))
					endSession()
					return
				}

				move1, err := parseMove(msg1)
				if err != nil {
					resp := WrongMoveResp{Msg: "wrong move, Move must be rock, paper, or scissors"}
					s.player1.Connection.WriteJSON(resp)
					continue // Wait for a valid move
				}

				player1Move = move1
				validPlayer1Move = true // Mark as valid once read successfully
			}

			// Read Player 2's move
			if !validPlayer2Move {
				_, msg2, err := s.player2.Connection.ReadMessage()
				if err != nil {
					s.player1.Connection.WriteMessage(websocket.TextMessage, []byte("player1 disconnected"))
					endSession()
					return
				}

				move2, err := parseMove(msg2)
				if err != nil {
					resp := WrongMoveResp{Msg: "wrong move, Move must be rock, paper, or scissors"}
					s.player2.Connection.WriteJSON(resp)
					continue // Wait for a valid move
				}

				player2Move = move2
				validPlayer2Move = true // Mark as valid once read successfully
			}
		}
		// Once both moves are valid, send them to the respective channels
		if !s.sessionEnded {
			s.player1move <- player1Move
			s.player2move <- player2Move
		}

	}
}

func (s *Session) determineRoundOutcome() {
	s.setRoundOutcome()
	s.setPlayerPoints()
	s.updateRound()
	s.generateRoundResponse()
	s.checkGameEnded()
}

type WinResp struct {
	WinMsg string `json:"message"`
	Winner string `json:"winner"`
}

func (s *Session) getWinResponse() *WinResp {
	return &WinResp{
		WinMsg: "session ended",
		Winner: s.getSessionWinner(),
	}
}

// use session with context wait for the cancel signal from the context then delete the session
