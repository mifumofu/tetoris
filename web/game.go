package main

import (
	"bytes"
	"encoding/gob"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"uuid"

	"github.com/gorilla/websocket"
	"github.com/mifumofu/tetoris/game"
)

var games = make(map[GameID]*MpGame)

type GameID string
type MpGame struct {
	id         GameID
	MaxPlayers int
	Players    map[client]bool
}

type client struct {
	ID      uuid.UUID
	conn    *websocket.Conn
	request *http.Request
	logger  *slog.Logger
}

func (g *MpGame) DisconnectClient(c client) {
	c.conn.Close()
	delete(g.Players, c)
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

// sending data
type MpUpdate struct {
	PlayerID uuid.UUID
	Action   game.Action
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

// sends an MpUpdate to all players in a game
func (g *MpGame) sendAllUpdate(id uuid.UUID, action game.Action) {
	g.sendAll(websocket.BinaryMessage, encodeAction(id, action))
}

// encodes an action for sending over a websocket
// returns encoding in bytes
func encodeAction(id uuid.UUID, action game.Action) []byte {
	var buf bytes.Buffer
	encodeActionWriter(&buf, id, action)
	return buf.Bytes()
}

// encodes an action for sending over a websocket
// uses an io.Writer
func encodeActionWriter(w io.Writer, id uuid.UUID, action game.Action) {
	enc := gob.NewEncoder(w)
	enc.Encode(MpUpdate{id, action})
}
