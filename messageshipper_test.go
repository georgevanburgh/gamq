package gamq

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/onsi/gomega"
)

func TestMessageShipper_SuccessfullyForwardsMessages(t *testing.T) {
	gomega.RegisterTestingT(t)

	underTest := MessageShipper{}

	inputChannel := make(chan *string, 0)

	writerBuffer := new(bytes.Buffer)
	dummyWriter := bufio.NewWriter(writerBuffer)
	closedChannel := make(chan bool)
	dummyClient := Client{Name: "Test", Writer: dummyWriter, Closed: &closedChannel}

	underTest.Initialize(inputChannel, &dummyClient)

	testMessage := "This is a test!"
	inputChannel <- &testMessage

	gomega.Eventually(func() string {
		return writerBuffer.String()
	}).Should(gomega.Equal(testMessage))
}
