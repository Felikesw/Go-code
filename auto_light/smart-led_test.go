/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author= lynn
 * Date= 2018/10/30
 * Despcription= test file
 *
 */

package portmanager_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"clc.hmu/app/portmanager/src/protocol"
	"clc.hmu/app/public"
)

func TestNewSmartLEDClient(t *testing.T) {
	port := "COM4"
	baudrate := 9600
	timeout := 5000

	client, err := protocol.NewSmartLEDClient(port, baudrate, timeout)
	if err != nil {
		t.Error(err)
	}

	t.Log(client)
}

func TestID(t *testing.T) {

	port := "COM4"
	baudrate := 9600
	timeout := 5000

	client, err := protocol.NewSmartLEDClient(port, baudrate, timeout)
	if err != nil {
		t.Error(err)
	}

	p := client.ID()
	t.Log("the client is= ", p)
}

func TestSample(t *testing.T) {
	port := "COM4"
	baudrate := 9600
	timeout := 3000

	client, err := protocol.NewSmartLEDClient(port, baudrate, timeout)
	if err != nil {
		t.Error(err)
	}
	var p public.SmartLEDOperationPayload
	p.Channel = "Delay0"
	p.Port = port
	p.BaudRate = baudrate
	p.Sender = []byte{0x00, 0x00, 0x01}
	p.Reciever = []byte{0x06, 0x3f, 0x58}

	d, _ := json.Marshal(p)

	v, err := client.Sample(string(d))
	if err != nil {
		t.Error(err)
	}

	fmt.Println("the value is: ", v)
}

func TestCommand(t *testing.T) {
	port := "COM4"
	baudrate := 9600
	timeout := 3000

	client, err := protocol.NewSmartLEDClient(port, baudrate, timeout)
	if err != nil {
		t.Error(err)
	}
	var p public.SmartLEDOperationPayload
	p.Channel = "setBrightness0"
	p.Port = port
	p.BaudRate = baudrate
	p.Sender = []byte{0x00, 0x00, 0x01}
	p.Reciever = []byte{0x06, 0x3f, 0x58}
	p.Brightness = 0x01

	d, _ := json.Marshal(p)

	v, err := client.Command(string(d))
	if err != nil {
		t.Error(err)
	}

	fmt.Println("the value is: ", v)
}
