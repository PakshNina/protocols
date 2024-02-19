package main

import (
	"net"

	log "github.com/sirupsen/logrus"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9393")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Remote address %v, local address %v", conn.RemoteAddr(), conn.LocalAddr())
	n, err := conn.Write([]byte("Hello"))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Send %d bytes", n)

	p := make([]byte, 5)
	n, err = conn.Read(p)
	log.Infof("Received %d bytes and message %s", n, string(p))
	if err = conn.Close(); err != nil {
		log.Fatal(err)
	}
	log.Info("Closed connection")
}
