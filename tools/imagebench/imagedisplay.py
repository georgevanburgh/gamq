import socket, os
import tempfile
import subprocess

# Global variables
HostAddress = "localhost"
HostPort = 48879
Protocol = "tcp"

# Helper function to check if a number is valid
def isNumber(givenObject):
    try:
        int(givenObject)
        return True
    except:
        return False

def getSocket():
    global HostAddress, HostPort
    if Protocol == "tcp":
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    else:
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect((HostAddress, HostPort))
    return s

def readConfig():
    global HostAddress, HostPort, Protocol

    # Get benchmark parameters
    protocol = raw_input("Protocol to use (tcp/udp): ")

    if protocol not in ["tcp", "udp"]:
        print "Invalid protocol"
        exit(-1)
    else:
        Protocol = protocol

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

def readImage():
    s = getSocket()
    s.sendall("sub abc\n")

    image = ""
    line = ""

    while line[-5:] != "\r\n.\r\n":
        print "Reading chunk"
        line = s.recv(4096)
        image += line

    return image

def displayImage(imageData):
    tempFile = tempfile.NamedTemporaryFile(suffix=".jpg", delete=False)
    tempFile.write(imageData)
    tempFile.close()
    viewer = subprocess.Popen(['display', tempFile.name])
    return viewer


readConfig()
displayedImage = ""
openImage = False
while True:
    imageData = readImage()
    if openImage:
        openImage.terminate()
        openImage.kill()
    openImage = displayImage(imageData)
