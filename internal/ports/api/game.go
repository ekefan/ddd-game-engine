package ports

import "context"

type GameServer interface {
	StartSession(ctx context.Context) error
	AddPlayer(ctx context.Context) error
	MakeMove(ctx context.Context) error
}
