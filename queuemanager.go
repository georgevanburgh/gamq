package gamq

import (
	"github.com/FireEater64/gamq/message"
	log "github.com/cihub/seelog"
)

type queueManager struct {
	queues                   map[string]*messageQueue
	metricsManager           *MetricsManager
	closeNotificationChannel chan *string
}

func newQueueManager() *queueManager {
	qm := queueManager{}

	qm.queues = make(map[string]*messageQueue)
	qm.closeNotificationChannel = make(chan *string, 10)

	qm.metricsManager = NewMetricsManager(&qm)

	go qm.listenForClosingQueues()

	log.Debug("Initialized QueueManager")

	return &qm
}

func (qm *queueManager) Publish(queueName string, message *message.Message) {
	// log.Debugf("Publishing message to %s: %s", queueName, message)

	queueToPublishTo := qm.getQueueSafely(queueName)
	queueToPublishTo.Publish(message)
}

func (qm *queueManager) Subscribe(queueName string, client *Client) {
	log.Infof("%s subscribed to %s", client.Name, queueName)

	queueToSubscribeTo := qm.getQueueSafely(queueName)
	queueToSubscribeTo.AddSubscriber(client)
}

func (qm *queueManager) CloseQueue(queueName string) {
	log.Infof("Closing %s", queueName)
	queueToClose := qm.getQueueSafely(queueName)
	queueToClose.Close()
}

func (qm *queueManager) getQueueSafely(queueName string) *messageQueue {
	queueToReturn, present := qm.queues[queueName]
	if !present {
		newQueue := newMessageQueue(queueName, qm.metricsManager.metricsChannel, qm.closeNotificationChannel)
		qm.queues[queueName] = newQueue
		queueToReturn = qm.queues[queueName]
	}

	return queueToReturn
}

func (qm *queueManager) listenForClosingQueues() {
	for {
		closingQueue := <-qm.closeNotificationChannel
		log.Debugf("Removing %s from active queues", *closingQueue)
		delete(qm.queues, *closingQueue)
	}
}
