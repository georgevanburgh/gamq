package gamq

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/FireEater64/gamq/message"
	"github.com/FireEater64/gamq/udp"
	log "github.com/cihub/seelog"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	unrecognisedCommandText   = "Unrecognised command"
	udpConnectionReadDeadline = 1
	helpString                = `Available commands:
	HELP: Prints this text
	PUB <queue> <message>: Publish <message> to <queue>
	SUB <queue>: Subscribe to messages on <queue>
	PINGREQ: Requests a PINGRESP from the server`
)

type ConnectionManager struct {
	wg         sync.WaitGroup
	qm         *queueManager
	rand       *rand.Rand
	tcpLn      net.Listener
	udpConn    *net.UDPConn
	udpClients map[string]*Client
	tcpClients int64
}

func NewConnectionManager() *ConnectionManager {
	manager := ConnectionManager{}

	// Initialize our random number generator (used for naming new connections)
	// A different seed will be used on each startup, for no good reason
	manager.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	manager.qm = newQueueManager()

	manager.udpClients = make(map[string]*Client)

	// Open TCP socket
	tcpAddr, tcpAddrErr := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", Configuration.Port))

	if tcpAddrErr != nil {
		panic("Invalid port configured")
	}

	tcpListener, tcpErr := net.ListenTCP("tcp", tcpAddr)
	manager.tcpLn = tcpListener

	if tcpErr != nil {
		log.Criticalf("Error whilst opening TCP socket: %s", tcpErr.Error())
		panic(tcpErr.Error())
	}

	// Listen on UDP socket
	udpAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", Configuration.Port))
	udpConn, _ := net.ListenUDP("udp", udpAddr)
	manager.udpConn = udpConn

	return &manager
}

func (manager *ConnectionManager) Start() {
	manager.wg.Add(2)
	go manager.listenOnUdpConnection()
	go manager.listenOnTcpConnection()
	manager.wg.Wait()
}

func (manager *ConnectionManager) listenOnUdpConnection() {
	var buffer [2048]byte

	// Listen forever
	// TODO: Revisit
	for {
		length, remoteAddr, err := manager.udpConn.ReadFromUDP(buffer[0:])
		if err != nil {
			panic(err.Error())
		}

		// Check if we've seen UDP packets from this address before - if so, reuse
		// existing client object
		client, ok := manager.udpClients[remoteAddr.String()]
		if !ok {
			log.Debug("New UDP client")
			writer := udp.NewUDPWriter(manager.udpConn, remoteAddr)
			bufferedWriter := bufio.NewWriter(writer)

			client = NewClient(strconv.Itoa(manager.rand.Int()), bufferedWriter, nil)
			manager.udpClients[remoteAddr.String()] = client
		} else {
			log.Debug("Found UDP client!")
		}

		// Log the number of bytes received
		manager.qm.metricsManager.metricsChannel <- NewMetric("bytesin.udp", "counter", int64(length))

		// TODO: Parse message, and check if we're expecting a message
		commandTokens := strings.Fields(string(buffer[:length]))
		var message []byte
		if commandTokens[0] == "pub" {
			// Use bytes.Equal until Go1.7 (https://github.com/golang/go/issues/14302)
			for {
				var err error
				length, _, err := manager.udpConn.ReadFromUDP(buffer[0:])

				if err != nil {
					return
				}

				// TODO: Is this cross platform? Needs testing
				if !bytes.Equal(buffer[:length], []byte{'.', '\r', '\n'}) {
					message = append(message, buffer[:length]...)
				} else {
					break
				}
			}
		}

		manager.parseClientCommand(commandTokens, &message, client)

		log.Debugf("Read %d bytes from %s: %s", length, remoteAddr, string(buffer[:length]))
	}
}

func (manager *ConnectionManager) listenOnTcpConnection() {
	defer manager.wg.Done()
	for {
		conn, err := manager.tcpLn.Accept()
		if err != nil {
			log.Errorf("An error occured whilst opening a TCP socket for reading: %s",
				err.Error())
		}
		log.Debug("A new TCP connection was opened.")

		manager.tcpClients++
		manager.updateClientMetric()

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

			*client.Closed <- true // TODO: This is blocking - shouldn't be
			break
		}

		// Tokenise the command line
		commandTokens := strings.Fields(string(line[:]))

		if len(commandTokens) == 0 {
			break
		}

		var message []byte

		if commandTokens[0] == "pub" {

			var line []byte

			// Use bytes.Equal until Go1.7 (https://github.com/golang/go/issues/14302)
			for {
				var err error
				line, err = client.Reader.ReadBytes(byte('\n'))

				if err != nil {
					return
				}

				// TODO: Is this cross platform? Needs testing
				if !bytes.Equal(line, []byte{'.', '\r', '\n'}) {
					message = append(message, line...)
				} else {
					log.Debug("End of message")
					break
				}
			}
		}

		// Parse command and (optionally) return response (if any)
		manager.parseClientCommand(commandTokens, &message, &client)

		// Log the number of bytes received (command + body)
		manager.qm.metricsManager.metricsChannel <- NewMetric("bytesin.tcp", "counter", int64(len(line)+len(message)))
	}

	log.Info("A connection was closed")
	manager.tcpClients--
	manager.updateClientMetric()
}

func (manager *ConnectionManager) updateClientMetric() {
	manager.qm.metricsManager.metricsChannel <- NewMetric("clients.tcp", "guage", manager.tcpClients)
}

func (manager *ConnectionManager) parseClientCommand(commandTokens []string, messageBody *[]byte, client *Client) {
	if len(commandTokens) == 0 {
		return
	}

	commandTokens[0] = strings.ToUpper(commandTokens[0])

	switch commandTokens[0] {
	case "HELP":
		manager.sendStringToClient(helpString, client)
	case "PUB":
		// TODO: Handle headers
		message := message.NewHeaderlessMessage(messageBody)
		manager.qm.Publish(commandTokens[1], message)
		if client.AckRequested {
			manager.sendStringToClient("PUBACK\n", client)
		}
	case "SUB":
		manager.qm.Subscribe(commandTokens[1], client)
	case "DISCONNECT":
		*client.Closed <- true
	case "PINGREQ":
		manager.sendStringToClient("PINGRESP", client)
	case "SETACK":
		if strings.ToUpper(commandTokens[1]) == "ON" {
			client.AckRequested = true
		} else {
			client.AckRequested = false
		}
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
