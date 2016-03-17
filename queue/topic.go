package queue

import (
	"github.com/FireEater64/gamq/message"
)

// Queue represents an in-memory queue, using a slice-basedgit stat data structure, and
// channel io
type Topic struct {
	Name           string
	InputChannel   chan *message.Message
	outputChannels map[string]chan *message.Message
	pointers       map[string]*queueMessage
	wakeLocks      map[string]chan struct{}
	head           *queueMessage
	tail           *queueMessage
	length         int
}

// NewQueue create a new queue, with given name
func NewTopic(givenName string) *Topic {
	t := Topic{}
	t.InputChannel = make(chan *message.Message)
	t.outputChannels = make(map[string]chan *message.Message)
	t.pointers = make(map[string]*queueMessage)
	t.wakeLocks = make(map[string]chan struct{})

	go t.receive()

	return &t
}

func (t *Topic) GetChannel(clientName string) chan *message.Message {
	channel, exists := t.outputChannels[clientName]
	if !exists {
		newChannel := make(chan *message.Message)
		t.outputChannels[clientName] = newChannel
		t.pointers[clientName] = t.tail
		t.wakeLocks[clientName] = make(chan struct{})
		go t.serviceChannel(clientName)
		channel = newChannel
	}

	return channel
}

func (t *Topic) serviceChannel(givenChannelName string) {
	for {
		if t.pointers[givenChannelName] == nil {
			<-t.wakeLocks[givenChannelName]
		}

		t.outputChannels[givenChannelName] <- t.pointers[givenChannelName].data
		t.pointers[givenChannelName] = t.pointers[givenChannelName].next
	}
}

// Run in a goroutine - receives and transmits messages on the channels
func (t *Topic) receive() {
receive:
	for {
		newMessage, ok := <-t.InputChannel
		if !ok {
			break receive // Someone closed our input channel - we're shutting down the pipeline
		}
		newQueueMessage := newQueueMessage(newMessage)

		if t.head == nil {
			t.tail = newQueueMessage
			t.length++
		} else {
			t.tail.next = newQueueMessage
			t.tail = newQueueMessage
			t.length++
		}
		// Set each consumer pointer to the head of the queue
		for name, _ := range t.pointers {
			if t.pointers[name] == nil {
				t.pointers[name] = newQueueMessage
				t.wakeLocks[name] <- struct{}{}
			}
		}
	}
}

func (t *Topic) PendingMessages() int {
	return t.length
}
