package multiplayer

import (
	"errors"
	"log/slog"
	"net/http"
	"uuid"

	"github.com/gorilla/websocket"
)

var (
	ErrInvalidMessage = errors.New("Invalid message sent!")
)

// game over websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

func ConnectWS(g *MpGame, w http.ResponseWriter, r *http.Request) {
	// limit to 2 people in a game
	if len(g.Players) >= g.MaxPlayers {
		slog.Warn("Rejecting connection as game at max players", "ip", r.RemoteAddr)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// do not let player join if not in lobby
	if g.state != StateLobby {
		slog.Warn("Rejecting connection as game not in lobby", "ip", r.RemoteAddr)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// upgrade client's request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("Failed to upgrade connection to websocket", "ip", r.RemoteAddr)
		return
	}

	cid := uuid.NewV4()
	c := client{
		conn:    ws,
		ID:      cid,
		request: r,
		logger: slog.With(slog.Group("conn",
			"ip", r.RemoteAddr,
			"user-agent", r.UserAgent(),
			"id", cid,
		)),
	}
	g.Players[c] = true

	// stop sending messages to client when they disconnect
	defer func() {
		delete(g.Players, c)
		ws.Close()
	}()

	// send update in game state to players
	g.sendLobbyUpdate(false)

	// handle incoming messages
	handleMpMessages(ws, c, g)
}

func handleMpMessages(ws *websocket.Conn, c client, g *MpGame) {
	c.logger = slog.With(slog.Group("conn",
		"ip", c.request.RemoteAddr,
		"user-agent", c.request.UserAgent(),
		"id", c.ID,
	))

	c.logger.Info("Connected")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.logger.Warn("failed to read message", "err", err)
			g.DisconnectClient(c)
			break
		}

		var cont bool
		switch g.state {
		case StateLobby:
			cont = handleMpLobbyMessages(msg, c, g)
		case StatePlaying:
			cont = handleMpPlayingMsgs(msg, c, g)
		}

		if cont == false {
			break
		}
	}
}
