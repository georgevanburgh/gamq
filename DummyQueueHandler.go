package gamq

import "fmt"

type DummyMessageHandler struct {
}

func (dmh *DummyMessageHandler) Initialize(input chan string) chan string {
	outputChannel := make(chan string)
	go func() {
		message := <-input
		fmt.Println(message)
		outputChannel <- message
	}()
	return outputChannel
}
