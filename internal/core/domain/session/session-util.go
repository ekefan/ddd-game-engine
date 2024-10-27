package session

import (
	"fmt"
	"strings"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
)

func parseMove(msg []byte) (domain.Move, error) {
	switch strings.ToLower(string(msg)) {
	case "rock":
		return domain.Rock, nil
	case "paper":
		return domain.Paper, nil
	case "scissor":
		return domain.Scissor, nil
	default:
		return -1, fmt.Errorf("invalid move: %s", msg)
	}
}

func (s *Session) gameEnded() bool {
    if s.match.Round >= MaxRound && s.player1.Points != s.player2.Points {
		s.sessionEnded = true
        return true
    }
    return false
}


func (s *Session) setPlayerPoints() {
	if s.player1.Move == s.player2.Move {
		s.match.RoundOutcome = domain.Draw
		return
	}
	s.updateRound()
	moveMap := map[domain.Move]domain.Move{
		domain.Rock:    domain.Scissor,
		domain.Paper:   domain.Rock,
		domain.Scissor: domain.Paper,
	}
	if moveMap[s.player1.Move] == s.player2.Move {
		s.player1.Points++
		s.match.RoundOutcome = domain.Player1Win
	} else {
		s.player2.Points++
		s.match.RoundOutcome = domain.Player2Win
	}
}

func (s *Session) generateResponse() {
	s.response = &domain.Response{
		Round:  s.match.Round,
		Winner: s.getRoundWinner(),
		RoundOutcome: s.getRoundOutcome(),
		SessionEnded: s.sessionEnded,
	}
}


func (s *Session) getRoundOutcome() string {
    return s.match.RoundOutcome.String()
}

func (s *Session) getRoundWinner() string {
    switch s.match.RoundOutcome {
    case domain.Draw:
        return "no winner"
    case domain.Player1Win:
        return s.player1.Name
    default:
        return s.player2.Name
    }
}

func (s *Session) updateRound() {
	s.match.Round++
}


func (s *Session) isSessionEnded() bool{
	s.mu.Lock()
	ended := s.sessionEnded
	s.mu.Unlock()
	return ended
}

func (s *Session) getSessionWinner() string {
	if s.player1.Points > s.player2.Points {
		return s.player1.Name
	}
	return s.player2.Name
}