package main

import (
	"fmt"

	"github.com/fishdontexist/chatroom/internal/api"
)

func main() {
	// Set up router
	router := api.SetupRoutes()

	fmt.Println("Starting WebSocket server on ws://localhost:8080/ws")

}
