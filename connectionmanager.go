package gamq

import (
	"bufio"
	"fmt"
	log "github.com/cihub/seelog"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	unrecognisedCommandText = "Unrecognised command"
	helpString              = `Available commands:
	HELP: Prints this text
	PUB <queue> <message>: Publish <message> to <queue>
	SUB <queue>: Subscribe to messages on <queue>
	PINGREQ: Requests a PINGRESP from the server`
)

type ConnectionManager struct {
	wg    sync.WaitGroup
	qm    QueueManager
	rand  *rand.Rand
	tcpLn net.Listener
	udpLn *net.Listener
}

func (manager *ConnectionManager) Initialize() {
	// Initialize our random number generator (used for naming new connections)
	// A different seed will be used on each startup, for no good reason
	manager.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	manager.qm = QueueManager{}
	manager.qm.Initialize()

	// Open TCP socket
	tcpAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", Configuration.Port)) // TODO: Handle error
	tcpListener, tcpErr := net.ListenTCP("tcp", tcpAddr)
	manager.tcpLn = tcpListener // Keep a a reference for unit testing purposes

	if tcpErr != nil {
		log.Criticalf("Error whilst opening TCP socket: %s", tcpErr.Error())
		panic(tcpErr.Error())
	}

	manager.wg.Add(1)
	go manager.listenOnConnection(tcpListener)

	// Listen on UDP socket
	udpAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", Configuration.Port))
	udpListener, _ := net.ListenUDP("udp", udpAddr)

	manager.wg.Add(1)
	go manager.listenOnUdpConnection(udpListener)

	manager.wg.Wait()
}

func (manager *ConnectionManager) listenOnUdpConnection(givenConnection *net.UDPConn) {
	var buffer [2048]byte

	// Listen forever
	for {
		length, remoteAddr, err := givenConnection.ReadFromUDP(buffer[0:])
		if err != nil {
			panic(err.Error())
		}

		writer := NewUdpWriter(givenConnection, remoteAddr)
		bufferedWriter := bufio.NewWriter(writer)

		client := NewClient(strconv.Itoa(manager.rand.Int()), bufferedWriter, nil)
		manager.parseClientCommand(string(buffer[:length]), client)

		log.Debugf("Read %d bytes from %s: %s", length, remoteAddr, string(buffer[:length]))
	}
}

func (manager *ConnectionManager) listenOnConnection(givenListener net.Listener) {
	defer manager.wg.Done()
	for {
		conn, err := givenListener.Accept()
		if err != nil {
			log.Errorf("An error occured whilst opening a TCP socket for reading: %s",
				err.Error())
		}
		log.Debug("A new TCP connection was opened.")
		manager.wg.Add(1)
		go manager.handleConnection(&conn)
	}
}

func (manager *ConnectionManager) handleConnection(conn *net.Conn) {
	defer manager.wg.Done()
	connReader := bufio.NewReader(*conn)
	connWriter := bufio.NewWriter(*conn)
	closedChannel := make(chan bool)
	client := Client{Name: strconv.Itoa(manager.rand.Int()),
		Writer: connWriter,
		Reader: connReader,
		Closed: &closedChannel}

	for {
		// Read until newline
		line, err := client.Reader.ReadString('\n')

		if err != nil {
			// Connection has been closed
			log.Debugf("%s closed connection", client.Name)
			*client.Closed <- true
			break
		}

		// Parse command and (optionally) return response (if any)
		manager.parseClientCommand(line, &client)
	}

	log.Info("A connection was closed")
}

func (manager *ConnectionManager) parseClientCommand(commandLine string, client *Client) {

	commandTokens := strings.Fields(string(commandLine[:len(commandLine)]))

	if len(commandTokens) == 0 {
		return
	}

	commandTokens[0] = strings.ToUpper(commandTokens[0])

	switch commandTokens[0] {
	case "HELP":
		manager.sendStringToClient(helpString, client)
	case "PUB":
		message := strings.Join(commandTokens[2:], " ")
		manager.qm.Publish(commandTokens[1], &message)
		// manager.sendStringToClient("PUBACK", client)
	case "SUB":
		manager.qm.Subscribe(commandTokens[1], client)
	case "DISCONNECT":
		*client.Closed <- true
	case "PINGREQ":
		manager.sendStringToClient("PINGRESP", client)
	case "CLOSE":
		manager.qm.CloseQueue(commandTokens[1])
	default:
		manager.sendStringToClient(unrecognisedCommandText, client)
	}
}

func (manager *ConnectionManager) sendStringToClient(toSend string, client *Client) {
	fmt.Fprintln(client.Writer, toSend)
	client.Writer.Flush()
}
