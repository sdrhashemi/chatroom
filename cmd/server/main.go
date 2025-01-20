package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Set up router
	router :=

		// Start the HTTP server
		fmt.Println("Starting WebSocket server on ws://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
