package main

import (
	"github.com/FireEater64/gamq"
)

func main() {
	reader := gamq.TcpReader{}
	reader.Initialize()
	reader.Completed.Wait()
}
