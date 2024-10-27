package main

import (
	"net/http"

	"github.com/ekefan/ddd-game-engine/internal/adapters/handlers/webapi"
	"github.com/ekefan/ddd-game-engine/internal/adapters/memory"
	"github.com/ekefan/ddd-game-engine/internal/core/service"
)

func main() {
	sessionRepo := memory.NewSessionRepository()
	gameService := service.NewGameService(sessionRepo)

	ws := webapi.NewWebServer(gameService)

	http.HandleFunc("/game", ws.CreateSession)
	http.HandleFunc("/game/play", ws.PlayGame)
	http.ListenAndServe(":8080", nil)
}
