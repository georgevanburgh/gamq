package queue

import (
	"github.com/FireEater64/gamq/message"
	"github.com/onsi/gomega"
	"testing"
)

func TestTopic_SingleSubscriber_SendsMessageCorrectly(t *testing.T) {
	underTest := NewTopic("test")

	outputChannel := underTest.GetChannel("foobar")

	testMessagePayload := []byte("TestMessage")
	testMessage := message.NewHeaderlessMessage(&testMessagePayload)
	underTest.InputChannel <- testMessage

	receivedMessage := <-outputChannel

	if receivedMessage != testMessage {
		t.Fail()
	}
}

func TestTopic_MultipleSubscribers_SendsMessagesCorrectly(t *testing.T) {
	gomega.RegisterTestingT(t)

	underTest := NewTopic("test")

	outputChannel1 := underTest.GetChannel("foo")
	outputChannel2 := underTest.GetChannel("bar")

	testMessagePayload := []byte("TestMessage")
	testMessage := message.NewHeaderlessMessage(&testMessagePayload)
	underTest.InputChannel <- testMessage

	gomega.Eventually(func() *message.Message {
		receivedMessage := <-outputChannel1
		return receivedMessage
	}).Should(gomega.Equal(testMessage))

	gomega.Eventually(func() *message.Message {
		receivedMessage := <-outputChannel2
		return receivedMessage
	}).Should(gomega.Equal(testMessage))
}
