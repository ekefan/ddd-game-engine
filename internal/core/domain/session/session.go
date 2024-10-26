package session

import (
	"errors"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ErrInvalidPlayer       = errors.New("a player name must not be an empty string")
	ErrInvalidMove         = errors.New("invalid move")
	ErrInvalidRoundOutcome = errors.New("invalid round outcome")
	ErrInvalidFlag         = errors.New("flag can either be 0 for Player1Flag and 1 for Player2Flag")
	ErrNoSession           = errors.New("no game session has been created")
	ErrPlayerMissing = errors.New("a session can start only with two players")
)

const (
	InitRound          = 1
	MaxRound           = 3
	DefaultPlayer1Name = "player1"
	DefaultPlayer2Name = "player2"
	Player1Flag        = 0
	Player2Flag        = 1
)

// Session represents a single instance of the game
//
// it would have two players playing in a match
//
// An instance is created with it's factory NewSession()
type Session struct {
	match   *domain.Match
	player1 *domain.Player
	player2 *domain.Player
}

type sessionConfiguration func (s *Session) error
func NewSession(cfgs ...sessionConfiguration) (Session, error) {
	match := &domain.Match{
		ID:    uuid.New(),
		Round: InitRound,
	}
	newSession := Session{
		match:   match,
	}
	for _, cfg := range cfgs {
		cfg(&newSession)
	}
	if newSession.player1 == nil || newSession.player2 == nil{
		return Session{}, ErrPlayerMissing
	}
	return newSession, nil
}


func WithPlayer1(playerConn *websocket.Conn, name string) sessionConfiguration {
	var playerName string
	if name == "" {
		playerName = DefaultPlayer1Name
	}
	player := domain.Player{
		ID: uuid.New(),
		Name: playerName,
		Connection: playerConn,
	}
	return func(s *Session) error {
		s.player1 = &player
		return nil
	}
}

func WithPlayer2(playerConn *websocket.Conn, name string) sessionConfiguration {
	var playerName string
	if name == "" {
		playerName = DefaultPlayer2Name
	}
	player := domain.Player{
		ID: uuid.New(),
		Name: playerName,
		Connection: playerConn,
	}
	return func(s *Session) error {
		s.player2 = &player
		return nil
	}
}

// GetID retuns the current session id which is the current match id
func (s *Session) GetID() uuid.UUID {
	if s.match == nil {
		return uuid.Nil
	}
	return s.match.ID
}

func (s *Session) GetRound() int {
	return s.match.Round
}

// UpdateRound update the current session's round
func (s *Session) UpdateRound() {
	s.match.Round++
}

// DetermineRoundOutcome checks the players move returns RoundOutcome
// on error RoundOutcome as -1, and associated error is return
func (s *Session) DetermineRoundOutCome() (domain.RoundOutcome, error) {
	if s.match == nil {
		return -1, ErrNoSession
	}
	var roundOutcome domain.RoundOutcome
	player1move, player2move := s.GetPlayersMoves()

	winMoveMapping := map[domain.Move]domain.Move{
		domain.Rock:    domain.Scissor,
		domain.Paper:   domain.Rock,
		domain.Scissor: domain.Paper,
	}
	if winMoveMapping[player1move] == player2move {
		roundOutcome = domain.Player1Win
	} else {
		roundOutcome = domain.Player2Win
	}
	return roundOutcome, nil
}

func (s *Session) SetRoundOutcome(roundOutcome domain.RoundOutcome) error {
	if !roundOutcome.IsValid() {
		return ErrInvalidRoundOutcome
	}
	s.match.RoundOutcome = roundOutcome
	return nil
}

func (s *Session) SessionAtMaxRound() bool {
	return s.GetRound() >= MaxRound
}


func (s *Session) GetRoundOutCome() domain.RoundOutcome {
	return s.match.RoundOutcome 
}

// SetPlayerName sets players names for a game session
//
// name is the player name to be set
//
// the flag is either 0 for Player1Flag or  1 for Player2Flag
func (s *Session) SetPlayerName(flag int, name string) error {
	switch flag {
	case 0:
		if name == "" {
			name = DefaultPlayer1Name
		}
		s.player1.Name = name
	case 1:
		if name == "" {
			name = DefaultPlayer2Name
		}
		s.player2.Name = name
	default:
		return ErrInvalidFlag
	}
	return nil
}

func (s *Session) GetPlayerName(flag int) (string, error) {
	switch flag {
	case Player1Flag:
		return s.player1.Name, nil
	case Player2Flag:
		return s.player2.Name, nil
	}

	return "", ErrInvalidFlag
}
func (s *Session) SetPlayerMove(flag int, move domain.Move) error {
	switch flag {
	case Player1Flag:
		if !move.IsValid() {
			return ErrInvalidMove
		}
		s.player1.Move = move
	case Player2Flag:
		if !move.IsValid() {
			return ErrInvalidMove
		}
		s.player2.Move = move
	default:
		return ErrInvalidFlag
	}
	return nil
}

// GEtMoves returns player1, and player2 moves
func (s *Session) GetPlayersMoves() (player1Move, player2Move domain.Move) {
	return s.player1.Move, s.player2.Move
}

func (s *Session) IncreasePlayerPoint(flag int) error {
	switch flag {
	case Player1Flag:
		s.player1.Points++
	case Player2Flag:
		s.player2.Points++
	default:
		return ErrInvalidFlag
	}
	return nil

}

func (s *Session) GetPlayersPoints() (player1Point, player2Point int16) {
	return s.player1.Points, s.player2.Points
}


// TODO: write tests for this function
func (s *Session) GetPlayersConn() (player1Conn, player2Conn *websocket.Conn) {
	return s.player1.Connection, s.player2.Connection
}