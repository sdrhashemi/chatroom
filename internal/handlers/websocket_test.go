package handlers_test

import (
	"testing"

	"github.com/fishdontexist/chatroom/internal/app/server"
	"github.com/fishdontexist/chatroom/internal/config"
	"github.com/gorilla/websocket"
)

func TestWeboSocketHandler(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	testChatApp, err := server.New(cfg.Nats.URL)
	if err != nil {
		t.Fatalf("Error creating chat app: %v", err)
	}
	testChatApp.SetupRoutes()

	conn, _, err := websocket.DefaultDialer.Dial(cfg.App.ServerURL, nil)
	if err != nil {
		t.Fatalf("Failed to connecto to Websocket server: %v", err)
	}
	defer conn.Close()

	// client sending a message
	err = conn.WriteMessage(websocket.TextMessage, []byte("test"))
	if err != nil {
		t.Fatalf("Failed to write message to websocket server: %v", err)
	}

	_, response, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read server response: %v", err)
	}

	if string(response) != "OK" {
		t.Errorf("Expected 'OK', got '%s'", string(response))
	}
}
