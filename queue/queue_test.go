package queue

import (
	"bytes"
	"github.com/FireEater64/gamq/message"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func TestQueue_CanSendAndReceiveBasicMessages(t *testing.T) {
	underTest := NewQueue("TestQueue")

	testMessagePayload := []byte("TestMessage")
	underTest.InputChannel <- (message.NewHeaderlessMessage(&testMessagePayload))

	receivedMessage := <-underTest.OutputChannel

	if !bytes.Equal(*receivedMessage.Body, testMessagePayload) {
		t.Fail()
	}
}

func TestQueue_ReceiveBeforeSend_ReturnsExpectedResult(t *testing.T) {
	gomega.RegisterTestingT(t)

	underTest := NewQueue("TestQueue")

	var receivedMessage *message.Message
	go func() {
		receivedMessage = <-underTest.OutputChannel
	}()

	time.Sleep(time.Millisecond * 10)

	testMessagePayload := []byte("TestMessage")
	underTest.InputChannel <- (message.NewHeaderlessMessage(&testMessagePayload))

	gomega.Eventually(func() *message.Message {
		return receivedMessage
	}).Should(gomega.Not(gomega.BeNil()))

	if !bytes.Equal(*receivedMessage.Body, testMessagePayload) {
		t.Fail()
	}
}

func BenchmarkQueueSendRecv(b *testing.B) {
	dummyMessagePayLoad := []byte("Test")
	dummyMessage := message.NewHeaderlessMessage(&dummyMessagePayLoad)

	underTest := NewQueue("Test")

	for i := 0; i < b.N; i++ {
		underTest.InputChannel <- dummyMessage
	}

	for i := 0; i < b.N; i++ {
		_ = <-underTest.OutputChannel
	}
}
