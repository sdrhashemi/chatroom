package connection_pool

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionPool struct {
	mu   sync.Mutex
	Pool map[string]*websocket.Conn
}

func New() *ConnectionPool {
	return &ConnectionPool{
		Pool: make(map[string]*websocket.Conn),
	}
}

func (cp *ConnectionPool) AddConnection(name string, conn *websocket.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.Pool[name] = conn
}

func (cp *ConnectionPool) RemoveConnection(name string, conn *websocket.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	delete(cp.Pool, name)
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

func (cp *ConnectionPool) UserNameExists(name string) bool {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	_, exists := cp.Pool[name]
	return exists
}

func (cp *ConnectionPool) GetUsers() []*websocket.Conn {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	connections := make([]*websocket.Conn, 0, len(cp.Pool))
	for _, conn := range cp.Pool {
		connections = append(connections, conn)
	}
	return connections
}

func (cp *ConnectionPool) GetUserNames() []string {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	userNames := make([]string, 0, len(cp.Pool))
	for index, _ := range cp.Pool {
		userNames = append(userNames, index)
	}
	return userNames
}
