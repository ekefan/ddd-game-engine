package webapi

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ekefan/ddd-game-engine/internal/core/domain"
	"github.com/ekefan/ddd-game-engine/internal/core/service"
	ports "github.com/ekefan/ddd-game-engine/internal/ports/api"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

type WebServer struct {
	gs ports.GameService
}

func NewWebServer(svc *service.GameService) *WebServer {
	return &WebServer{
		gs: svc,
	}
}

type CreateSessionResp struct {
	SessionID uuid.UUID `json:"session_id"`
}

// TODO: Write test for this handlers
func (ws *WebServer) CreateSession(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("can't wait in lobby", "reason", err)
		return
	}
	player1 := &domain.Player{
		ID:         uuid.New(),
		Name:       domain.DefaultPlayer1Name,
		Connection: socket,
	}
	sess, err := ws.gs.CreateSession(player1)
	if err != nil {
		slog.Error("can't create session", "reason", err)
		return
	}
	resp := CreateSessionResp{
		SessionID: sess.GetID(),
	}
	socket.WriteJSON(resp)
	// TODO: set time out to close connection when another player hasn't joined in a while
}

func (ws *WebServer) PlayGame(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("can't start game session", "reason", err)
		return
	}
	defer socket.Close()
	// validate client request
	urlValues, ok := r.URL.Query()["session_id"]
	if !ok {
		slog.Error("can't play game", "reason", "invalid request parameters")
		socket.WriteMessage(websocket.TextMessage,
			[]byte("Invalid request parameters, should be localhost/game/play?session_id=<session_id you copied"))
		return
	}
	// get session id from request
	sessionID, err := uuid.Parse(urlValues[0])
	if err != nil {
		slog.Error("can't play game", "reason", "invalid request")
		socket.WriteMessage(websocket.TextMessage, []byte("Invalid request, session id is invalid"))
		socket.Close()
		return
	}
	player2 := &domain.Player{
		ID:         uuid.New(),
		Name:       domain.DefaultPlayer2Name,
		Connection: socket,
	}
	err = ws.gs.PlayGame(sessionID, player2)
	if err != nil {
		fmt.Println("error playing game:", err)
	}
}
