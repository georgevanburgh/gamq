package gamq

import (
	"github.com/FireEater64/gamq/message"
	"github.com/FireEater64/gamq/queue"
	"sync/atomic"
	"time"

	log "github.com/cihub/seelog"
)

type messageQueue struct {
	Name                   string
	messageInput           chan *message.Message
	messageOutput          chan *message.Message
	queue                  *queue.Queue
	metrics                chan<- *Metric
	closing                chan<- *string
	subscribers            map[string]*messageShipper
	running                bool
	messagesSentLastSecond uint64 // messagesSentLastSecond should never be > 0
}

func newMessageQueue(queueName string, metricsChannel chan<- *Metric, closingChannel chan<- *string) *messageQueue {
	q := messageQueue{Name: queueName}
	q.messageInput = make(chan *message.Message)
	q.subscribers = make(map[string]*messageShipper)
	q.metrics = metricsChannel
	q.closing = closingChannel

	q.queue = queue.NewQueue(queueName)

	// Launch the metrics handler and unsubscribed listener
	go q.logMetrics()
	q.running = true

	return &q
}

func (q *messageQueue) Close() {
	log.Debugf("Closing %s", q.Name)

	// Close all subscribers
	for _, subscriber := range q.subscribers {
		q.closeSubscriber(subscriber.ClientName)
	}

	q.closing <- &q.Name
	q.running = false
}

func (q *messageQueue) Publish(givenMessage *message.Message) {
	q.queue.InputChannel <- givenMessage
	atomic.AddUint64(&q.messagesSentLastSecond, 1)
}

func (q *messageQueue) AddSubscriber(givenSubscriber *Client) {
	go q.listenForDisconnectingSubscribers(givenSubscriber)
	q.subscribers[givenSubscriber.Name] = newMessageShipper(q.queue.OutputChannel, givenSubscriber, q.metrics, q.Name)
}

func (q *messageQueue) listenForDisconnectingSubscribers(givenSubscriber *Client) {
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

func (q *messageQueue) closeSubscriber(givenSubscriberName string) {
	// Close the subscribers channel, and remove the subscriber from our array
	log.Debugf("%s unsubscribed from %s", givenSubscriberName, q.Name)
	q.subscribers[givenSubscriberName].CloseChannel <- true
	delete(q.subscribers, givenSubscriberName)
}

func (q *messageQueue) logMetrics() {
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
		q.metrics <- &Metric{Name: q.Name + ".pending", Value: int64(q.queue.PendingMessages()), Type: "guage"}
	}
}
