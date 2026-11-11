package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handleMpConnections)

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
