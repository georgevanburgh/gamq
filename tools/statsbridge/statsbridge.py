import socket

UDP_IP = "127.0.0.1"
UDP_PORT = 8125
TCP_IP = "127.0.0.1"
TCP_PORT = 22222

incommingSocket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
incommingSocket.bind((UDP_IP, UDP_PORT))

outgoingSocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
outgoingSocket.connect((TCP_IP, TCP_PORT))

print "Listening on: {}".format(UDP_PORT)

while True:
    data, addr = incommingSocket.recvfrom(1024)
    outgoingSocket.sendall(data)
    print "Relayed: {}".format(data)
