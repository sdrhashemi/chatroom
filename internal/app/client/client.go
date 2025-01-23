package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fishdontexist/chatroom/pkg/client/ui"
	"github.com/fishdontexist/chatroom/pkg/message"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	done chan struct{}
}

func New(serverURL string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		done: make(chan struct{}),
	}, nil
}

func (c *Client) Start() {
	defer c.conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	if !c.handleUsernameSetup() {
		fmt.Println("Failed to set up username. Exiting...")
		return
	}

	go c.readMessages()
	go c.writeMessages()

	select {
	case <-c.done:
	case <-interrupt:
		log.Println("Interrupt signal received, shutting down...")
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}
}

func (c *Client) readMessages() {
	defer close(c.done)
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		msg, err := message.Deserialize(data)
		if err != nil {
			ui.DisplayMessage(string(data))
			continue
		}
		switch msg.Type {
		case "users":
			if users, ok := msg.Data.([]interface{}); ok {
				ui.DisplayUsers(users)
			} else {
				ui.DisplayMessage("Invalid users list recieved. cannot show")
			}
		default:
			ui.DisplayMessage(fmt.Sprintf("Unknown message type: %s", msg.Type))
		}
	}
}

func (c *Client) writeMessages() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if text == "/exit" {
			fmt.Println("Exiting...")
			c.done <- struct{}{}
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}
	}
}

func (c *Client) handleUsernameSetup() bool {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading server message: %v", err)
			return false
		}

		ui.DisplayMessage(string(data))

		username := ui.PromptUserName()
		if username == "" {
			continue
		}

		err = c.conn.WriteMessage(websocket.TextMessage, []byte(username))
		if err != nil {
			log.Printf("Error sending username: %v", err)
			return false
		}

		_, response, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading server response: %v", err)
			return false
		}

		if string(response) == "OK" {
			log.Println("Username accepted.")
			return true
		}

		ui.DisplayMessage(string(response))
	}
}
