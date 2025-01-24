package client

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fishdontexist/chatroom/pkg/client/ui"
	"github.com/fishdontexist/chatroom/pkg/message"
	"github.com/gorilla/websocket"
)

const UsernameAccepted = "Username Accepted"

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
		ui.DisplayError("Failed to set up username. Exiting...", true)
		return
	}

	go c.readMessages()
	go c.writeMessages()

	select {
	case <-c.done:
	case <-interrupt:
		ui.DisplayError("Interrupt signal received, shutting down...", true)
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}
}

func (c *Client) readMessages() {
	defer close(c.done)
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			ui.DisplayError(fmt.Sprintf("Error reading message: %v", err), true)
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
		case "chatroom":
			ui.DisplayNeutral(msg.Data.(string))
		case "ack":
			ui.DisplaySuccess("[Sent]")
		default:
			ui.DisplayError(fmt.Sprintf("Unknown message type: %s", msg.Type), false)
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
			ui.DisplayError(fmt.Sprintf("Error sending message: %v", err), false)
			continue
		}
	}
}

func (c *Client) handleUsernameSetup() bool {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			ui.DisplayError(fmt.Sprintf("Error reading server message: %v", err), false)
			return false
		}

		ui.DisplayNeutral(string(data))

		username := ui.PromptUserName()
		if username == "" {
			continue
		}

		err = c.conn.WriteMessage(websocket.TextMessage, []byte(username))
		if err != nil {
			ui.DisplayError(fmt.Sprintf("Error sending username: %v", err), false)
			return false
		}

		_, response, err := c.conn.ReadMessage()
		if err != nil {
			ui.DisplayError(fmt.Sprintf("Error reading server response: %v", err), false)
			return false
		}

		if string(response) == UsernameAccepted {
			ui.DisplaySuccess("Username Accepted")
			return true
		}

		ui.DisplayMessage(string(response))
	}
}
