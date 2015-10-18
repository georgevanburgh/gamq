package main

import (
	"fmt"

	"github.com/FireEater64/gamq/encoders"
)

func main() {
	encoder := encoders.TestEncoder{}

	messageToEncode := encoders.Message{Body: "abc", Headers: make(map[string]string)}
	messageToEncode.Headers["Foo"] = "Bar"

	encodedMessage := encoder.Encode(messageToEncode)

	fmt.Println(encodedMessage)

	decodedMessage := encoder.Decode(encodedMessage)

	for key, value := range decodedMessage.Headers {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println(decodedMessage.Body)
}
