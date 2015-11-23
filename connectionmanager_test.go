package gamq

import (
	"bufio"
	"bytes"
	"testing"
)

func TestConnectionManager_parseClientCommand_helpMessageReturns(t *testing.T) {
	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	underTest.parseClientCommand([]string{"HELP"}, &mockClient)

	if buf.String() == UNRECOGNISEDCOMMANDTEXT {
		t.Fail()
	}
}

func TestConnectionManager_parseClientCommand_isCaseInsensitive(t *testing.T) {
	underTest := ConnectionManager{}

	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	mockClient := Client{Name: "Mock", Writer: bufWriter}

	underTest.parseClientCommand([]string{"help"}, &mockClient)

	if buf.String() == UNRECOGNISEDCOMMANDTEXT {
		t.Fail()
	}
}
