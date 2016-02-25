package message

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewMessage_CreatesSuccessfully(t *testing.T) {
	testPayload := []byte("TestMessage")
	testHeader := make(map[string]string)
	testHeader["abc"] = "123"

	underTest := NewMessage(&testHeader, &testPayload)

	if (*underTest.Head)["abc"] != "123" {
		t.Fail()
	}

	if !bytes.Equal(testPayload, *underTest.Body) {
		t.Fail()
	}
}

func TestNewHeaderlessMessage_CreatesSuccessfully(t *testing.T) {
	testPayload := []byte("TestMessage")

	underTest := NewHeaderlessMessage(&testPayload)

	if (*underTest.Head) == nil {
		t.Fail()
	}

	if !bytes.Equal(testPayload, *underTest.Body) {
		t.Fail()
	}
}

func Benchmark(b *testing.B) {
	var message *Message

	dummyHeaders := make(map[string]string)
	dummyBody := []byte("abc123")

	for i := 0; i < b.N; i++ {
		message = NewMessage(&dummyHeaders, &dummyBody)
	}

	fmt.Println(message)
}
