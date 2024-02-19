package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	log "github.com/sirupsen/logrus"
)

const (
	addr = "https://127.0.0.1:9999"
)

func main() {
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		QuicConfig: &quic.Config{},
	}
	defer func() {
		if err := roundTripper.Close(); err != nil {
			log.Error(err, "roundTripper")
		}
	}()
	hclient := &http.Client{
		Transport: roundTripper,
	}
	resp, err := hclient.Get(fmt.Sprintf("%s%s", addr, "/"))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(resp.Body)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()
	log.Infof("Response for main page: '%s'", string(b))
	resp, err = hclient.Get(fmt.Sprintf("%s%s", addr, "/hello"))
	b, err = io.ReadAll(resp.Body)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()
	log.Infof("Response for main page '%s'", string(b))
}
