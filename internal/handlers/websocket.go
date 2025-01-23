package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	connection_pool "github.com/fishdontexist/chatroom/pkg/connection_pool"
	"github.com/fishdontexist/chatroom/pkg/message"
	nats_lib "github.com/fishdontexist/chatroom/pkg/nats"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

const ConnectionReadDeadline = 30

type Handler struct {
	Pool      *connection_pool.ConnectionPool
	Publisher *nats_lib.Publisher
}

func New(pool *connection_pool.ConnectionPool, publisher *nats_lib.Publisher) *Handler {
	h := &Handler{
		Pool:      pool,
		Publisher: publisher,
	}
	h.subscribeToNATS()
	return h
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

	userName, err := h.captureClientName(ws)
	if err != nil {
		log.Printf("Error capturing client username: %v", err)
		ws.Close()
		return
	}

	h.Pool.AddConnection(userName, ws)

	log.Println("Client connected successfully!")
	joindMessage := message.Message{
		Type: "chatroom",
		Data: fmt.Sprintf("*%s has joined the chat*", userName),
	}
	serilizedJoindMessage, err := joindMessage.Serialize()
	if err != nil {
		log.Printf("Error serializing join message: %v", err)
	}
	h.Publisher.Publish("chat", string(serilizedJoindMessage))
	h.reader(userName, ws)

}

func (h *Handler) reader(username string, conn *websocket.Conn) {
	defer func() {
		// Remove connection from pool when disconnected
		defer conn.Close()
		h.Pool.RemoveConnection(username, conn)
		leaveMessage := message.Message{
			Type: "chatroom",
			Data: fmt.Sprintf("*%s has left the chat*", username),
		}
		finalResponse, err := leaveMessage.Serialize()
		if err != nil {
			log.Print("Error serializing leave message: %v", err)
		}
		h.Publisher.Publish("chat", string(finalResponse))
		log.Printf("Client disconnected!")
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", username, err)
			continue
		}

		log.Printf("Message received from %s: %s", username, string(message))

		if string(message) == "#users" {
			users := h.Pool.GetUserNames()
			response := map[string]interface{}{
				"type": "users",
				"data": users,
			}
			jsonUsersData, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshaling users list to JSON: %v", err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, jsonUsersData); err != nil {
				log.Printf("Error sending users list to client: %v", err)
			}
			continue
		}
		// Publish message to nats_lib
		h.Publisher.Publish("chat", fmt.Sprintf("%s: %s", username, string(message)))

	}
}

func (h *Handler) captureClientName(conn *websocket.Conn) (string, error) {

	for {

		err := conn.WriteMessage(websocket.TextMessage, []byte("Please enter a unique name: "))
		if err != nil {
			return "", fmt.Errorf("failed to prompt client for name: %v", err)

		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("error reading name: %v", err)
		}

		clientUserName := strings.TrimSpace(string(message))
		if clientUserName == "" {
			err := conn.WriteMessage(websocket.TextMessage, []byte("Name cannot be empty. enter your name again: "))
			if err != nil {
				return "", fmt.Errorf("failed to send empty name message: %v", err)
			}
			continue
		}

		// check the username to be unique
		exists := h.Pool.UserNameExists(clientUserName)

		if exists {
			err = conn.WriteMessage(websocket.TextMessage, []byte("Username already exists, enter another name: "))
			if err != nil {
				return "", fmt.Errorf("failed to send duplicate name message: %v", err)
			}
			continue

		}

		// name is unique
		err = conn.WriteMessage(websocket.TextMessage, []byte("OK"))

		return clientUserName, nil

	}
}

func (h *Handler) subscribeToNATS() {
	_, err := h.Publisher.NATSConnection().Subscribe("chat", func(msg *nats.Msg) {
		log.Printf("Message received from nats_lib: %s", string(msg.Data))

		// Broadcast to all connected users
		h.Pool.BroadcastMessageToClients(msg.Data)
	})
	if err != nil {
		log.Fatalf("Error subscribing to nats_lib: %v", err)
	}
}
