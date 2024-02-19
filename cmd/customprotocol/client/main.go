package main

import (
	"bufio"
	"crypto/tls"
	"unsafe"

	log "github.com/sirupsen/logrus"

	p "protocols/internal/customprotocol"
)

const (
	protocol = "custom-protocol"
	network  = "tcp"
	addr     = "127.0.0.1:9494"
)

func main() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos: []string{
			protocol,
		},
	}
	conn, err := tls.Dial(network, addr, tlsConfig)
	defer func() {
		if err = conn.Close(); err != nil {
			log.Error(err)
		}
	}()
	if err != nil {
		log.Error(err)
	}

	rdr := bufio.NewReader(conn)
	wtr := bufio.NewWriter(conn)

	// Creating bytes to send
	msg := p.CustomMessage{
		Field1: 0,
		Field2: 2,
		Field3: 3,
	}
	ptr := (*byte)(unsafe.Pointer(&msg))
	bytes := unsafe.Slice(ptr, unsafe.Sizeof(msg))
	n, err := wtr.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
	if err = wtr.Flush(); err != nil {
		log.Fatal(err)
	}
	log.Infof("Send %d bytes to the '%s'", n, protocol)
	if err = conn.CloseWrite(); err != nil {
		log.Fatal(err)
	}

	var resp p.CustomResponse
	responseBytes := make([]byte, unsafe.Sizeof(resp))
	n, err = rdr.Read(responseBytes)
	if err != nil {
		log.Fatal(err)
	}
	result := (*p.CustomResponse)(unsafe.Pointer(&responseBytes[0]))
	log.Infof("Received %d bytes from the server with the response Status: %v", n, result.Status)
}
