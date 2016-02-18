package gamq

import (
	"sync/atomic"
	"time"

	log "github.com/cihub/seelog"
)

type Queue struct {
	Name                   string
	messageInput           chan *string
	messageOutput          chan *string
	metrics                chan<- *Metric
	closing                chan<- *string
	subscribers            map[string]*MessageShipper
	running                bool
	messagesSentLastSecond uint64 // messagesSentLastSecond should never be > 0
}

func (q *Queue) Initialize(metricsChannel chan<- *Metric, closingChannel chan<- *string) {
	q.messageInput = make(chan *string)
	q.subscribers = make(map[string]*MessageShipper)
	q.metrics = metricsChannel
	q.closing = closingChannel

	messageHandler1 := DummyMessageHandler{}
	messageHandler2 := DummyMessageHandler{}

	// Hook the flow together
	q.messageOutput = messageHandler2.Initialize(
		messageHandler1.Initialize(q.messageInput))

	// Launch the metrics handler and unsubscribed listener
	go q.logMetrics()
	q.running = true
}

func (q *Queue) Close() {
	log.Debugf("Closing %s", q.Name)

	// Close all subscribers
	for _, subscriber := range q.subscribers {
		q.closeSubscriber(subscriber.ClientName)
	}

	q.closing <- &q.Name
	q.running = false
}

func (q *Queue) Publish(givenMessage *string) {
	q.messageInput <- givenMessage
	atomic.AddUint64(&q.messagesSentLastSecond, 1)
}

func (q *Queue) AddSubscriber(givenSubscriber *Client) {
	go q.listenForDisconnectingSubscribers(givenSubscriber)
	messageShipper := MessageShipper{}
	messageShipper.Initialize(q.messageOutput, givenSubscriber)
	q.subscribers[givenSubscriber.Name] = &messageShipper
}

func (q *Queue) listenForDisconnectingSubscribers(givenSubscriber *Client) {
	disconnectMessage := <-*givenSubscriber.Closed
	if disconnectMessage {
		q.closeSubscriber(givenSubscriber.Name)

		// Close the queue if we have no more subscribers
		// TODO: Should also check for pending messages/publishers
		if len(q.subscribers) == 0 && len(q.messageOutput) == 0 {
			log.Debugf("No subscribers left on queue %s - closing", q.Name)
			q.Close()
		}
	}

	// Other subscribers care about knowing the channel is closed
	*givenSubscriber.Closed <- disconnectMessage
}

func (q *Queue) closeSubscriber(givenSubscriberName string) {
	// Close the subscribers channel, and remove the subscriber from our array
	log.Debugf("%s unsubscribed from %s", givenSubscriberName, q.Name)
	q.subscribers[givenSubscriberName].CloseChannel <- true
	delete(q.subscribers, givenSubscriberName)
}

func (q *Queue) logMetrics() {
	for _ = range time.Tick(time.Second) {
		// Die with the handler
		if !q.running {
			break
		}

		// If this is the metrics queue - don't log metrics
		if q.Name == metricsQueueName {
			break
		}

		// Send various metrics
		currentMessageRate := atomic.SwapUint64(&q.messagesSentLastSecond, 0)

		q.metrics <- &Metric{Name: q.Name + ".messagerate", Value: int64(currentMessageRate), Type: "counter"}
		q.metrics <- &Metric{Name: q.Name + ".subscribers", Value: int64(len(q.subscribers)), Type: "guage"}
		q.metrics <- &Metric{Name: q.Name + ".pending", Value: int64(len(q.messageOutput)), Type: "guage"}
	}
}
