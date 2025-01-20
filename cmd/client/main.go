package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

func main() {
	// Connect to WebSocket server
	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket server")

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket Server!"))
	if err != nil {
		log.Fatal("WriteMessage failed:", err)
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("ReadMessage failed:", err)
		return
	}
	fmt.Printf("Received from server: %s\n", msg)
}
