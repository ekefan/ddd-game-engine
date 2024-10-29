# ddd-game-engine

###### _project is not complete_
The project features a domain driven game engine for a rock paper and scissor game, serving two applications a http server and a cli server at the same time

Currently only the http server is available and can be spinned of on your local machine by runing this command when you clone this repository to your machine:

```bash
go run cmd/web/main.go
```

The application exposes an interface for the game engine to the two applications using principles from hexagonal architecture.

## The game play

- Two players are connected to the server using web sockets
- A player must initiate the connection by sending a request to the game endpoint, that player receives a session id with which the second player can join before the game starts
- The game server hosts the game service which manages game sessions between two connected players
- Each player takes turns to make moves and at each round the results of the round is given back to the players
- At the end of the game which happens after at least 3 valid rounds the total result of the game is given back to the players and the game services ends the session.

## Future Updates

- CLI application to fully support synchronization with web application
- Player customization
- Good unit tests
