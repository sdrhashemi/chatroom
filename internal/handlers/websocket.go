package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Connection pool for managing users
var activeConnections = struct {
	sync.Mutex
	pool []*websocket.Conn
}{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected successfully!")

	activeConnections.Lock()
	activeConnections.pool = append(activeConnections.pool, ws)
	activeConnections.Unlock()
	reader(ws)

}

func reader(conn *websocket.Conn) {
	defer func() {
		// Remove connection from pool when disconnected
		defer conn.Close()
		activeConnections.Lock()
		for i, connection := range activeConnections.pool {
			if connection == conn {
				activeConnections.pool = append(activeConnections.pool[:i], activeConnections.pool[i+1:]...)
				break
			}
		}
		activeConnections.Unlock()
		log.Printf("Client disconnected!")
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))
		// broadcastMessage(messageType, p)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
