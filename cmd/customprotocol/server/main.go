package main

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"unsafe"

	log "github.com/sirupsen/logrus"

	"protocols/internal/certificate"
	p "protocols/internal/customprotocol"
)

const (
	protocol = "custom-protocol"
	network  = "tcp"
	addr     = "127.0.0.1:9494"
	chunk    = 8
)

func main() {
	cert, err := certificate.GetCertificate()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Created server certificate")
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ServerName:   "custom-server",
		NextProtos: []string{
			protocol,
		},
	}
	log.Infof("Created tls config for %s protocol", protocol)
	conn, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Started listen address %s", addr)
	srv := &http.Server{
		Addr: addr,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){
			protocol: handleCustomProtocol,
		},
	}
	if err = srv.Serve(tls.NewListener(conn, tlsConfig)); err != nil {
		log.Fatal(err)
	}
}

func handleCustomProtocol(_ *http.Server, conn *tls.Conn, _ http.Handler) {
	rdr := bufio.NewReader(conn)
	wtr := bufio.NewWriter(conn)
	buffer := make([]byte, chunk)
	var receivedBytes []byte
	for {
		for {
			_, err := rdr.Read(buffer)
			receivedBytes = append(receivedBytes, buffer...)
			if err != nil {
				if err != io.EOF {
					log.Error(err)
				}
				break
			}
		}
		receivedMessage := (*p.CustomMessage)(unsafe.Pointer(&receivedBytes[0]))
		log.Infof("Got frame from client: Field1 %v, Field2 %v, Field3 %v", receivedMessage.Field1, receivedMessage.Field2, receivedMessage.Field3)

		response := handleMessage(receivedMessage)
		addrByte := (*byte)(unsafe.Pointer(response))
		responseBytes := unsafe.Slice(addrByte, unsafe.Sizeof(response.Status))
		if _, err := wtr.Write(responseBytes); err != nil {
			log.Error(err)
			continue
		}
		err := wtr.Flush()
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Response to client sent with Status %v", response.Status)

		// Waiting another message
		_, err = rdr.Peek(1)
		if err != nil {
			if err == bufio.ErrBufferFull {
				continue
			}
			break
		}
	}
}

func handleMessage(message *p.CustomMessage) *p.CustomResponse {
	response := &p.CustomResponse{}
	if message.Field1 > 0 && message.Field2 > 0 && message.Field3 > 0 {
		response.Status = 0
	} else {
		response.Status = 1
	}
	return response
}
