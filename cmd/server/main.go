package main

import (
	"fmt"

	"github.com/fishdontexist/chatroom/internal/api"
)

func main() {
	natsURL := "nats://localhost:4222"
	chatApp, err := app.NewApp(natsURL)
	if err!= nil {
		log.Fatalf("Error creating app: %v", err)
	}
	defer chatApp.Connection.Close()

	chatApp.SetupRoutes()

	chatApp.StartServer(":8080")

}
