package gamq

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/onsi/gomega"
)

const (
	TEST_QUEUE_NAME = "TestQueue"
)

// Check that messages sent to a queue are eventually sent to consumers
func TestQueue_sendMessage_messageReceivedSuccessfully(t *testing.T) {
	// Need gomega for async testing
	gomega.RegisterTestingT(t)

	underTest := Queue{Name: TEST_QUEUE_NAME}
	testMessage := "Testing!"

	dummyMetricsPipe := make(chan<- *Metric)
	underTest.Initialize(dummyMetricsPipe)

	writerBuffer := new(bytes.Buffer)
	dummyWriter := bufio.NewWriter(writerBuffer)
	closedChannel := make(chan bool)
	dummyClient := Client{Name: "Test", Writer: dummyWriter, Closed: &closedChannel}

	// Add the subscription
	underTest.AddSubscriber(&dummyClient)

	// Queue the message
	underTest.Publish(&testMessage)

	gomega.Eventually(func() string {
		return writerBuffer.String()
	}).Should(gomega.Equal(testMessage))
}

func TestQueue_initialize_completesSuccessfully(t *testing.T) {
	underTest := Queue{Name: TEST_QUEUE_NAME}

	dummyMetricsPipe := make(chan<- *Metric)
	underTest.Initialize(dummyMetricsPipe)

	// Queue should be named correctly
	if underTest.Name != TEST_QUEUE_NAME {
		t.Fail()
	}
}
