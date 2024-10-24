package rps

// RoundOutcome is a value-object that represents the outcome of a round
// 
// Draw == 0
//
// Player1Win == 1
// 
// Player2Win == 2
//
// Any other Round outcome is considered invalid, these outcomes have a string representation
type RoundOutcome int

const (
	Draw RoundOutcome = iota
	Player1Win
	Player2Win
)

func (ro RoundOutcome) String() string {
	return [3]string{"Draw", "Player1 win", "Player2 win"}[ro]
}

// IsValid checks if round outcome is a valid i.e either draw, player1win or player2win
func (ro RoundOutcome) IsValid() bool {
	return ro == Draw || ro == Player1Win || ro == Player2Win
}