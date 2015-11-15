package gamq

import (
	"testing"
)

const (
	TEST_QUEUE_NAME = "TestQueue"
)

func TestQueue_initialize_completesSuccessfully(t *testing.T) {
	underTest := Queue{Name: TEST_QUEUE_NAME}

	underTest.Initialize()

	// Queue should be named correctly
	if underTest.Name != TEST_QUEUE_NAME {
		t.Fail()
	}

	// Messages channel should be initialized
	if underTest.Messages == nil {
		t.Fail()
	}

	// Subscribers channel should be initialized
	if underTest.Subscribers == nil {
		t.Fail()
	}
}
