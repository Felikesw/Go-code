package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

func main() {

	cfg, err := File()
	if err != nil {
		log.Fatal("Failed to get data frame: ", err)
	}

	//connetcting to the serial port
	c := &serial.Config{Name: cfg.Port, Baud: cfg.BaudRate, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	defer s.Close()
	if err != nil {
		log.Fatal("Failed to open port: ", err)
	}

	scan := bufio.NewScanner(os.Stdin)
	var chunks []byte

	//see what the user want to do
	fmt.Println("\n\nPlease select an action (p = power on; r = read; w = write; q = quit): ")

	scan.Scan()
	action := scan.Text()

	for action != "q" {

		switch action {
		case "p":
			//power on, ReadPort will scan for 3 seconds, make sure the device is on before that
			log.Println("\nAwaiting for the device to power on...")

			chunks = ReadPort(s)
			if chunks == nil {
				log.Fatal("Did not received any data, check your connection")
			} else {
				log.Printf("%x", chunks)
				log.Println("Powered on data received")
			}

		case "r":
			log.Println("\nReading data...")
			chunks = ReadPort(s)
			if chunks == nil {
				log.Println("Did not received any data, check your connection")
			} else {
				log.Printf("Received data: %x", chunks)
			}

		case "w":
			WritePort(s, cfg)

			log.Println("Awaiting for the return data")
			chunks = ReadPort(s)
			if chunks == nil {
				log.Println("Did not received any data, check your connection")
			} else {
				log.Printf("Received data: %x", chunks)
			}

		case "q":
			os.Exit(1)

		default:
			log.Println("Invalid function code, please select one of the following: w, r, q")
		}

		fmt.Println("\n\nPlease select an action (r = read; w = write; q = quit): ")
		scan.Scan()
		action = scan.Text()
	}
}

//ReadPort reads data from the port
func ReadPort(port *serial.Port) []byte {
	var chunks []byte
	buf := make([]byte, 128)
	var err error
	n := 1

	for n != 0 {
		n, err = port.Read(buf)
		if err != nil && err != io.EOF {
			log.Println("Corrupted data: ", chunks)
			log.Fatal("Failed to read data: ", err)
		}

		for _, v := range buf[:n] {
			chunks = append(chunks, v)
		}
	}

	return chunks
}

//WritePort writes data to the port
func WritePort(port *serial.Port, cfg *Cfg) {
	//getting the data frame
	msg, err := Light(cfg)
	if err != nil {
		log.Fatal(err)
	}

	_, err = port.Write(msg)
	if err != nil {
		log.Fatal("Failed to write data: ", err)
	}
}
