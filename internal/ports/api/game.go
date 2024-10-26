package ports

type GameService interface {
	WaitInLobby() error //start route post
	PlayGame() error    //play route web socket
}
