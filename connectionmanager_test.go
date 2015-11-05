package gamq

import (
	"os"
	"testing"
)

func TestConnectionManager_parseClientCommand_helpMessageReturns(t *testing.T) {
	underTest := ConnectionManager{}

	response := underTest.parseClientCommand("HELP", os.Stdout.Write())

	if response == UNRECOGNISEDCOMMANDTEXT {
		t.Fail()
	}
}

func TestConnectionManager_parseClientCommand_isCaseInsensitive(t *testing.T) {
	underTest := ConnectionManager{}

	response := underTest.parseClientCommand("help")

	if response == UNRECOGNISEDCOMMANDTEXT {
		t.Fail()
	}
}
