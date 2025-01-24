package nats_lib

import (
	"log"

	"github.com/nats-io/nats.go"
)

type IPublisher interface {
	Publish(subject string, message string)
	NATSConnection() *nats.Conn
}

type Publisher struct {
	conn *nats.Conn
}

func New(natsURL string) (*Publisher, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: nc}, nil
}

func (p *Publisher) Publish(subject string, message string) {
	// convert message to []byte with heloper method
	messageToByte := stringToBytes(message)
	if err := p.conn.Publish(subject, messageToByte); err != nil {
		log.Printf("Error publishing the message to NATS: %v", err)
	}
}

func (p *Publisher) NATSConnection() *nats.Conn {
	return p.conn
}

func stringToBytes(message string) []byte {
	return []byte(message)
}

func (p *Publisher) Close() {
	p.conn.Close()
}
