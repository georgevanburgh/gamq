package gamq

import (
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
	// log.Debugf("Publishing message to %s: %s", queueName, message)

	queueToPublishTo := qm.getQueueSafely(queueName)
	queueToPublishTo.Messages <- &message
}

func (qm *QueueManager) Subscribe(queueName string, client *Client) {
	log.Infof("%s subscribed to %s", client.Name, queueName)

	queueToSubscribeTo := qm.getQueueSafely(queueName)
	queueToSubscribeTo.Subscribers <- client
}

func (qm *QueueManager) CloseQueue(queueName string) {
	log.Infof("Closing %s", queueName)
	queueToClose := qm.getQueueSafely(queueName)
	delete(qm.queues, queueName)
	queueToClose.Close()
}

func (qm *QueueManager) getQueueSafely(queueName string) *Queue {
	queueToReturn, present := qm.queues[queueName]
	if !present {
		newQueue := Queue{Name: queueName}
		newQueue.Initialize()
		qm.queues[queueName] = &newQueue
		queueToReturn = qm.queues[queueName]
	}

	return queueToReturn
}
