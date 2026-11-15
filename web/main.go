package main

import (
	"log/slog"
	"net/http"
	"os"
	"uuid"

	"github.com/gorilla/websocket"
	"github.com/mifumofu/tetoris/game"
)

func main() {

	if os.Getenv("DEBUG") != "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Info("Enabled debug logging")
	} else {
		// use json handler unless debugging
		slog.SetDefault(
			slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		)
	}

	http.HandleFunc("/mp/{id}/ws", handleMpConnections)
	http.HandleFunc("/mp/create", handleCreateGame)

	addr := ":8080"
	slog.Info("Server started", "addr", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("ListenAndServe", "err", err)
	}
}

var games = make(map[GameID]*MpGame)

type GameID string
type MpGame struct {
	id          GameID
	state       state
	MaxPlayers  int
	Players     map[client]bool
	InitialGame game.Game
}

type state int8

const (
	StateLobby state = iota
	StatePlaying
	StateEnd
)

type client struct {
	ID      uuid.UUID
	conn    *websocket.Conn
	request *http.Request
	logger  *slog.Logger
}
