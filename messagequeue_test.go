package gamq

import (
	"testing"
)

func TestMessageQueue_PushAndPullItem_ReturnsCorrectItem(t *testing.T) {
	underTest := NewMessageQueue()
	testString := "test"

	underTest.Push(&testString)
	result := underTest.Poll()

	if testString != *result {
		t.Fail()
	}
}

func TestMessageQueue_PushAndPeekAndPullItem_ReturnsCorrectItem(t *testing.T) {
	underTest := NewMessageQueue()
	testString := "test"

	underTest.Push(&testString)
	result1 := underTest.Peek()
	result2 := underTest.Poll()

	if testString != *result1 {
		t.Fail()
	}

	if testString != *result2 {
		t.Fail()
	}
}

func TestMessageQueue_Len_ReturnsCorrectQuantity(t *testing.T) {
	underTest := NewMessageQueue()
	testString := "test"
	messagesToPush := 10

	for i := 0; i < messagesToPush; i++ {
		underTest.Push(&testString)
	}

	reportedLength := underTest.Len()

	if int(reportedLength) != messagesToPush {
		t.Fail()
	}
}

func BenchmarkMessageQueue_Push(t *testing.B) {
	underTest := NewMessageQueue()
	testString := "test"

	for i := 0; i < t.N; i++ {
		underTest.Push(&testString)
	}
}

func BenchmarkMessageQueue_PushAndPop(t *testing.B) {
	underTest := NewMessageQueue()
	testString := "test"

	for i := 0; i < t.N; i++ {
		underTest.Push(&testString)
	}
	for i := 0; i < t.N; i++ {
		underTest.Poll()
	}
}
