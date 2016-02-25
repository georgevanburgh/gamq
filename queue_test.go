package gamq

import (
	"bufio"
	"bytes"
	"github.com/FireEater64/gamq/message"
	"strings"
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

	testMessagePayload := []byte("Testing!")
	testMessage := message.NewHeaderlessMessage(&testMessagePayload)

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
	underTest.Publish(testMessage)

	gomega.Eventually(func() []byte {
		return writerBuffer.Bytes()
	}).Should(gomega.Equal(testMessagePayload))
}

func TestQueue_sendMessage_generatesMetrics(t *testing.T) {
	// More async testing
	gomega.RegisterTestingT(t)

	// We should receive metrics ending in these names from a queue during
	// normal operation
	expectedMetricNames := [...]string{"messagerate", "subscribers", "pending"}

	// Mocking
	dummyMetricsChannel := make(chan *Metric)
	dummyClosingChannel := make(chan *string)

	underTest := Queue{}
	underTest.Initialize(dummyMetricsChannel, dummyClosingChannel)

	// After a subscriber is added, we should start receiving metrics
	dummySubscriber := Client{Closed: new(chan bool)}
	underTest.AddSubscriber(&dummySubscriber)

	seenMetricNames := make(map[string]bool)
	go func() {
		for {
			metric := <-dummyMetricsChannel
			metricNameChunks := strings.Split(metric.Name, ".")
			finalMetricName := metricNameChunks[len(metricNameChunks)-1]
			seenMetricNames[finalMetricName] = true
		}
	}()

	// Check we've received metrics ending in all the expected names
	// NOTE: It might take longer than the default gomega 1 second timeout to
	// receive all the metrics we're expecting
	gomega.Eventually(func() bool {
		toReturn := true
		for _, metricName := range expectedMetricNames {
			if !seenMetricNames[metricName] {
				toReturn = false
			}
		}
		return toReturn
	}, "5s").Should(gomega.BeTrue()) //  Timeout upped to 5 seconds
}

// A unsubscribing client should not be considered for message delivery
func TestQueue_sendMessageAfterUnsubscribe_messageReceivedSuccessfully(t *testing.T) {
	// Need gomega for async testing
	gomega.RegisterTestingT(t)

	underTest := Queue{Name: TEST_QUEUE_NAME}
	testMessagePayload := []byte("Testing!")
	testMessage := message.NewHeaderlessMessage(&testMessagePayload)

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
	underTest.Publish(testMessage)

	// Bit of a hack - only one of the subscribers will get the message,
	// and we don't know which one
	gomega.Eventually(func() []byte {
		if writerBuffer1.String() == "" {
			return writerBuffer2.Bytes()
		} else {
			return writerBuffer1.Bytes()
		}
	}).Should(gomega.Equal(testMessagePayload))

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
	underTest.Publish(testMessage)

	gomega.Eventually(func() []byte {
		return writerBuffer2.Bytes()
	}).Should(gomega.Equal(testMessagePayload))
}

func TestQueue_xPendingMetrics_producesCorrectMetric(t *testing.T) {
	// Need gomega for async testing
	gomega.RegisterTestingT(t)

	numberOfMessagesToSend := 10

	underTest := Queue{Name: TEST_QUEUE_NAME}
	testMessagePayload := []byte("Testing!")
	testMessage := message.NewHeaderlessMessage(&testMessagePayload)

	dummyMetricsPipe := make(chan *Metric)
	dummyClosingPipe := make(chan *string)
	underTest.Initialize(dummyMetricsPipe, dummyClosingPipe)

	for i := 0; i < numberOfMessagesToSend; i++ {
		underTest.Publish(testMessage)
	}

	// Eventually, we should see `numberOfMessagesToSend` pending messages
	gomega.Eventually(func() int {
		metric := <-dummyMetricsPipe
		if strings.Contains(metric.Name, "pending") {
			return int(metric.Value)
		} else {
			return -1
		}
	}, "5s").Should(gomega.Equal(numberOfMessagesToSend))
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
