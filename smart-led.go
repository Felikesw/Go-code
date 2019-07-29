/*
*
* Copyright 2018 huayuan-iot
*
* Author: Johnson
* Date: 2019/07/29
* Despcription: smart LED implement
*
 */

package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"clc.hmu/app/public"
	"github.com/tarm/serial"
)

func main() {}

// SmartLEDClientID id
var SmartLEDClientID = "smartledclient"

// SmartLEDClient client
type SmartLEDClient struct {
	ClientID  string
	ControlLn byte

	Port *serial.Port
}

// NewSmartLEDClient new client
func NewSmartLEDClient(port string, baudrate, timeout int) (*SmartLEDClient, error) {
	cfg := &serial.Config{Name: port, Baud: baudrate, ReadTimeout: time.Millisecond * time.Duration(timeout)}
	fmt.Println("cfg:", cfg)
	sp, err := serial.OpenPort(cfg)
	if err != nil {
		return &SmartLEDClient{}, err
	}

	var client SmartLEDClient
	client.ClientID = SmartLEDClientID
	client.Port = sp

	return &client, nil
}

//ID specified client's ID, use for searching
func (light *SmartLEDClient) ID() string {
	return light.ClientID + string(int(light.ControlLn))
}

//Sample sample, get values
func (light *SmartLEDClient) Sample(payload string) (string, error) {

	p, err := DecodeSmartLEDPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	var cfg public.SmartLEDPayload

	//connetcting to the serial port
	client, err := NewSmartLEDClient(cfg.Port, cfg.BaudRate, cfg.TimeOut)
	if err != nil {
		return "", err
	}

	s := client.Port
	err = WritePort(s, &cfg)
	if err != nil {
		return "", err
	}

	chunks, err := ReadPort(s)
	if err != nil {
		return "", err
	}

	if chunks == nil {
		return "", fmt.Errorf("Did not received any data, check your connection")
	}

	var value string
	switch p.Channel {
	case "Power":
		value = string(chunks[16])

	case "Brightness":
		value = string(chunks[17])

	case "ColorTemp":
		value = string(chunks[18])

	case "Color":
		value = string(chunks[19:23])

	case "Auto":
		value = string(chunks[23])

	case "Somebody":
		value = string(chunks[24])

	case "Nobody":
		value = string(chunks[25])

	case "Chained":
		value = string(chunks[26])

	case "Transition":
		value = string(chunks[27])

	case "Delay":
		value = string(chunks[28])

	case "LightType":
		value = string(chunks[29])

	default:
		return "invalite channel id", nil
	}

	return value, nil
}

//Command command, set values
func (light *SmartLEDClient) Command(payload string) (string, error) {

	p, err := DecodeSmartLEDPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	var cfg public.SmartLEDPayload

	//+0x80
	//connetcting to the serial port
	client, err := NewSmartLEDClient(cfg.Port, cfg.BaudRate, cfg.TimeOut)
	if err != nil {
		return "", err
	}

	s := client.Port

	err = WritePort(s, &cfg)
	if err != nil {
		return "", err
	}

	chunks, err := ReadPort(s)
	if err != nil {
		return "", err
	}

	if chunks == nil {
		return "", fmt.Errorf("Did not received any data, check your connection")
	}

	var value string
	switch p.Channel {
	case "Power":
		value = string(chunks[16])

	case "Brightness":
		value = string(chunks[17])

	case "ColorTemp":
		value = string(chunks[18])

	case "Color":
		value = string(chunks[19:23])

	case "Auto":
		value = string(chunks[23])

	case "Somebody":
		value = string(chunks[24])

	case "Nobody":
		value = string(chunks[25])

	case "Chained":
		value = string(chunks[26])

	case "Transition":
		value = string(chunks[27])

	case "Delay":
		value = string(chunks[28])

	case "LightType":
		value = string(chunks[29])

	default:
		return "invalite channel id", nil
	}

	return value, nil
}

