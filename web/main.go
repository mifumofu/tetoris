package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/mp/{id}/ws", handleMpConnections)
	http.HandleFunc("/mp/create", handleCreateGame)

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
