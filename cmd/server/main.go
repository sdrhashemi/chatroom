package main

import (
	"log"

	"github.com/fishdontexist/chatroom/internal/app"
)

func main() {
	natsURL := "nats://localhost:4222"
	chatApp, err := app.New(natsURL)
	if err != nil {
		log.Fatalf("Error creating app: %v", err)
	}
	defer chatApp.Publisher.Close()

	chatApp.SetupRoutes()

	chatApp.StartServer(":8080")

}
