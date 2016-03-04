package queue

import (
	"github.com/FireEater64/gamq/message"
)

// Queue represents an in-memory queue, using a slice-basedgit stat data structure, and
// channel io
type Queue struct {
	Name          string
	InputChannel  chan *message.Message
	OutputChannel chan *message.Message
	head          *queueMessage
	tail          *queueMessage
	length        int
}

// NewQueue create a new queue, with given name
func NewQueue(givenName string) *Queue {
	q := Queue{}
	q.InputChannel = make(chan *message.Message)
	q.OutputChannel = make(chan *message.Message)

	go q.pump()

	return &q
}

// Run in a goroutine - receives and transmits messages on the channels
func (q *Queue) pump() {
pump:
	for {
		// If we have no messages - block until we receive one
		if q.head == nil {
			newMessage, ok := <-q.InputChannel
			newQueueMessage := newQueueMessage(newMessage)
			if !ok {
				break pump // Someone closed our input channel
			}
			q.head = newQueueMessage
			q.tail = newQueueMessage
			q.length++
		}

		select {

		case newMessage, ok := <-q.InputChannel:
			if !ok {
				break pump // Someone closed our input channel - we're shutting down the pipeline
			}
			newQueueMessage := newQueueMessage(newMessage)
			q.tail.next = newQueueMessage
			q.tail = newQueueMessage
			q.length++

		case q.OutputChannel <- q.head.data:
			q.head = q.head.next
			q.length--
		}
	}

	// Finish sending remaining values
	nextMessageToSend := q.head
	for nextMessageToSend != nil {
		q.OutputChannel <- nextMessageToSend.data
		nextMessageToSend = nextMessageToSend.next
	}

	close(q.OutputChannel)
}

func (q *Queue) PendingMessages() int {
	return q.length
}
