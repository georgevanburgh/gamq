package main

import (
	"github.com/FireEater64/gamq"
)

func main() {
	// Set up a done channel, that's shared by the whole pipeline.
	// Closing this channel will kill all pipeline goroutines
	//done := make(chan struct{})
	//defer close(done)

	connectionManager := gamq.ConnectionManager{}
	connectionManager.Initialize()
}
