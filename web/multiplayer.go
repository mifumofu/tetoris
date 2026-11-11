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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

var clients = make(map[client]bool) // Track active clients

type client struct {
	id   int
	conn *websocket.Conn
}

func handleMpConnections(w http.ResponseWriter, r *http.Request) {
	// limit to 2 people in a game
	if len(clients) >= 2 {
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

	clients[c] = true

	// stop sending messages to client when they disconnect
	defer func() {
		delete(clients, c)
		ws.Close()
	}()

	fmt.Println("num of clients", len(clients))

	// handle incoming messages
	handleMpMessages(ws, c)
}

func handleMpMessages(ws *websocket.Conn, c client) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			delete(clients, c)
			break
		}

		a, err := parseMpMessage(msg)
		if err != nil {
			fmt.Println("failed to parse message:", err)
			delete(clients, c)
			break
		}

		fmt.Printf("%v: %v\n", c.id, a.String())

		// Broadcast message to all clients
		for client := range clients {
			if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("broadcast error:", err)
				client.conn.Close()
				delete(clients, c)
			}
		}
	}
}

var ErrInvalidMessage = errors.New("Invalid message sent!")

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
