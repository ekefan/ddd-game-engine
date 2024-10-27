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
