package gamq

import ()

type Queue struct {
	Name        string
	subscribers []*Client
}

func (q *Queue) Initialize() {
	q.subscribers = make([]*Client, 0)
}

func (q *Queue) AddSubscriber(givenSubscriber *Client) {
	_ = "breakpoint"
	q.subscribers = append(q.subscribers, givenSubscriber)
	_ = "breakpoint"
}

func (q *Queue) RemoveSubscriber(givenSubscriber *Client) {
}

func (q *Queue) PublishMessage(givenMessage string) {
	// Write the message to each subscriber (currently no queuing takes place)
	for _, subscriber := range q.subscribers {
		subscriber.Writer.WriteString(givenMessage)
		subscriber.Writer.Flush()
	}
}
