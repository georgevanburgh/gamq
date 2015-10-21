package encoders

import (
	"testing"
)

func TestEncoder_Encode(t *testing.T) {
	underTest := TestEncoder{}

	messageToEncode := Message{Body: "abc", Headers: make(map[string]string)}
	messageToEncode.Headers["Foo"] = "Bar"

	output := underTest.Encode(messageToEncode)
	if output != "Foo:Bar\nBODY\nabc" {
		t.Error("Encoded output incorrect")
	}
}

func TestEncoder_Decode(t *testing.T) {
	underTest := TestEncoder{}

	messageToDecode := "Foo:Bar\nBODY\nabc"

	decodedMessage := underTest.Decode(messageToDecode)

	if decodedMessage.Headers["Foo"] != "Bar" {
		t.Error("Header not decoded correctly")
	}

	if decodedMessage.Body != "abc" {
		t.Errorf("Body could not be decoded correctly")
	}
}

func BenchmarkEncoder_Encode(b *testing.B) {
	underTest := TestEncoder{}

	messageToEncode := Message{Body: "abc", Headers: make(map[string]string)}
	messageToEncode.Headers["Foo"] = "Bar"

	for n := 0; n < b.N; n++ {
		underTest.Encode(messageToEncode)
	}
}

func BenchmarkEncoder_Decode(b *testing.B) {
	underTest := TestEncoder{}

	messageToDecode := "Foo:Bar\nBODY\nabc"

	for n := 0; n < b.N; n++ {
		underTest.Decode(messageToDecode)
	}
}
