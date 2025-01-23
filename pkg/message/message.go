package message

import "encoding/json"

type Message struct {
	Type string
	Data interface{}
}

func Serialize(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

func Deserialize(data []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return msg, err
}
