package main

import (
	"log"

	"github.com/fishdontexist/chatroom/internal/app/server"
	"github.com/fishdontexist/chatroom/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	chatApp, err := server.New(cfg.Nats.URL)
	if err != nil {
		log.Fatalf("Error creating app: %v", err)
	}
	defer chatApp.Publisher.Close()

	chatApp.SetupRoutes()

	chatApp.StartServer(cfg.App.Port)

}
