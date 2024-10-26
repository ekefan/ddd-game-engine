package handler

import (
	// "github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/ekefan/ddd-game-engine/internal/core/service"
	ports "github.com/ekefan/ddd-game-engine/internal/ports/api"
)

type WebServer struct {
	gsvc ports.GameService
}

func NewHttpServer(svc *service.GameService) *WebServer {
	return &WebServer{
		gsvc: svc,
	}
}