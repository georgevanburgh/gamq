package gamq

import (
	"bufio"
	"bytes"
	"github.com/FireEater64/gamq/message"
	"testing"

	"github.com/onsi/gomega"
)

func TestMessageShipper_SuccessfullyForwardsMessages(t *testing.T) {
	gomega.RegisterTestingT(t)

	underTest := MessageShipper{}

	inputChannel := make(chan *message.Message, 0)

	writerBuffer := new(bytes.Buffer)
	dummyWriter := bufio.NewWriter(writerBuffer)
	closedChannel := make(chan bool)
	dummyClient := Client{Name: "Test", Writer: dummyWriter, Closed: &closedChannel}

	underTest.Initialize(inputChannel, &dummyClient)

	testMessagePayload := []byte("This is a test!")
	testMessage := message.NewHeaderlessMessage(&testMessagePayload)
	inputChannel <- testMessage

	gomega.Eventually(func() []byte {
		return writerBuffer.Bytes()
	}).Should(gomega.Equal(testMessagePayload))
}
