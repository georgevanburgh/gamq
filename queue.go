package gamq

import ()

type Queue struct {
	Name        string
	Messages    chan<- *string
	Subscribers chan<- *Client
}

func (q *Queue) Initialize() {
	messages := make(chan *string)
	subscribers := make(chan *Client)

	messageHandler1 := DummyMessageHandler{}
	messageHandler2 := DummyMessageHandler{}
	messageShipper := MessageShipper{}

	// Hook the flow together
	messageShipper.Initialize(
		messageHandler2.Initialize(
			messageHandler1.Initialize(messages)), subscribers)

	q.Messages = messages
	q.Subscribers = subscribers
}

func (q *Queue) Close() {
	close(q.Messages)
	close(q.Subscribers)
}
