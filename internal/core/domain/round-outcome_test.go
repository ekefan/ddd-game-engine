package domain

import (
	"testing"
)

func TestRoundOutcomeString(t *testing.T) {
	var outcome RoundOutcome
	possibleOutcomes := [3]string{"Draw", "Player1 win", "Player2 win"}
	for range 3 {
		if possibleOutcomes[outcome] != outcome.String() {
			t.Errorf("invalid round outcome: %d", outcome)
		}
		outcome++
	}
}


func TestRoundOutcomeIsValid(t *testing.T) {
	type testCase struct {
		name string
		roundOutcome RoundOutcome
		expectedResult bool
	}

	testCases := []testCase{
		{
			name: "valid case",
			roundOutcome: Draw,
			expectedResult: true,
		}, {
			name: "lower invalid case",
			roundOutcome: -1,
			expectedResult: false,
		}, {
			name: "upper invalid case",
			roundOutcome: 3,
			expectedResult: false,
		}, {
			name: "second valid case",
			roundOutcome: Player1Win,
			expectedResult: true,
		},{
			name: "last valid case",
			roundOutcome: Player2Win,
			expectedResult: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			res := tc.roundOutcome.IsValid()
			if res != tc.expectedResult{
				t.Errorf("expected %v, got %v", tc.expectedResult, res)
			}
			
		})
	}
}