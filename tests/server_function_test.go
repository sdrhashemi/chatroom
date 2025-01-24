package tests

import (
	"testing"

	"github.com/fishdontexist/chatroom/internal/config"
	"github.com/fishdontexist/chatroom/pkg/message"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}
	conn, _, err := websocket.DefaultDialer.Dial(cfg.App.ServerURL, nil)
	if err != nil {
		t.Fatalf("Error dialing websocket: %v", err)
	}
	defer conn.Close()

	// username setup
	username := "testuser"
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading message: %v", err)
	}
	assert.Equal(t, "Enter a unique name: ", string(data))

	// send the username to server
	err = conn.WriteMessage(websocket.TextMessage, []byte(username))
	if err != nil {
		t.Fatalf("Error writing message: %v", err)
	}
	_, data, err = conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading message: %v", err)
	}
	assert.Equal(t, "Username Accepted", string(data))
	_, data, err = conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading message: %v", err)
	}
	joindMessage, err := message.Deserialize(data)
	if err != nil {
		t.Fatalf("Error deserializing message: %v", err)
	}
	if joindMessage.Type != "chatroom" {
		t.Fatalf("Expected message type to be chatroom, got: %s", joindMessage.Type)
	}
	assert.Equal(t, "*testuser has joined the chat*", joindMessage.Data)
	// user write a message to the chatroom and he shoudl recieve the same message and acknowledgement
	testUserNewMessage := "Hello, World!"
	err = conn.WriteMessage(websocket.TextMessage, []byte(testUserNewMessage))
	if err != nil {
		t.Fatalf("Error writing message: %v", err)
	}
	_, data, err = conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading message: %v", err)
	}
	receivedServerAck, err := message.Deserialize(data)
	if err != nil {
		t.Fatalf("Error deserializing message: %v", err)
	}
	if receivedServerAck.Type != "ack" {
		t.Fatalf("Expected message type to be ack, got: %s", receivedServerAck.Type)
	}
	assert.Equal(t, "OK", receivedServerAck.Data)
}
