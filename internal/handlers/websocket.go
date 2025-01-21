package handlers

import (
	"fmt"
	"log"
	"net/http"

	connection_pool "github.com/fishdontexist/chatroom/pkg/connection_pool"
	"github.com/fishdontexist/chatroom/pkg/nats"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Pool      *connection_pool.ConnectionPool
	Publisher *nats.Publisher
}

func New(pool *connection_pool.ConnectionPool, publisher *nats.Publisher) *Handler {
	return &Handler{
		Pool:      pool,
		Publisher: publisher,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request, pool *connection_pool.ConnectionPool) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected successfully!")

	pool.AddConnection(ws)
	reader(ws, pool)

}

func reader(conn *websocket.Conn, pool *connection_pool.ConnectionPool) {
	defer func() {
		// Remove connection from pool when disconnected
		defer conn.Close()
		pool.RemoveConnection(conn)
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
