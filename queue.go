package gamq

import (
	log "github.com/cihub/seelog"
	"sync/atomic"
	"time"
)

type Queue struct {
	Name                   string
	messages               chan *string
	subscribers            []*Client
	running                bool
	messagesSentLastSecond uint64 // messagesSentLastSecond should never be > 0
}

func (q *Queue) Initialize() {
	q.messages = make(chan *string)
	q.subscribers = make([]*Client, 0)

	messageHandler1 := DummyMessageHandler{}
	messageHandler2 := DummyMessageHandler{}
	messageShipper := MessageShipper{}

	// Hook the flow together
	messageShipper.Initialize(
		messageHandler2.Initialize(
			messageHandler1.Initialize(q.messages)), &q.subscribers)

	// Launch the metrics handler
	go q.logMetrics()
	q.running = true
}

func (q *Queue) Close() {
	close(q.messages)
	q.running = false
}

func (q *Queue) Publish(givenMessage *string) {
	q.messages <- givenMessage
	atomic.AddUint64(&q.messagesSentLastSecond, 1)
}

func (q *Queue) AddSubscriber(givenSubscriber *Client) {
	q.subscribers = append(q.subscribers, givenSubscriber)
}

func (q *Queue) logMetrics() {
	for _ = range time.Tick(time.Second) {
		// Die with the handler
		if !q.running {
			break
		}

		// Print out the number of messages published last second
		currentValue := atomic.LoadUint64(&q.messagesSentLastSecond)
		log.Infof("%d/second", currentValue)
		atomic.StoreUint64(&q.messagesSentLastSecond, 0)
	}
}
