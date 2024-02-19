package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	log "github.com/sirupsen/logrus"

	"protocols/internal/certificate"
)

const (
	addr = "127.0.0.1:9999"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandle)
	mux.HandleFunc("/hello", helloHandle)
	cert, err := certificate.GetCertificate()
	if err != nil {
		log.Error(err, "certificate")
	}
	srv := http3.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*cert},
			ServerName:   "http3-server",
		},
		QuicConfig:      &quic.Config{},
		Handler:         mux,
		EnableDatagrams: true,
	}
	log.Info("Starting server...")
	if err = srv.ListenAndServe(); err != nil {
		log.Error(err)
	}
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	if _, err := fmt.Fprintf(w, "Main page!"); err != nil {
		log.Error(err, "mainHandle")
	}
	log.Info("Served main handler")
}

func helloHandle(w http.ResponseWriter, req *http.Request) {
	if _, err := fmt.Fprintf(w, "<h>Hello!</h1>"); err != nil {
		log.Error(err, "helloHandle")
	}
	log.Info("Served hello handler")
}
