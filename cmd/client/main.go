package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
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

	// fmt.Print("Please Enter your name: ")
	// credentialScanner := bufio.NewScanner(os.Stdin)
	// var userName string

	// for {
	// 	if credentialScanner.Scan() {
	// 		userName = strings.TrimSpace(credentialScanner.Text())
	// 		if userName != "" {
	// 			break
	// 		}
	// 		fmt.Print("Name cannot be empty, enter your name again: ")
	// 	} else {
	// 		log.Println("Error reading input, existing...")
	// 		return
	// 	}
	// }
	// if err := conn.WriteMessage(websocket.TextMessage, []byte(userName)); err != nil {
	// 	log.Println("Error sending message: ", err)
	// 	return
	// }

	if !handleUsernameSetup(conn) {
		log.Println("Failed to set up username. Exiting...")
		return
	}

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
			if text == "" {
				fmt.Print("> ")
				continue
			}
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

func handleUsernameSetup(conn *websocket.Conn) bool {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		_, serverMessage, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading server message: %v", err)
			return false
		}

		fmt.Println(string(serverMessage))

		if scanner.Scan() {
			username := strings.TrimSpace(scanner.Text())
			if username == "" {
				fmt.Println("Username cannot be empty, try again: ")
				continue
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte(username))
			if err != nil {
				log.Printf("Error sending username to server: %v", err)
				return false
			}

			// wait for server response
			_, response, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading server response: %v", err)
				return false
			}

			if string(response) == "OK" {
				log.Println("Username accepted.")
				return true
			}

			// if server rejected username
			fmt.Println(string(response))

		} else {
			log.Println("Error reading input. Exiting...")
			return false
		}
	}
}
