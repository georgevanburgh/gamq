package gamq

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestConnectionManager_parseClientCommand_helpMessageReturns(t *testing.T) {
	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	underTest.parseClientCommand("HELP", &mockClient)

	if buf.String() == UNRECOGNISEDCOMMANDTEXT {
		t.Fail()
	}
}

func TestConnectionManager_parseClientCommand_isCaseInsensitive(t *testing.T) {
	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	underTest.parseClientCommand("help", &mockClient)

	if buf.String() == UNRECOGNISEDCOMMANDTEXT {
		t.Fail()
	}
}

func TestConnectionManager_parseClientCommand_invalidCommandProcessedCorrectly(t *testing.T) {
	underTest := ConnectionManager{}

	dummyWriterBuffer := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(dummyWriterBuffer)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	underTest.parseClientCommand("fdkfjadkfh", &mockClient)

	if dummyWriterBuffer.String() != UNRECOGNISEDCOMMANDTEXT+"\n" {
		t.Fail()
	}
}

func TestConnectionManager_emptyClientCommand_handledGracefully(t *testing.T) {
	underTest := ConnectionManager{}

	mockClient := Client{}

	underTest.parseClientCommand("", &mockClient)
}

func TestConnectionManager_whenInitialized_acceptsConnectionsCorrectly(t *testing.T) {
	// Choose a high port, so we don't need sudo to run tests
	config := Config{}
	config.Port = 55555
	SetConfig(&config)

	underTest := ConnectionManager{}
	go underTest.Initialize()

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
