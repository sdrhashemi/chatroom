package server

import "github.com/gorilla/mux"

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ws", handlers.HandleConnection).Methods("GET")
	return router
}
