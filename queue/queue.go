package queue

import (
	"github.com/FireEater64/gamq/message"
)

// Queue represents an in-memory queue, using a linked-list data structure, and
// channel io
type Queue struct {
	Name          string
	InputChannel  chan *message.Message
	OutputChannel chan *message.Message
	list          []*message.Message
}

// NewQueue create a new queue, with given name
func NewQueue(givenName string) *Queue {
	q := Queue{}
	q.InputChannel = make(chan *message.Message)
	q.OutputChannel = make(chan *message.Message)
	q.list = []*message.Message{}

	go q.pump()

	return &q
}

// Run in a goroutine - receives and transmits messages on the channels
func (q *Queue) pump() {
pump:
	for {
		// If we have no messages - block until we receive one
		if len(q.list) == 0 {
			newMessage, ok := <-q.InputChannel
			if !ok {
				break pump // Someone closed our input channel
			}
			q.list = append(q.list, newMessage)
		}

		select {
		case newMessage, ok := <-q.InputChannel:
			if !ok {
				break pump // Queue is shutting down
			}
			q.list = append(q.list, newMessage)
		case q.OutputChannel <- q.list[0]:
			q.list = q.list[1:]
		}
	}

	// Finish sending remaining values
	for _, value := range q.list {
		q.OutputChannel <- value
	}

	close(q.OutputChannel)
}
