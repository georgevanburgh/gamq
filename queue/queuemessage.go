package queue

import (
	"github.com/FireEater64/gamq/message"
)

type queueMessage struct {
	next *queueMessage
	data *message.Message
}

func newQueueMessage(givenMessage *message.Message) *queueMessage {
	return &queueMessage{data: givenMessage}
}
