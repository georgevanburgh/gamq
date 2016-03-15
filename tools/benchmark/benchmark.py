#!/usr/bin/env python

# Python benchmark for gamq

import time
import socket
import threading

# Global variables
HostAddress = "localhost"
HostPort = 48879
Protocol = ""
AckMessages = False
NumberOfMessages = 0

# Helper function to check if a number is valid
def isNumber(givenObject):
    try:
        int(givenObject)
        return True
    except:
        return False


def getSocket(protocol):
    if protocol == "tcp":
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    elif protocol == "udp":
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    else:
        print "Invalid protocol: {}".format(protocol)
        exit(-1)
    s.connect((HostAddress, HostPort))
    return s


def writeThread():
    s = getSocket(Protocol)

    if AckMessages:
        s.sendall("setack on\n")

    startTime = time.clock()
    for i in range(0, int(NumberOfMessages), 1):
        s.sendall("pub abc\n")
        s.sendall("{}\n".format(i))
        s.sendall(".\r\n")
        if AckMessages:
            response = s.recv(8)
            if response[:6] != "PUBACK":
                print "Error whilst publishing {}, got response: {}".format(i, response)
    endTime = time.clock()
    print "Took {} seconds to write {} messages".format((endTime - startTime), NumberOfMessages)
    time.sleep(2)
    s.close()


def readThread():
    s = getSocket("tcp")
    startTime = time.clock()
    s.sendall("sub abc\n")

    for i in range(0, int(NumberOfMessages), 1):
        response = ""
        while response[-3:] != ".\r\n":
            response += s.recv(1)
        response = response.translate(None, ".\r\n")
        if int() != int(i):
            print "Expected {}, got {}".format(i, response)
    endTime = time.clock()
    print "Took {} seconds to read {} messages".format((endTime - startTime), NumberOfMessages)
    time.sleep(2)
    s.close()


def readConfig():
    global AckMessages, NumberOfMessages, HostAddress, HostPort, Protocol

    # Get benchmark parameters
    protocol = raw_input("Protocol to use (tcp/udp): ")

    if protocol not in ["tcp", "udp"]:
        print "Invalid protocol"
        exit(-1)
    else:
        Protocol = protocol

    numberOfMessages = raw_input("Number of messages to send: ")

    if not isNumber(numberOfMessages):
        print "Invalid number"
        exit(-1)
    else:
        NumberOfMessages = int(numberOfMessages)

    ackMessages = raw_input("Ack messages (y/n): ")

    AckMessages = (ackMessages == "y")

    hostAddress = raw_input("Host to connect to: ")

    if hostAddress == "":
        print "Defaulting to localhost"
    else:
        HostAddress = hostAddress

    hostPort = raw_input("Port to connect to: ")

    if hostPort == "":
        print "Defaulting to 48879"
    elif isNumber(hostPort):
        HostPort = hostPort
    else:
        print "Invalid number"
        exit(-1)

readConfig()
writeThread = threading.Thread(target=writeThread)
readThread = threading.Thread(target=readThread)
readThread.daemon = True
writeThread.daemon = True
writeThread.start()
readThread.start()

while threading.active_count() > 1:
    time.sleep(1)
