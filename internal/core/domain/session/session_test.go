package session_test

import (
	"math/rand"
	"testing"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// randomValidRoundOutcome generates a random valid round outcome
func randomValidRoundOutcome() domain.RoundOutcome {
	possibleOutcomes := [3]domain.RoundOutcome{domain.Draw, domain.Player1Win, domain.Player2Win}
	randIdx := rand.Intn(len(possibleOutcomes))
	return possibleOutcomes[randIdx]
}

func randomValidMove() domain.Move {
	possibleMoves := [3]domain.Move{
		domain.Paper,
		domain.Rock,
		domain.Scissor,
	}
	randIdx := rand.Intn(len(possibleMoves))
	return possibleMoves[randIdx]
}
func createNewSession(t *testing.T) session.Session {
	gs := session.NewSession(&domain.Player{})
	require.NotEmpty(t, gs)
	assert.Equal(t, gs.GetRound(), session.InitRound)
	return *gs
}

func TestNewSession(t *testing.T) {
	createNewSession(t)
}

func TestGetID(t *testing.T) {
	gs := createNewSession(t)
	require.NotEmpty(t, gs.GetID())
}

func TestGetRound(t *testing.T) {
	gs := createNewSession(t)
	require.Equal(t, gs.GetRound(), session.InitRound)
}