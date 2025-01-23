package connection_pool

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionPool struct {
	mu   sync.Mutex
	Pool []*websocket.Conn
}

func New() *ConnectionPool {
	return &ConnectionPool{
		Pool: make([]*websocket.Conn, 0),
	}
}

func (cp *ConnectionPool) AddConnection(conn *websocket.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.Pool = append(cp.Pool, conn)
}

func (cp *ConnectionPool) RemoveConnection(conn *websocket.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	for i, c := range cp.Pool {
		if c == conn {
			cp.Pool = append(cp.Pool[:i], cp.Pool[i+1:]...)
			break
		}
	}
}

func (cp *ConnectionPool) BroadcastMessageToClients(message []byte) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for _, conn := range cp.Pool {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error writing message to client: %v", err)
		}
	}
}

func (cp *ConnectionPool) GetUsers() []*websocket.Conn {
	return cp.Pool
}
