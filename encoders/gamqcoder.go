package encoders

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type Encoder interface {
	Encode()
	Decode()
}

type TestEncoder struct {
}

//
func (t TestEncoder) Encode(m Message) string {
	var buffer bytes.Buffer

	for index, element := range m.Headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", index, element))
	}
	buffer.WriteString("BODY\n")

	buffer.WriteString(m.Body)

	return buffer.String()
}

func (t TestEncoder) Decode(s string) Message {
	scanner := bufio.NewScanner(bytes.NewBufferString(s))
	var toReturn Message = Message{}
	toReturn.Headers = make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "BODY" {
			break
		}
		if line != "" {
			parts := strings.SplitN(line, ":", 2)
			toReturn.Headers[parts[0]] = parts[1]
		}
	}

	scanner.Scan()
	toReturn.Body += scanner.Text()

	return toReturn
}

type Message struct {
	Body    string
	Headers map[string]string
}
