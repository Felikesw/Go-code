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

	// err := Edit()
	// if err != nil {
	// 	log.Fatal("Failed to edit cfg.json")
	// }

	cfg, err := File()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	//connetcting to the serial port
	c := &serial.Config{Name: cfg.Port, Baud: cfg.BaudRate, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal("Failed to open port: ", err)
	}
	defer s.Close()

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
			log.Print("\nAwaiting for the device to power on...")

			chunks = ReadPort(s)
			if chunks == nil {
				log.Fatal("Did not received any data, check your connection")
			} else {
				log.Printf("%x", chunks)
				log.Print("Powered on data received")
			}

		case "r":

			temp := Read()

			WritePort(s, temp)

			log.Print("Reading data...")
			chunks = ReadPort(s)
			if chunks == nil {
				log.Print("Did not received any data, check your connection")
			} else {
				log.Printf("Received data: %x", chunks)
			}

			Stats(chunks)

		case "w":

			fmt.Println("Update with new data(y/n/t): ")
			scan.Scan()
			action = scan.Text()
			temp := cfg

			if action == "y" {
				log.Print("Editing cfg.json...")
				temp, err = File()
				if err != nil {
					log.Print("Error accessing cfg.json: ", err)
				}
				log.Print("cfg.json updated\n")
			} else if action == "t" {
				_ = Edit()
				temp, _ = File()
			}

			WritePort(s, temp)

			log.Print("Awaiting for the return data")
			chunks = ReadPort(s)
			if chunks == nil {
				log.Print("Did not received any data, check your connection")
			} else {
				log.Printf("Received data: %x", chunks)

				err = Save(chunks)
				if err != nil {
					log.Println("Failed to save the input: ", err)
				}

				Stats(chunks)
			}

		case "q":
			os.Exit(1)

		default:
			log.Print("Invalid function code, please select one of the following: w, r, q")
		}

		fmt.Print("\n\nPlease select an action (r = read; w = write; q = quit): ")
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
			log.Print("Corrupted data: ", chunks)
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
		log.Println(err)
	}

	_, err = port.Write(msg)
	if err != nil {
		log.Fatal("Failed to write data: ", err)
	}
}

//Stats outputs all the stats
func Stats(chunks []byte) {
	fmt.Printf("\nSender: %x\nReceiver: %x\nControl Line: %x\nPower: %x\nBrightness: %x\nColorTemp: %x\nColor: %x\nAuto: %x\nSomebody: %x\nNobody: %x\nChained: %x\nTransition: %x\nDelay: %x\nLightType: %x\n", chunks[7:10], chunks[10:13], chunks[15], chunks[16], chunks[17], chunks[18], chunks[19:23], chunks[23], chunks[24], chunks[25], chunks[26], chunks[27], chunks[28], chunks[29])
}
