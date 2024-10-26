package domain

import (
	"testing"
)

func TestMoveString(t *testing.T) {
	var move Move
	possibleMoves := [3]string{"Rock", "Paper", "Scissor"}
	for range 3 {
		if possibleMoves[move] != move.String() {
			t.Errorf("invalid move: %d", move)
		}
		move++
	}
}

func TestMoveIsValid(t *testing.T) {
	type testCase struct {
		name string
		move Move
		expectedResult bool
	}

	testCases := []testCase{
		{
			name: "valid case",
			move: Rock,
			expectedResult: true,
		}, {
			name: "lower invalid case",
			move: -1,
			expectedResult: false,
		}, {
			name: "upper invalid case",
			move: 3,
			expectedResult: false,
		}, {
			name: "second valid case",
			move: Paper,
			expectedResult: true,
		},{
			name: "last valid case",
			move: Scissor,
			expectedResult: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			res := tc.move.IsValid()
			if res != tc.expectedResult{
				t.Errorf("expected %v, got %v", tc.expectedResult, res)
			}
			
		})
	}
}