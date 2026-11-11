package main

import (
	"errors"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/mifumofu/tetoris/game"
)

var (
	ErrInvalidMessage = errors.New("Invalid message sent!")
)

// game over websocket

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

type client struct {
	id      int
	conn    *websocket.Conn
	request *http.Request
}

func handleMpConnections(w http.ResponseWriter, r *http.Request) {
	// get game provided in websocket
	id := r.PathValue("id")
	g, exists := games[GameID(id)]
	if !exists {
		slog.Warn("Non existant game connect attempt", "ip", r.RemoteAddr)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// limit to 2 people in a game
	if len(g.Players) >= g.MaxPlayers {
		slog.Warn("Rejecting connection as game at max players", "ip", r.RemoteAddr)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// upgrade client's request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("Failed to upgrade connection to websocket", "ip", r.RemoteAddr)
		return
	}

	c := client{
		conn:    ws,
		id:      rand.Int(),
		request: r,
	}
	g.Players[c] = true

	// stop sending messages to client when they disconnect
	defer func() {
		delete(g.Players, c)
		ws.Close()
	}()

	// handle incoming messages
	handleMpMessages(ws, c, g)
}

func handleMpMessages(ws *websocket.Conn, c client, g *MpGame) {
	logger := slog.With(slog.Group("conn",
		"ip", c.request.RemoteAddr,
		"user-agent", c.request.UserAgent(),
		"id", c.id,
	))

	logger.Info("Connected")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			logger.Warn("failed to read message", "err", err)
			g.DisconnectClient(c)
			break
		}

		a, err := parseMpMessage(msg)
		if err != nil {
			logger.Warn("failed to parse message", "err", err)
			g.DisconnectClient(c)
			break
		}

		logger.Debug("msg", "action", a.String())

		// Broadcast message to all clients
		for client := range g.Players {
			if err := client.conn.WriteMessage(websocket.TextMessage, []byte(a.String())); err != nil {
				logger.Warn("broadcast error", "err", err)
				g.DisconnectClient(c)
			}
		}
	}
}

// clients should be sending the int of the game.Actions enum action
// they are performing in their game
func parseMpMessage(b []byte) (game.Action, error) {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return game.ActionDown, ErrInvalidMessage
	}

	a := game.Action(i)

	// check if valid
	// TODO: find a better way to see if action is valid
	if i > int(game.ActionInstadown) || i < 0 {
		return a, ErrInvalidMessage
	}

	return a, nil
}
