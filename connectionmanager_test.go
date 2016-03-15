package gamq

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/onsi/gomega"
)

func TestConnectionManager_parseClientCommand_helpMessageReturns(t *testing.T) {
	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	var emptyMessage []byte

	underTest.parseClientCommand([]string{"HELP"}, &emptyMessage, &mockClient)

	if buf.String() == unrecognisedCommandText {
		t.Fail()
	}
}

func TestConnectionManager_parseClientCommand_isCaseInsensitive(t *testing.T) {
	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	var emptyMessage []byte

	underTest.parseClientCommand([]string{"help"}, &emptyMessage, &mockClient)

	if buf.String() == unrecognisedCommandText {
		t.Fail()
	}
}

func TestConnectionManager_setAck_setsClientAckFlag(t *testing.T) {
	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	var emptyMessage []byte

	underTest.parseClientCommand([]string{"setack", "on"}, &emptyMessage, &mockClient)

	if !mockClient.AckRequested {
		t.Fail()
	}
}

func TestConnectionManager_parseEmptyCommand_doesNotCrash(t *testing.T) {
	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	var emptyMessage []byte

	underTest.parseClientCommand([]string{}, &emptyMessage, &mockClient)

	if buf.String() != "" {
		t.Fail()
	}
}

func TestConnectionManager_subscribeToQueue_addsClientToListOfSubscribers(t *testing.T) {

	// Choose a high port, so we don't need sudo to run tests
	config := Config{}
	config.Port = 55556
	SetConfig(&config)

	underTest := NewConnectionManager()

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := NewClient("Mock", bufWriter, nil)

	var emptyMessage []byte

	underTest.parseClientCommand([]string{"sub", "test"}, &emptyMessage, mockClient)

	queue := underTest.qm.getQueueSafely("test")

	if len(queue.subscribers) != 1 {
		t.Fail()
	}

	*mockClient.Closed <- true
}

func TestConnectionManager_disconnectCommand_removesClient(t *testing.T) {
	gomega.RegisterTestingT(t)

	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := NewClient("Mock", bufWriter, nil)
	closedChannel := make(chan bool, 1)
	mockClient.Closed = &closedChannel

	t.Log("Disconnecting")

	var emptyMessage []byte

	underTest.parseClientCommand([]string{"disconnect"}, &emptyMessage, mockClient)

	gomega.Eventually(func() bool {
		closed := <-*(mockClient.Closed)
		return closed
	}).Should(gomega.BeTrue())
}

func TestConnectionManager_parseClientCommand_invalidCommandProcessedCorrectly(t *testing.T) {
	underTest := ConnectionManager{}

	dummyWriterBuffer := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(dummyWriterBuffer)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	var emptyMessage []byte

	underTest.parseClientCommand([]string{"fdkfjadkfh"}, &emptyMessage, &mockClient)

	if dummyWriterBuffer.String() != unrecognisedCommandText+"\n" {
		t.Fail()
	}
}

func TestConnectionManager_emptyClientCommand_handledGracefully(t *testing.T) {
	underTest := ConnectionManager{}

	dummyWriterBuffer := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(dummyWriterBuffer)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	var emptyMessage []byte

	underTest.parseClientCommand([]string{""}, &emptyMessage, &mockClient)

	if dummyWriterBuffer.String() != unrecognisedCommandText+"\n" {
		t.Fail()
	}
}

func TestConnectionManager_whenInitialized_acceptsConnectionsCorrectly(t *testing.T) {
	gomega.RegisterTestingT(t)

	// Choose a high port, so we don't need sudo to run tests
	config := Config{}
	config.Port = 55555
	SetConfig(&config)

	underTest := NewConnectionManager()
	go underTest.Start()

	gomega.Eventually(func() net.Listener {
		return underTest.tcpLn
	}).ShouldNot(gomega.BeNil())

	testConn, err := net.Dial("tcp", "localhost:55555")
	if err != nil || testConn == nil {
		t.Fail()
	}

	fmt.Fprintf(testConn, "PINGREQ\n")
	response, err := bufio.NewReader(testConn).ReadString('\n')

	if err != nil || response != "PINGRESP\n" {
		t.Fail()
	}

	testConn.Close()
}
