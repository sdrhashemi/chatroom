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

const (
	ConnectionReadDeadline = 30
	AcknowledgmentMessage  = "OK"
	UsernameAccepted       = "Username Accepted"
)

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
		log.Println("Websocket upgrade failed: ", err)
		return
	}

	userName, err := h.captureClientName(ws)
	if err != nil {
		log.Printf("Error capturing client username: %v", err)
		ws.Close()
		return
	}

	h.Pool.AddConnection(userName, ws)
	log.Println("Client connected successfully:", userName)

	// message broadcast to all active users that someone has joind just now
	h.joindMessagePublish(userName)
	h.reader(userName, ws)

}

func (h *Handler) reader(username string, conn *websocket.Conn) {
	defer func() {
		// Remove connection from pool when disconnected
		defer conn.Close()
		h.Pool.RemoveConnection(username, conn)
		// publish the message that someone has left
		h.leaveMessagePublish(username)
		log.Printf("Client disconnected!")
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", username, err)
			continue
		}

		log.Printf("Message received from %s: %s", username, string(message))

		// Acknowlegement message revcieved to the client
		err = sendClientAcknowledgment(conn)
		if err != nil {
			log.Printf("Error sending acknowledgment to %s: %v", username, err)
		}

		if string(message) == "#users" {
			err = h.handleUserCommand(conn)
			if err != nil {
				log.Println("Error marshaling or sending the users list to the clinet: ", err)
			}
			continue
		}
		// Publish message to nats_lib
		h.Publisher.Publish("chat", fmt.Sprintf("%s: %s", username, string(message)))

	}
}

func (h *Handler) captureClientName(conn *websocket.Conn) (string, error) {

	for {

		err := sendMessageViaWebsocket(conn, "Enter a unique name: ")
		if err != nil {
			return "", fmt.Errorf("error sending username prompt: %v", err)
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("error reading name: %v", err)
		}

		clientUserName := strings.TrimSpace(string(message))
		if clientUserName == "" {
			err := sendMessageViaWebsocket(conn, "Name cannot be empty.")
			if err != nil {
				return "", fmt.Errorf("failed to send empty name message: %v", err)
			}
			continue
		}

		// check the username to be unique
		exists := h.Pool.UserNameExists(clientUserName)

		if exists {
			err := sendMessageViaWebsocket(conn, "Name already taken.")
			if err != nil {
				return "", fmt.Errorf("failed to send name taken message: %v", err)
			}
			continue

		}

		// name is unique
		err = sendMessageViaWebsocket(conn, UsernameAccepted)

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

func (h *Handler) joindMessagePublish(username string) {
	joindMessage := message.Message{
		Type: "chatroom",
		Data: fmt.Sprintf("*%s has joined the chat*", username),
	}
	serilizedJoindMessage, err := joindMessage.Serialize()
	if err != nil {
		log.Printf("Error serializing join message: %v", err)
	}
	h.Publisher.Publish("chat", string(serilizedJoindMessage))

}

func (h *Handler) leaveMessagePublish(username string) {

	leaveMessage := message.Message{
		Type: "chatroom",
		Data: fmt.Sprintf("*%s has left the chat*", username),
	}
	finalResponse, err := leaveMessage.Serialize()
	if err != nil {
		log.Print("Error serializing leave message: %v", err)
	}
	h.Publisher.Publish("chat", string(finalResponse))
}

func (h *Handler) handleUserCommand(conn *websocket.Conn) error {
	users := h.Pool.GetUserNames()
	response := map[string]interface{}{
		"type": "users",
		"data": users,
	}
	jsonUsersData, err := json.Marshal(response)
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, jsonUsersData); err != nil {
		return err
	}
	return nil
}

func usernamePrompt(conn *websocket.Conn) error {
	prompt := map[string]string{
		"type": "username",
		"data": "enter name",
	}
	jsonPrompt, err := json.Marshal(prompt)
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, jsonPrompt); err != nil {
		return err
	}
	return nil
}
func sendClientAcknowledgment(conn *websocket.Conn) error {
	okResponse := map[string]string{
		"type": "ack",
		"data": AcknowledgmentMessage,
	}
	jsonOkResponse, err := json.Marshal(okResponse)
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, jsonOkResponse); err != nil {
		return err
	}
	return nil
}

func sendMessageViaWebsocket(conn *websocket.Conn, message string) error {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return err
	}
	return nil
}
