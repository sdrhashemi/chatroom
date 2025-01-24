package handlers_test

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebscoketHandler(t *testing.T) {
	wsURL := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Error dialing websocket: %v", err)
	}
	defer conn.Close()

}
