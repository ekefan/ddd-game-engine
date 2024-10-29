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
	ErrInvalidPlayer       = errors.New("a player must have an active socket connection")
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
	player1move  chan domain.PlayerMove
	player2move  chan domain.PlayerMove
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
		player1move: make(chan domain.PlayerMove),
		player2move: make(chan domain.PlayerMove),
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

func (s *Session) GetResponse() *domain.Response {
	s.generateResponse()
	return s.response
}

// TODO: refactor
func (s *Session) Write(endSession context.CancelFunc) {
	for !s.isSessionEnded() {
		// receive moves from player connection
		move1 := <-s.player1move
		fmt.Printf(" player 1 move received from session:%v\n", s.GetID())
		move2 := <-s.player2move
		fmt.Printf(" player 2 move received from session:%v\n", s.GetID())
		s.player1.Move = move1.Move
		s.player2.Move = move2.Move

		// determine RoundOutcome
		s.DetermineRoundOutcome()
		if s.sessionEnded {
			fmt.Println("session ended first check", s.sessionEnded)
			s.player1.Connection.WriteJSON(s.response)
			s.player2.Connection.WriteJSON(s.response)

			fmt.Println("game ended")
			s.player1.Connection.WriteJSON(s.GetWinResponse())
			s.player2.Connection.WriteJSON(s.GetWinResponse())
			fmt.Println("last connection sent")

			endSession()
			return
		}
		s.player1.Connection.WriteJSON(s.response)
		s.player2.Connection.WriteJSON(s.response)
		fmt.Println("session ended second check", s.sessionEnded)
	}
	s.player1.Connection.Close()
	s.player2.Connection.Close()
}

func (s *Session) Read(endSession context.CancelFunc) {
	type WrongMoveResp struct {
		Msg string `json:"msg"`
	}

	for !s.isSessionEnded() {
		player1Move := &domain.PlayerMove{}
		player2Move := &domain.PlayerMove{}
		validPlayer1Move := false
		validPlayer2Move := false

		// Attempt to read moves from both players
		for !validPlayer1Move || !validPlayer2Move {
			// Read Player 1's move
			if !validPlayer1Move {
				_, msg1, err := s.player1.Connection.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err) {
						fmt.Printf("session with id: %v has a disconnected player", s.GetID())
					}
					endSession()
					return
				}

				move1, err := parseMove(msg1)
				if err != nil {
					resp := WrongMoveResp{Msg: "wrong move, Move must be rock, paper, or scissors"}
					s.player1.Connection.WriteJSON(resp)
					continue // Wait for a valid move
				}

				player1Move = &domain.PlayerMove{
					Move: move1,
					Conn: s.player1.Connection,
				}
				validPlayer1Move = true // Mark as valid once read successfully
			}

			// Read Player 2's move
			if !validPlayer2Move {
				_, msg2, err := s.player2.Connection.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err) {
						fmt.Printf("session with id: %v should be closed", s.GetID())
						return
					}
					endSession()
					return
				}

				move2, err := parseMove(msg2)
				if err != nil {
					resp := WrongMoveResp{Msg: "wrong move, Move must be rock, paper, or scissors"}
					s.player2.Connection.WriteJSON(resp)
					continue // Wait for a valid move
				}


				player2Move = &domain.PlayerMove{
					Move: move2,
					Conn: s.player2.Connection,
				}
				validPlayer2Move = true // Mark as valid once read successfully
			}
		}

		// Once both moves are valid, send them to the respective channels
		s.player1move <- *player1Move
		s.player2move <- *player2Move
	}
}

func (s *Session) DetermineRoundOutcome() {
	s.setPlayerPoints()
	if s.gameEnded() {
		s.generateResponse()
		return
	}
	s.generateResponse()
}

type WinResp struct {
	WinMsg string `json:"message"`
	Winner string `json:"winner"`
}

func (s *Session) GetWinResponse() *WinResp {
	return &WinResp{
		WinMsg: "session ended",
		Winner: s.getSessionWinner(),
	}
}


// use session with context wait for the cancel signal from the context then delete the session