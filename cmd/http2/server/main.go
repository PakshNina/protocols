package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"

	c "protocols/internal/certificate"
)

// curl -k -v --http1.1 https://localhost:9090/hello
func main() {
	cert, err := c.GetCertificate()
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{
		MaxVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{*cert},
	}

	httpServer := http.Server{
		Addr:      ":9090",
		TLSConfig: tlsConfig,
	}
	http2Server := http2.Server{}
	_ = http2.ConfigureServer(&httpServer, &http2Server)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(httpServer.ListenAndServeTLS(
		"", "",
	))
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	log.Infof("Request connection: %s, path: %s", req.Proto, req.URL.Path[1:])
	w.Header().Set("Content-Type", "text/html")
	if _, err := fmt.Fprintf(w, "Hello!"); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

