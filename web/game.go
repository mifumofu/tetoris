package main

import (
	"log/slog"
	"math/rand"
	"net/http"
)

var games = make(map[GameID]*MpGame)

type GameID string
type MpGame struct {
	MaxPlayers int
	Players    map[client]bool
}

func (g *MpGame) DisconnectClient(c client) {
	c.conn.Close()
	delete(g.Players, c)
}

func newGameId() GameID {
	gameIdCodeChars := []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'o', 'm', 'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '7', '8', '9', '7', '8', '9',
	}

	var code string
	for range 6 {
		code += string(gameIdCodeChars[rand.Intn(len(gameIdCodeChars))])
	}
	return GameID(code)
}

func newGame() (GameID, *MpGame) {
	return newGameId(), new(MpGame{
		MaxPlayers: 2,
		Players:    make(map[client]bool),
	})
}

// create new game
func handleCreateGame(w http.ResponseWriter, r *http.Request) {
	i, g := newGame()

	slog.Info("Created new game", "id", i)
	games[i] = g
	w.Write([]byte(i))
}
