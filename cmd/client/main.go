package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {

	serverURL := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatalf("Error connecting to Websocket server: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to Websocket server")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message from Websocket server: %v", err)
				return
			}
			fmt.Printf("\r%s\n> ", message)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("> ")
		for scanner.Scan() {
			text := scanner.Text()
			if text == "/exit" {
				log.Println("Exiting...")
				interrupt <- syscall.SIGINT
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("Error sending message: ", err)
				return
			}
			fmt.Print("> ")

		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from stdin: %v", err)
		}
	}()

	select {
	case <-done:
	case <-interrupt:
		log.Println("Interrupt signal received, shutting down...")
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Error sending close message: ", err)

		}
	}
}
