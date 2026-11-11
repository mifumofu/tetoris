package main

import (
	"log/slog"
	"net/http"
	"os"
)

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

	http.HandleFunc("/mp/{id}/ws", handleMpConnections)
	http.HandleFunc("/mp/create", handleCreateGame)

	addr := ":8080"
	slog.Info("Server started", "addr", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("ListenAndServe", "err", err)
	}
}
