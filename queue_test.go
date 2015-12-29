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
	dummyClosingPipe := make(chan<- *string)
	underTest.Initialize(dummyMetricsPipe, dummyClosingPipe)

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

// A unsubscribing client should not be considered for message delivery
func TestQueue_sendMessageAfterUnsubscribe_messageReceivedSuccessfully(t *testing.T) {
	// Need gomega for async testing
	gomega.RegisterTestingT(t)

	underTest := Queue{Name: TEST_QUEUE_NAME}
	testMessage := "Testing!"

	dummyMetricsPipe := make(chan<- *Metric)
	dummyClosingPipe := make(chan<- *string)
	underTest.Initialize(dummyMetricsPipe, dummyClosingPipe)

	writerBuffer1 := new(bytes.Buffer)
	dummyWriter1 := bufio.NewWriter(writerBuffer1)
	closedChannel1 := make(chan bool)
	dummyClient1 := Client{Name: "Test1", Writer: dummyWriter1, Closed: &closedChannel1}

	writerBuffer2 := new(bytes.Buffer)
	dummyWriter2 := bufio.NewWriter(writerBuffer2)
	closedChannel2 := make(chan bool)
	dummyClient2 := Client{Name: "Test2", Writer: dummyWriter2, Closed: &closedChannel2}

	// Add the subscription
	underTest.AddSubscriber(&dummyClient1)
	underTest.AddSubscriber(&dummyClient2)

	// Queue the message
	underTest.Publish(&testMessage)

	// Bit of a hack - only one of the subscribers will get the message,
	// and we don't know which one
	gomega.Eventually(func() string {
		if writerBuffer1.String() == "" {
			return writerBuffer2.String()
		} else {
			return writerBuffer1.String()
		}
	}).Should(gomega.Equal(testMessage))

	// We'll be reusing these buffers
	writerBuffer1.Reset()
	writerBuffer2.Reset()

	// Close one client
	*dummyClient1.Closed <- true

	// Should remove the client from the map
	gomega.Eventually(func() bool {
		return underTest.subscribers[dummyClient1.Name] == nil
	}).Should(gomega.BeTrue())

	// Now send a message - the remaining client should receive it without issue
	underTest.Publish(&testMessage)

	gomega.Eventually(func() string {
		return writerBuffer2.String()
	}).Should(gomega.Equal(testMessage))
}

func TestQueue_initialize_completesSuccessfully(t *testing.T) {
	underTest := Queue{Name: TEST_QUEUE_NAME}

	dummyMetricsPipe := make(chan<- *Metric)
	dummyClosingPipe := make(chan<- *string)
	underTest.Initialize(dummyMetricsPipe, dummyClosingPipe)

	// Queue should be named correctly
	if underTest.Name != TEST_QUEUE_NAME {
		t.Fail()
	}
}
