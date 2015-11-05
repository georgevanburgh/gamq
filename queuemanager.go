package gamq

import (
	"fmt"
)

type QueueManager struct {
	queues map[string]*Queue
}

func (qm *QueueManager) Initialize() {
	qm.queues = make(map[string]*Queue)
	fmt.Println("Initialized")
}

func (qm *QueueManager) Publish(queueName string, message string) {
	fmt.Printf("Publishing message to %s: %s\n", queueName, message)

	_ = "breakpoint"

	queueToPublishTo := qm.getQueueSafely(queueName)
	queueToPublishTo.PublishMessage(message)
}

func (qm *QueueManager) Subscribe(queueName string, client *Client) {
	fmt.Printf("Subscribing to %s\n", queueName)

	// var queueToSubscribeTo Queue
	// var present bool

	queueToSubscribeTo := qm.getQueueSafely(queueName)
	queueToSubscribeTo.AddSubscriber(client)
}

func (qm *QueueManager) getQueueSafely(queueName string) *Queue {
	_ = "breakpoint"
	queueToReturn, present := qm.queues[queueName]
	if !present {
		newQueue := Queue{Name: queueName}
		newQueue.Initialize()
		qm.queues[queueName] = &newQueue
		queueToReturn = qm.queues[queueName]
	}

	return queueToReturn
}
