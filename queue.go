package gamq

import ()

type Queue struct {
	Name        string
	Messages    chan<- string
	Subscribers chan<- *Client
}

func (q *Queue) Initialize() {
	messages := make(chan string)
	subscribers := make(chan *Client)

	messageHandler1 := DummyMessageHandler{}
	messageHandler2 := DummyMessageHandler{}
	messageHandler3 := DummyMessageHandler{}

	// Hook the flow together
	messageHandler3.Initialize(messageHandler2.Initialize(messageHandler1.Initialize(messages)))

	q.Messages = messages
	q.Subscribers = subscribers
}

func (q *Queue) Close() {
	close(q.Messages)
	close(q.Subscribers)
}
