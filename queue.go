package gamq

import (
	log "github.com/cihub/seelog"
	"sync/atomic"
	"time"
)

type Queue struct {
	Name                   string
	messages               chan *string
	subscribers            map[string]*Client
	running                bool
	messagesSentLastSecond uint64 // messagesSentLastSecond should never be > 0
}

func (q *Queue) Initialize() {
	q.messages = make(chan *string)
	q.subscribers = make(map[string]*Client)

	messageHandler1 := DummyMessageHandler{}
	messageHandler2 := DummyMessageHandler{}
	messageShipper := MessageShipper{}

	// Hook the flow together
	messageShipper.Initialize(
		messageHandler2.Initialize(
			messageHandler1.Initialize(q.messages)), &q.subscribers)

	// Launch the metrics handler and unsubscribed listener
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
	go q.listenForDisconnectingSubscribers(givenSubscriber)
	q.subscribers[givenSubscriber.Name] = givenSubscriber
}

func (q *Queue) listenForDisconnectingSubscribers(givenSubscriber *Client) {
	disconnectMessage := <-*givenSubscriber.Closed
	if disconnectMessage {
		// Remove the subscriber
		log.Debugf("%s unsubscribed from %s", givenSubscriber.Name, q.Name)
		delete(q.subscribers, givenSubscriber.Name)

		// Close the queue if we have no more subscribers
		// TODO: Should also check for pending messages/publishers
		if len(q.subscribers) == 0 {
			log.Debugf("No subscribers left on queue %s - closing", q.Name)
			q.Close()
		}
	}

	// Other subscribers care about knowing the channel is closed
	*givenSubscriber.Closed <- disconnectMessage
}

func (q *Queue) logMetrics() {
	for _ = range time.Tick(time.Second) {
		// Die with the handler
		if !q.running {
			break
		}

		// Print out various metrics
		currentValue := atomic.LoadUint64(&q.messagesSentLastSecond)
		log.Infof("%d/second", currentValue)
		log.Infof("%d subscribers", len(q.subscribers))
		atomic.StoreUint64(&q.messagesSentLastSecond, 0)
	}
}
