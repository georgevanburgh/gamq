package gamq

import (
	"bufio"
	"bytes"
	"testing"
)

func TestNewClient_ReturnsValidClient(t *testing.T) {

	testName := "TestClient"
	buf := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buf)
	bufReader := bufio.NewReader(buf)

	underTest := NewClient(testName, bufWriter, bufReader)

	if underTest.Closed == nil {
		t.Error("Closed channel nil")
		t.Fail()
	}

	if underTest.Name != testName {
		t.Errorf("Expected name to be %s, received %s", testName, underTest.Name)
		t.Fail()
	}

	if underTest.Writer != bufWriter {
		t.Error("Writer not set correctly")
		t.Fail()
	}

	if underTest.Reader != bufReader {
		t.Error("Reader not set correctly")
		t.Fail()
	}
}
