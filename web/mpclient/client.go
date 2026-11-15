package mpclient

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/mifumofu/tetoris/game"
)

// server
var Server = "127.0.0.1:8080"

type MpGame struct {
	game game.Game
	done chan struct{}
	conn *websocket.Conn
}

// creates a new game and joins it
func Create() (MpGame, error) {
	code, err := create()
	if err != nil {
		return MpGame{}, err
	}

	return Connect(code)
}

// creates a new game
func create() (string, error) {
	u := url.URL{Scheme: "https", Host: Server, Path: "/mp/create"}
	res, err := http.Post(u.String(), "text/plain", nil)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(res.Body)
	return string(b), err
}

// connects to a game
func Connect(code string) (MpGame, error) {
	// connection url
	u := url.URL{Scheme: "ws", Host: Server, Path: "/mp/" + code + "/ws"}
	log.Printf("connecting to %s", u.String())

	// dial websocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return MpGame{}, err
	}
	done := make(chan struct{})

	// deal with messages
	g := MpGame{
		conn: conn,
		done: done,
	}

	g.game = game.NewGame(game.WithActionListener(g.actionListener))

	go g.receive()

	return g, nil
}

// receives messages
func (g MpGame) receive() {
	defer close(g.done)
	for {
		_, message, err := g.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}

// sends an action to the server
func (g MpGame) Send(a game.Action) {
	err := g.conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(int(a))))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

// disconnects from the server/game
func (g MpGame) Disconnect() {
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
}
