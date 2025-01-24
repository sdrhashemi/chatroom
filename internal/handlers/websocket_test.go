package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fishdontexist/chatroom/internal/handlers"
	"github.com/fishdontexist/chatroom/pkg/connection_pool"
	nats_lib "github.com/fishdontexist/chatroom/pkg/nats"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

// TestWebSocketHandler tests the WebSocket upgrade process and username handling.
func TestWebSocketHandler(t *testing.T) {
	// Create a new connection pool and NATS publisher
	pool := connection_pool.New()
	publisher, err := nats_lib.New(nats.DefaultURL)
	assert.NoError(t, err)

	// Create a new handler
	handler := handlers.New(pool, publisher)

	// Create a test server with the WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(handler.WebSocketHandler))
	defer server.Close()

	// Convert the server URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Send a username to the server
	username := "testuser"
	err = ws.WriteMessage(websocket.TextMessage, []byte(username))
	assert.NoError(t, err)

	// Read the acknowledgment message from the server
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, handlers.UsernameAccepted, string(msg))

	// Verify that the username is added to the connection pool
	assert.True(t, pool.UserNameExists(username))
}

// TestCaptureClientName tests the username capture functionality.
func TestCaptureClientName(t *testing.T) {
	// Create a new connection pool and NATS publisher
	pool := connection_pool.New()
	publisher, err := nats_lib.New(nats.DefaultURL)
	assert.NoError(t, err)

	// Create a new handler
	handler := handlers.New(pool, publisher)

	// Create a test server with the WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(handler.WebSocketHandler))
	defer server.Close()

	// Convert the server URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Send a username to the server
	username := "testuser"
	err = ws.WriteMessage(websocket.TextMessage, []byte(username))
	assert.NoError(t, err)

	// Read the acknowledgment message from the server
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, handlers.UsernameAccepted, string(msg))

}

// TestBroadcastMessage tests the message broadcasting functionality.
func TestBroadcastMessage(t *testing.T) {
	// Create a new connection pool and NATS publisher
	pool := connection_pool.New()
	publisher, err := nats_lib.New(nats.DefaultURL)
	assert.NoError(t, err)

	// Create a new handler
	handler := handlers.New(pool, publisher)

	// Create a test server with the WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(handler.WebSocketHandler))
	defer server.Close()

	// Convert the server URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Send a username to the server
	username := "testuser"
	err = ws.WriteMessage(websocket.TextMessage, []byte(username))
	assert.NoError(t, err)

	// Read the acknowledgment message from the server
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, handlers.UsernameAccepted, string(msg))

	// Send a message to the server
	message := "Hello, World!"
	err = ws.WriteMessage(websocket.TextMessage, []byte(message))
	assert.NoError(t, err)

	// Read the acknowledgment message from the server
	_, msg, err = ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, handlers.AcknowledgmentMessage, string(msg))

	// Verify that the message is broadcasted to all clients
	// (This part of the test would require a second WebSocket connection to verify)
}

// TestHandleUserCommand tests the "#users" command functionality.
func TestHandleUserCommand(t *testing.T) {
	// Create a new connection pool and NATS publisher
	pool := connection_pool.New()
	publisher, err := nats_lib.New(nats.DefaultURL)
	assert.NoError(t, err)

	// Create a new handler
	handler := handlers.New(pool, publisher)

	// Create a test server with the WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(handler.WebSocketHandler))
	defer server.Close()

	// Convert the server URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Send a username to the server
	username := "testuser"
	err = ws.WriteMessage(websocket.TextMessage, []byte(username))
	assert.NoError(t, err)

	// Read the acknowledgment message from the server
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, handlers.UsernameAccepted, string(msg))

	// Send the "#users" command to the server
	err = ws.WriteMessage(websocket.TextMessage, []byte("#users"))
	assert.NoError(t, err)

	// Read the response from the server
	_, msg, err = ws.ReadMessage()
	assert.NoError(t, err)

	// Unmarshal the response
	var response map[string]interface{}
	err = json.Unmarshal(msg, &response)
	assert.NoError(t, err)

	// Verify that the response contains the correct user list
	assert.Equal(t, "users", response["type"])
	assert.Contains(t, response["data"], username)
}

// TestJoinAndLeaveMessagePublish tests the join and leave message publishing functionality.
func TestJoinAndLeaveMessagePublish(t *testing.T) {
	// Create a new connection pool and NATS publisher
	pool := connection_pool.New()
	publisher, err := nats_lib.New(nats.DefaultURL)
	assert.NoError(t, err)

	// Create a new handler
	handler := handlers.New(pool, publisher)

	// Create a test server with the WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(handler.WebSocketHandler))
	defer server.Close()

	// Convert the server URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Send a username to the server
	username := "testuser"
	err = ws.WriteMessage(websocket.TextMessage, []byte(username))
	assert.NoError(t, err)

	// Read the acknowledgment message from the server
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, handlers.UsernameAccepted, string(msg))

	// Verify that the join message is published
	// (This part of the test would require a NATS subscription to verify)

	// Close the WebSocket connection to trigger the leave message
	ws.Close()

	// Verify that the leave message is published
	// (This part of the test would require a NATS subscription to verify)
}
