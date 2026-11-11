package main

import (
	"errors"
	"fmt"
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
	id   int
	conn *websocket.Conn
}

func handleMpConnections(w http.ResponseWriter, r *http.Request) {
	// get game provided in websocket
	id := r.PathValue("id")
	g, exists := games[GameID(id)]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// limit to 2 people in a game
	if len(g.Players) >= g.MaxPlayers {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// upgrade client's request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := client{
		conn: ws,
		id:   rand.Int(),
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
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			g.DisconnectClient(c)
			break
		}

		a, err := parseMpMessage(msg)
		if err != nil {
			fmt.Println("failed to parse message:", err)
			g.DisconnectClient(c)
			break
		}

		fmt.Printf("%v: %v\n", c.id, a.String())

		// Broadcast message to all clients
		for client := range g.Players {
			if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("broadcast error:", err)
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
