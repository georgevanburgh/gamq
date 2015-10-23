package gamq

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

const (
	TCPPORT                 = 48879
	HELPSTRING              = "Help. Me.\n"
	UNRECOGNISEDCOMMANDTEXT = "Unrecognised command\n"
)

type ConnectionManager struct {
	wg sync.WaitGroup
	qm QueueManager
}

func (manager *ConnectionManager) Initialize() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", TCPPORT))

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
		manager.wg.Add(1)
		go manager.handleConnection(&conn)
	}
}

func (manager *ConnectionManager) handleConnection(conn *net.Conn) {
	defer manager.wg.Done()
	connReader := bufio.NewReader(*conn)
	connWriter := bufio.NewWriter(*conn)

	for {
		line, err := connReader.ReadBytes('\n')

		if err != nil {
			// Connection has been closed
			break
		}

		stringLine := string(line[:len(line)])
		fmt.Print(stringLine)

		// Parse command and (optionally) return response (if any)
		response := manager.parseClientCommand(stringLine, connWriter)
		if response != "" {
			connWriter.WriteString(response)
			connWriter.Flush()
		}
	}

	fmt.Println("Connection closed")
}

func (manager *ConnectionManager) parseClientCommand(command string, writer io.Writer) string {
	commandTokens := strings.Fields(command)
	switch strings.ToUpper(commandTokens[0]) {
	case "HELP":
		return HELPSTRING
	case "PUB":
		manager.qm.Publish(commandTokens[1], strings.Join(commandTokens[2:], " "))
	case "SUB":
		manager.qm.Subscribe(commandTokens[1], &writer)
	default:
		return UNRECOGNISEDCOMMANDTEXT
	}

	return ""
}
