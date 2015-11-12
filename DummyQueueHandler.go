package gamq

import "fmt"

type DummyMessageHandler struct {
}

func (dmh *DummyMessageHandler) Initialize(input chan string) chan string {
	outputChannel := make(chan string)
	go func() {
		for {
			message, more := <-input
			if more {
				fmt.Println(message)
				outputChannel <- message
			} else {
				// We're done
				fmt.Println("Closing output channel")
				close(outputChannel)
				return
			}
		}
	}()
	return outputChannel
}
