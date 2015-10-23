package gamq

import (
	"bufio"
	"fmt"
)

type QueueManager struct {
}

func (qm QueueManager) Initialize() {

}

func (qm QueueManager) Publish(queueName string, message string) {
	fmt.Printf("Publishing message to %s: %s", queueName, message)
}

func (qm QueueManager) Subscribe(queueName string, writer *bufio.Writer) {
	fmt.Printf("Subscribing to %s", queueName)
}
