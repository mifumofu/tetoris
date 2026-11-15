package main

import (
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/mifumofu/tetoris/game"
)

func (g *MpGame) DisconnectClient(c client) {
	c.conn.Close()
	delete(g.Players, c)

	if g.state == StateLobby {
		time.Sleep(50 * time.Millisecond) // avoids: panic: concurrent write to websocket connection
		slog.Info("Sending lobby update")
		g.sendLobbyUpdate(false)
	}
}

// returns a string of length l with its contents being random
// alphanumeric characters
func alphanumericRand(l int) string {
	gameIdCodeChars := []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'o', 'm', 'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '7', '8', '9', '7', '8', '9',
	}

	var code string
	for range l {
		code += string(gameIdCodeChars[rand.Intn(len(gameIdCodeChars))])
	}
	return code
}

func newGameId() GameID {
	return GameID(alphanumericRand(6))
}

func newGame() (GameID, *MpGame) {
	return newGameId(), new(MpGame{
		MaxPlayers:  2,
		Players:     make(map[client]bool),
		state:       StateLobby,
		InitialGame: game.NewGame(),
	})
}

// create new game
func handleCreateGame(w http.ResponseWriter, r *http.Request) {
	i, g := newGame()

	slog.Info("Created new game", "id", i)
	games[i] = g
	w.Write([]byte(i))
}

// sends a message to all players in a game
func (g *MpGame) sendAll(messageType int, data []byte) {
	for p := range g.Players {
		go func() {
			err := p.conn.WriteMessage(messageType, data)
			if err != nil {
				p.logger.Warn("broadcast error", "err", err)
				g.DisconnectClient(p)
			}
		}()
	}
}
