package main

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func main() {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP)
	if err != nil {
		log.Fatal("Can not create socket", err)
	}
	log.Info("Created socket")
	if err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		log.Fatal("Can not set opt", err)
	}
	if err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 1); err != nil {
		log.Fatal("Can not set opt", err)
	}
	log.Info("Set options")
	sa := &unix.SockaddrInet4{
		Port: 9393,
		Addr: [4]byte{127, 0, 0, 1},
	}
	if err = unix.Bind(fd, sa); err != nil {
		log.Fatal("could not bind address with socker")
	}
	log.Infof("Bind socket to address %v:%d", sa.Addr, sa.Port)
	if err = unix.Listen(fd, 10); err != nil {
		log.Fatal("Could not listen")
	}
	log.Info("Start listing to incoming connections")
	for {
		fdc, _, err := unix.Accept(fd)
		if err != nil {
			continue
		}
		log.Infof("Accepted client with descriptor %d", fdc)
		p := make([]byte, 5)
		n, err := unix.Read(fdc, p)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("got %d bytes: %s", n, string(p))
		n, err = unix.Write(fdc, p)
		log.Infof("Send %d bytes: %s", n, string(p))
	}
}
