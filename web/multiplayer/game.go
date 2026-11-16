package multiplayer

import (
	"log/slog"
	"math/rand"
	"net/http"
	"uuid"

	"github.com/gorilla/websocket"
	"github.com/mifumofu/tetoris/game"
)

type GameID string
type MpGame struct {
	id          GameID
	state       State
	MaxPlayers  int
	Players     map[client]bool
	InitialGame game.Game
}

type State int8

const (
	StateLobby State = iota
	StatePlaying
	StateEnd
)

type client struct {
	ID      uuid.UUID
	send    chan []byte
	conn    *websocket.Conn
	request *http.Request
	logger  *slog.Logger
}

func (g *MpGame) DisconnectClient(c client) {
	c.conn.Close()
	delete(g.Players, c)

	if g.state == StateLobby {
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

func NewGame() (GameID, *MpGame) {
	return newGameId(), new(MpGame{
		MaxPlayers:  2,
		Players:     make(map[client]bool),
		state:       StateLobby,
		InitialGame: game.NewGame(),
	})
}

// sends a message to all players in a game
func (g *MpGame) sendAll(messageType int, data []byte) {
	for p := range g.Players {
		p.send <- data
	}
}
