package api

import (
	"github.com/fishdontexist/chatroom/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ws", handlers.WebSocketHandler).Methods("GET")
	return router
}
