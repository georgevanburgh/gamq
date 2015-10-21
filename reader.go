package gamq

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

const (
	PORT = 4321
)

type Reader interface {
	Initialize()
}

type TcpReader struct {
	Completed sync.WaitGroup
}

func (reader TcpReader) Initialize() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	if err != nil {
		fmt.Printf("An error occured whilst opening a socket for reading: %s",
			err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("An error occured whilst opening a socket for reading: %s",
				err.Error())
		}

		channel := make(chan string)
		reader.Completed.Add(1)
		go handleConnection(channel, &conn, &reader)
	}
}

func handleConnection(out chan string, conn *net.Conn, reader *TcpReader) {
	defer close(out)
	defer reader.Completed.Done()

	bufferedReader := bufio.NewReader(*conn)

	for {
		line, err := bufferedReader.ReadBytes('\n')

		if err != nil {
			// Connection has been closed
			return
		}
		fmt.Println(line)
	}
}
