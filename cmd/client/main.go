package main

import (
	"log"

	"github.com/fishdontexist/chatroom/internal/app/client"
	"github.com/fishdontexist/chatroom/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error reading the config")
		return
	}
	client, err := client.New(cfg.App.ServerURL)
	if err != nil {
		log.Fatalf("Error connecting to WebSocket server: %v", err)
	}

	client.Start()
}
