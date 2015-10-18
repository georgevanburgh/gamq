package main
import (
	"github.com/FireEater64/gamq/encoders"
	"fmt"
)

func main() {
	encoder := encoders.TestEncoder{}
	encodedMessage := encoder.Encode(encoders.Message{Body: "abc", Headers: 123 })

	fmt.Println(encodedMessage)

	decodedMessage := encoder.Decode(encodedMessage)
	fmt.Println(decodedMessage.Body)
}
