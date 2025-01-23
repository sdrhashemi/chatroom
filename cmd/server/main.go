package main

import (
	"log"

	"github.com/fishdontexist/chatroom/internal/app"
	"github.com/fishdontexist/chatroom/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("../../config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	chatApp, err := app.New(cfg.Nats.URL)
	if err != nil {
		log.Fatalf("Error creating app: %v", err)
	}
	defer chatApp.Publisher.Close()

	chatApp.SetupRoutes()

	chatApp.StartServer(cfg.App.Port)

}