//WritePort writes data to the port
func WritePort(port *serial.Port, cfg *public.SmartLEDPayload) error {
	//getting the data frame
	head := []byte{0xa5, 0xa5, 0xa5, 0xa5, 0x03}
	version := byte(0x01)

	frame := DataFrame(head, version, cfg)

	_, err := port.Write(frame)
	if err != nil {
		return err
	}
	return nil
}

//ReadPort reads data from the port
func ReadPort(port *serial.Port) ([]byte, error) {
	var chunks []byte
	buf := make([]byte, 128)
	n := 1

	for n != 0 {
		n, err := port.Read(buf)
		if err != nil && err != io.EOF {
			return chunks, err
		}

		for _, v := range buf[:n] {
			chunks = append(chunks, v)
		}
	}

	return chunks, nil
}

//DecodeSmartLEDPayload decodes the payload
func DecodeSmartLEDPayload(payload string) (public.SmartLEDPayload, error) {
	var p public.SmartLEDPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}
	return p, nil
}

//DataFrame creates the data frame
func DataFrame(head []byte, ver byte, cfg *public.SmartLEDPayload) []byte {
	data := head
	var temp0 int
	var temp1 int
	if cfg.FnCode == 0x40 {
		data = append(data, 0x1c)
		temp0 = 0x1c
		temp1 = 0x1c
	} else {
		data = append(data, 0x1b)
		temp0 = 0x1b
		temp1 = 0x1b
	}

	data = append(data, ver)
	temp0 += int(ver)
	temp1 ^= int(ver)

	for _, v := range cfg.Sender {
		data = append(data, v)
		temp0 += int(v)
		temp1 ^= int(v)
	}

	for _, v := range cfg.Reciever {
		data = append(data, v)
		temp0 += int(v)
		temp1 ^= int(v)
	}

	data = append(data, cfg.Number)
	temp0 += int(cfg.Number)
	temp1 ^= int(cfg.Number)

	data = append(data, cfg.FnCode)
	temp0 += int(cfg.FnCode)
	temp1 ^= int(cfg.FnCode)

	data = append(data, cfg.ControlLn)
	temp0 += int(cfg.ControlLn)
	temp1 ^= int(cfg.ControlLn)

	data = append(data, cfg.Power)
	temp0 += int(cfg.Power)
	temp1 ^= int(cfg.Power)

	data = append(data, cfg.Brightness)
	temp0 += int(cfg.Brightness)
	temp1 ^= int(cfg.Brightness)

	data = append(data, cfg.ColorTemp)
	temp0 += int(cfg.ColorTemp)
	temp1 ^= int(cfg.ColorTemp)

	for _, v := range cfg.Color {
		data = append(data, v)
		temp0 += int(v)
		temp1 ^= int(v)
	}

	data = append(data, cfg.Auto)
	temp0 += int(cfg.Auto)
	temp1 ^= int(cfg.Auto)

	data = append(data, cfg.Somebody)
	temp0 += int(cfg.Somebody)
	temp1 ^= int(cfg.Somebody)

	data = append(data, cfg.Nobody)
	temp0 += int(cfg.Nobody)
	temp1 ^= int(cfg.Nobody)

	data = append(data, cfg.Chained)
	temp0 += int(cfg.Chained)
	temp1 ^= int(cfg.Chained)

	data = append(data, cfg.Transition)
	temp0 += int(cfg.Transition)
	temp1 ^= int(cfg.Transition)

	data = append(data, cfg.Delay)
	temp0 += int(cfg.Delay)
	temp1 ^= int(cfg.Delay)

	data = append(data, cfg.LightType)
	temp0 += int(cfg.LightType)
	temp1 ^= int(cfg.LightType)

	if cfg.FnCode == 0x40 {
		for _, v := range cfg.Addition {
			data = append(data, v)
			temp0 += int(v)
			temp1 ^= int(v)
		}
	}

	var check0 = make([]byte, 8)
	binary.BigEndian.PutUint64(check0, uint64(int32(temp0)))
	data = append(data, check0[7])

	var check1 = make([]byte, 8)
	binary.BigEndian.PutUint64(check1, uint64(int64(temp1)))

	data = append(data, check1[7])

	return data

}
