package message

import (
	"time"
)

type Message struct {
	Body       *[]byte
	Head       *map[string]string
	receivedAt time.Time
}

func NewMessage(givenHead *map[string]string, givenBody *[]byte) *Message {
	return &Message{Head: givenHead, Body: givenBody, receivedAt: time.Now()}
}

func NewHeaderlessMessage(givenBody *[]byte) *Message {
	emptyHeader := make(map[string]string)
	return &Message{Head: &emptyHeader, Body: givenBody, receivedAt: time.Now()}
}
