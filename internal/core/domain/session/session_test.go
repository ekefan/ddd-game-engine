package session_test

import (
	"math/rand"
	"testing"

	"github.com/ekefan/ddd-game-engine/internal/core/domain/session"
	"github.com/ekefan/ddd-game-engine/internal/core/rps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// randomValidRoundOutcome generates a random valid round outcome
func randomValidRoundOutcome() rps.RoundOutcome {
	possibleOutcomes := [3]rps.RoundOutcome{rps.Draw, rps.Player1Win, rps.Player2Win}
	randIdx := rand.Intn(len(possibleOutcomes))
	return possibleOutcomes[randIdx]
}

func randomValidMove() rps.Move {
	possibleMoves := [3]rps.Move{
		rps.Paper,
		rps.Rock,
		rps.Scissor,
	}
	randIdx := rand.Intn(len(possibleMoves))
	return possibleMoves[randIdx]
}
func createNewSession(t *testing.T) session.Session {
	gs, err := session.NewSession()
	require.NoError(t, err)
	require.NotEmpty(t, gs)
	return gs
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

func TestUpdateRound(t *testing.T) {
	gs := createNewSession(t)
	initRound := gs.GetRound()
	gs.UpdateRound()
	updatedRound := gs.GetRound()
	require.Equal(t, 1, updatedRound - initRound)
}
func TestSetRoundOutcome(t *testing.T) {
	type testCase struct {
		name         string
		roundOutcome rps.RoundOutcome
		expectedErr  error
	}
	testCases := []testCase{
		{
			name:         "first valid roundOutcome",
			roundOutcome: randomValidRoundOutcome(),
			expectedErr:  nil,
		}, {
			name:         "second valid roundOutcome",
			roundOutcome: randomValidRoundOutcome(),
			expectedErr:  nil,
		}, {
			name:         "third valid roundOutcome",
			roundOutcome: randomValidRoundOutcome(),
			expectedErr:  nil,
		}, {
			name:         "lower invalid roundOutcome",
			roundOutcome: -1,
			expectedErr:  session.ErrInvalidRoundOutcome,
		}, {
			name:         "upper invalid roundOutcome",
			roundOutcome: 3,
			expectedErr:  session.ErrInvalidRoundOutcome,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gs := createNewSession(t)
			err := gs.SetRoundOutcome(tc.roundOutcome)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSetPlayerName(t *testing.T) {
	type testCase struct {
		expectedErr  error
		expectedName string
		test         string
		name         string
		flag         int
	}
	testCases := []testCase{
		{
			test:         "empty name string for player 1",
			name:         "",
			flag:         session.Player1Flag,
			expectedErr:  nil,
			expectedName: session.DefaultPlayer1Name,
		}, {
			test:         "non empty name string",
			name:         "eben",
			flag:         session.Player2Flag,
			expectedErr:  nil,
			expectedName: "eben",
		}, {
			test:         "invalid flag",
			name:         "",
			flag:         -1,
			expectedErr:  session.ErrInvalidFlag,
			expectedName: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			gs := createNewSession(t)
			err := gs.SetPlayerName(tc.flag, tc.name)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}


func TestSetPlayerMove(t *testing.T) {
	type testCase struct {
		name string
		move rps.Move
		flag int
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "set valid move 1",
			move: randomValidMove(),
			flag: session.Player1Flag,
			expectedErr: nil,
		},{
			name: "set lower invalid move",
			move: -1,
			flag: session.Player1Flag,
			expectedErr: session.ErrInvalidMove,
		},{
			name: "set valid move 2",
			move: randomValidMove(),
			flag: session.Player2Flag,
			expectedErr: nil,
		},{
			name: "set upper invalid move",
			move: 3,
			flag: session.Player2Flag,
			expectedErr: session.ErrInvalidMove,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.name, func(t *testing.T){
			gs := createNewSession(t)
			err := gs.SetPlayerMove(tc.flag, tc.move)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}