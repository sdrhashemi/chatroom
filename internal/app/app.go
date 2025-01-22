package app

import (
	"log"
	"net/http"

	"github.com/fishdontexist/chatroom/internal/handlers"
	"github.com/fishdontexist/chatroom/pkg/connection_pool"
	"github.com/fishdontexist/chatroom/pkg/nats"
)

type App struct {
	Handler    *handlers.Handler
	Connection *connection_pool.ConnectionPool
	Publisher  *nats.Publisher
}

func New(natsURL string) (*App, error) {
	publisher, err := nats.New(natsURL)
	if err != nil {
		return nil, err
	}

	pool := connection_pool.New()

	handler := handlers.New(pool, publisher)

	return &App{
		Handler:    handler,
		Connection: pool,
		Publisher:  publisher,
	}, nil
}

func (a *App) SetupRoutes() {
	http.HandleFunc("/ws", a.Handler.WebSocketHandler)
}

func (a *App) StartServer(addr string) {
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
