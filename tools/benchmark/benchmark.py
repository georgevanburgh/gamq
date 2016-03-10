#!/usr/bin/env python

# Python benchmark for gamq

import time
import socket
import threading

# Global variables
HostAddress = ""
HostPort = 48879
Protocol = ""
NumberOfMessages = 0

# Helper function to check if a number is valid
def isNumber(givenObject):
    try:
        int(givenObject)
        return True
    except:
        return False


def getSocket():
    if Protocol == "tcp":
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    # else:
    #     s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect((HostAddress, HostPort))
    return s


def writeThread():
    s = getSocket()
    startTime = time.clock()
    for i in range(0, int(NumberOfMessages), 1):
        command = "pub abc {}\n".format(i)
        print command,
        s.sendall(command)
    endTime = time.clock()
    s.close()
    print "Took {} seconds to write {} messages".format((endTime - startTime), NumberOfMessages)


def readThread():
    s = getSocket()
    startTime = time.clock()
    s.sendall("sub abc\n")
    for i in range(0, int(NumberOfMessages), 1):
        response = s.recv(10)
        print response
        if int(response) != int(i):
            print "Expected {}, got {}".format(i, response)
    endTime = time.clock()
    s.close()
    print "Took {} seconds to read {} messages".format((endTime - startTime), NumberOfMessages)


def readConfig():
    global NumberOfMessages, HostAddress, HostPort, Protocol

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
readThread = threading.Thread(target=readThread)
writeThread = threading.Thread(target=writeThread)
readThread.start()
writeThread.start()
readThread.join()
writeThread.join()
