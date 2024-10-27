package domain

import "github.com/gorilla/websocket"

// Move is a value-object that represents a move a player can make
//
// Rock == 0
//
// Paper == 1
//
// Scissor == 2
//
// Any other Move is considered invalid, these outcomes have a string representation
type Move int

const (
	Rock Move = iota
	Paper
	Scissor
)

func (m Move) String() string {
	return [...]string{"Rock", "Paper", "Scissor"}[m]
}

// IsValid checks if a move is valid, i.e one of Rock, Paper, Scissors
func (m Move) IsValid() bool {
	return m == Rock || m == Paper || m == Scissor
}

type PlayerMove struct {
	Move Move
	Conn *websocket.Conn
}