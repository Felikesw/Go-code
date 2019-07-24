package main

import (
	"log"

	"github.com/tarm/serial"
)

func main() {

	cfg, err := File()
	if err != nil {
		log.Fatal(err)
	}

	msg, err := Light(cfg)
	if err != nil {
		log.Fatal(err)
	}

	c := &serial.Config{Name: cfg.Port, Baud: cfg.BaudRate}
	s, err := serial.OpenPort(c)
	defer s.Close()

	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write(msg)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data received: %x", buf[:n])
}
