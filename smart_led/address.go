package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/tarm/serial"
)

//GetAdd gets the address of the module
func main() {

	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Do you know the address that is consist of three bytes(y/n):")

	scanner.Scan()
	answer := scanner.Text()
	var store []byte

	//if "n" then power up the module, and get the address from the power up info
	if answer == "n" {
		//connetcting to the serial port
		c := &serial.Config{Name: "COM4", Baud: 9600, ReadTimeout: time.Second * 5}
		s, err := serial.OpenPort(c)
		if err != nil {
			log.Fatal("Failed to open port: ", err)
		}
		defer s.Close()

		log.Print("\nAwaiting for the device to power on...")

		addre := ReadPort(s)
		if store = addre; addre == nil {
			log.Fatal("Did not received any data, check your connection")
		}

	} else if answer == "y" { //if "y", then prefix the user entried with 0x00

		scanner := bufio.NewScanner(os.Stdin)
		log.Println("Please enter the three-byte address (each byte followed by an enter):")

		store = append(store, 0x00)

		for i := 0; i < 3; i++ {
			scanner.Scan()
			answer := scanner.Text()
			set, err := strconv.ParseInt(answer, 16, 0)
			if err != nil {
				log.Fatal("Error: ", err)
			}
			store = append(store, uint8(set))
		}

	} else {
		log.Fatal("Error: invalid function code")
	}

	//turn the []byte into binary then to an int
	log.Println(store)
	buf := bytes.NewBuffer(store)
	var x int32
	binary.Read(buf, binary.BigEndian, &x)
	log.Println("The address is: ", x)

}

//ReadPort reads data from the port then returns the sender address
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

	return chunks[7:10]
}
