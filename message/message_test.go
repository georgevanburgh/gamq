package message

import (
	"fmt"
	"testing"
)

func Benchmark(b *testing.B) {
	var message *Message

	dummyHeaders := make(map[string]string)
	dummyBody := []byte("abc123")

	for i := 0; i < b.N; i++ {
		message = NewMessage(&dummyHeaders, &dummyBody)
	}

	fmt.Println(message)
}
