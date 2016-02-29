package message

import (
	"time"
)

type Message struct {
	Body       *[]byte
	Head       *map[string]string
	ReceivedAt time.Time
}

func NewMessage(givenHead *map[string]string, givenBody *[]byte) *Message {
	return &Message{Head: givenHead, Body: givenBody, ReceivedAt: time.Now()}
}

func NewHeaderlessMessage(givenBody *[]byte) *Message {
	emptyHeader := make(map[string]string)
	return &Message{Head: &emptyHeader, Body: givenBody, ReceivedAt: time.Now()}
}
