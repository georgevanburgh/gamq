package gamq

import (
	"testing"
)

func TestConnectionManager_parseClientCommand_helpMessageReturns(t *testing.T) {
	underTest := ConnectionManager{}

	response := underTest.parseClientCommand("HELP")

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
