import subprocess, socket, os
# Global variables
HostAddress = "localhost"
HostPort = 48879
Protocol = "tcp"
ImagePath = ""
UseWebcam = False

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
    global ImagePath, HostAddress, HostPort, Protocol, UseWebcam

    # Get benchmark parameters
    protocol = raw_input("Protocol to use (tcp/udp): ")

    if protocol not in ["tcp", "udp"]:
        print "Invalid protocol"
        exit(-1)
    else:
        Protocol = protocol

    useWebcam = raw_input("Use webcam? (y/n): ")
    UseWebcam = (useWebcam == "y")

    if not UseWebcam:
        imagePath = raw_input("Path of image to send: ")

        if not os.path.isfile(imagePath):
            print "Invalid path"
            exit(-1)
        else:
            ImagePath = imagePath

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

def readImageFromWebcam():
    subprocess.call(['imagesnap', '-w', str(1), '-q', '/tmp/photo.jpg'])
    imageData = readImageData('/tmp/photo.jpg')
    os.remove('/tmp/photo.jpg')
    return imageData

def readImageData(givenImagePath):
    try:
        file = open(givenImagePath, "rb")
        imageData = file.read()
    except:
        print "Error whilst reading file!"
    finally:
        if file:
            file.close()
    return imageData

def sendImageToBroker(givenImageData):
    s = getSocket()
    s.sendall("pub abc\n")
    s.sendall(givenImageData)
    s.sendall("\r\n")
    s.sendall(".\r\n")

readConfig()
if UseWebcam:
    while True:
        imageData = readImageFromWebcam()
        sendImageToBroker(imageData)
else:
    imageData = readImageData(ImagePath)
    sendImageToBroker(imageData)
