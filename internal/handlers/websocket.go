package handlers

import (
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

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected successfully!")

	h.Pool.AddConnection(ws)
	h.reader(ws)

}

func (h *Handler) reader(conn *websocket.Conn) {
	defer func() {
		// Remove connection from pool when disconnected
		defer conn.Close()
		h.Pool.RemoveConnection(conn)
		log.Printf("Client disconnected!")
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		log.Printf("Message received: %s", p)
		// broadcastMessage(messageType, p)
		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	log.Println(err)
		// 	return
		// }

		// Publish message to NATS
		h.Publisher.Publish("chat", string(p))

	}
}
