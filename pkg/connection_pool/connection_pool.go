package connection_pool

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionPool struct {
	mu   sync.Mutex
	pool []*websocket.Conn
}

func New() *ConnectionPool {
	return &ConnectionPool{
		pool: make([]*websocket.Conn, 0),
	}
}

func (cp *ConnectionPool) AddConnection(conn *websocket.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.pool = append(cp.pool, conn)
}

func (cp *ConnectionPool) RemoveConnection(conn *websocket.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	for i, c := range cp.pool {
		if c == conn {
			cp.pool = append(cp.pool[:i], cp.pool[i+1:]...)
			break
		}
	}
}

func (cp *ConnectionPool) GetUsers() []*websocket.Conn {
	return cp.pool
}
