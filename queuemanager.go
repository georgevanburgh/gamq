package gamq

import (
	"fmt"
	log "github.com/cihub/seelog"
)

type QueueManager struct {
	queues map[string]*Queue
}

func (qm *QueueManager) Initialize() {
	qm.queues = make(map[string]*Queue)
	log.Debug("Initialized QueueManager")
}

func (qm *QueueManager) Publish(queueName string, message string) {
	fmt.Printf("Publishing message to %s: %s\n", queueName, message)

	_ = "breakpoint"

	queueToPublishTo := qm.getQueueSafely(queueName)
	queueToPublishTo.Messages <- message
}

func (qm *QueueManager) Subscribe(queueName string, client *Client) {
	log.Infof("%s subscribed to %s", client.Name, queueName)

	queueToSubscribeTo := qm.getQueueSafely(queueName)
	queueToSubscribeTo.Subscribers <- client
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
