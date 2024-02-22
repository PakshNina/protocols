package main

import (
	"crypto/tls"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
)

func main() {
	client := &http.Client{}
	tlsConfig := &tls.Config{
		MaxVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}
	client.Transport = &http2.Transport{
		TLSClientConfig: tlsConfig,
	}
	resp, err := client.Get("https://localhost:9090/hello")
	if err != nil {
		log.Error(err, "GET")
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()
	if err != nil {
		log.Fatalf("Failed reading response body: %s", err)
	}
	log.Infof("Got response %d: %s %s", resp.StatusCode, resp.Proto, string(body))
}
