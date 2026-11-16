package main

import (
	_ "embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/mifumofu/tetoris/web/multiplayer"
)

//go:embed index.html
var indexFile []byte

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

	http.HandleFunc("GET /{$}", handleRoot)
	http.HandleFunc("/mp/{id}/ws", handleMpConnections)
	http.HandleFunc("/mp/create", handleCreateGame)

	addr := ":8080"
	slog.Info("Server started", "addr", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("ListenAndServe", "err", err)
	}
}

var games = make(map[multiplayer.GameID]*multiplayer.MpGame)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(indexFile)
}

// create new game
func handleCreateGame(w http.ResponseWriter, r *http.Request) {
	i, g := multiplayer.NewGame()

	slog.Info("Created new game", "id", i)
	games[i] = g
	w.Write([]byte(i))
}

func handleMpConnections(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	g, exists := games[multiplayer.GameID(id)]
	if !exists {
		slog.Warn("Non existant game connect attempt", "ip", r.RemoteAddr)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	multiplayer.ConnectWS(g, w, r)
}
