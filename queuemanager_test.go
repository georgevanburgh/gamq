package gamq

import (
	"testing"

	"github.com/onsi/gomega"
)

const (
	TestQueueName = "testing"
)

// Reusing the same (closed) channel name shouldn't make us crash
// Wait between sub/unsub operations to make sure metrics are sent proparly
func TestQueueManager_queuesClosed_removedFromMap(t *testing.T) {
	config := Config{}
	SetConfig(&config)

	underTest := NewQueueManager()

	// Create dummy clients
	dummyClient := Client{}
	dummyClient.Name = "Dummy"
	closedChannel := make(chan bool)
	dummyClient.Closed = &closedChannel

	// Subscribe
	underTest.Subscribe(TestQueueName, &dummyClient)

	// Close the queue
	*dummyClient.Closed <- true

	// Check that TestQueueName is removed from the QueueManager map
	gomega.Eventually(func() bool {
		return underTest.queues[TestQueueName] != nil
	}).Should(gomega.BeTrue())
}
