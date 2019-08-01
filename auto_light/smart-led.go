/*
*
* Copyright 2018 huayuan-iot
*
* Author: Johnson
* Date: 2019/07/29
* Despcription: smart LED implement
*
 */

package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"clc.hmu/app/public"
	"github.com/tarm/serial"
)

// SmartLEDClientID id
var SmartLEDClientID = "smartledclient"

// SmartLEDClient client
type SmartLEDClient struct {
	ClientID string

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
	return light.ClientID
}

//Sample sample, get values
func (light *SmartLEDClient) Sample(payload string) (string, error) {

	cfg, err := DecodeSmartLEDPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	//determing which control line is wanted
	path, err := strconv.Atoi(cfg.Channel[len(cfg.Channel)-1:])
	if err != nil {
		return "", err
	}

	if path == 0 {
		cfg.ControlLn = 0x00
	} else if path == 1 {
		cfg.ControlLn = 0x01
	} else {
		x := "invalid channel id"
		return "", fmt.Errorf("failed to recognize the control line, errmsg [%v]", x)
	}

	cfg.Channel = cfg.Channel[:len(cfg.Channel)-1]
	read := SendFrmae(&cfg)

	err = WritePort(light.Port, read)
	if err != nil {
		return "", err
	}

	chunks, err := ReadPort(light.Port)
	if err != nil {
		return "", err
	}

	if chunks == nil {
		return "", fmt.Errorf("Did not received any data, check your connection")
	}

	value, err := receive(cfg.Channel, chunks)
	if err != nil {
		return "", err
	}

	//fmt.Printf("%v : %v\n", cfg.Channel, value)
	return value, nil
}

//Command command, set values
func (light *SmartLEDClient) Command(payload string) (string, error) {

	cfg, err := DecodeSmartLEDPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	//adding 0x08 to the value modifying

	//determing which control line is wanted
	path, err := strconv.Atoi(cfg.Channel[len(cfg.Channel)-1:])
	if err != nil {
		return "", err
	}

	if path == 0 {
		cfg.ControlLn = 0x00
	} else if path == 1 {
		cfg.ControlLn = 0x01
	} else {
		x := "invalid channel id"
		return "", fmt.Errorf("failed to recognize the control line, errmsg [%v]", x)
	}

	edit := SendFrmae(&cfg)
	edit.Channel = cfg.Channel[:len(cfg.Channel)-1]

	switch edit.Channel {
	case "setPower":
		edit.LightType = 0x80
		edit.Power = cfg.Power + 0x80

	case "setBrightness":
		edit.LightType = 0x81
		edit.Brightness = cfg.Brightness + 0x80

	case "setColorTemp":
		edit.LightType = 0x82
		edit.ColorTemp = cfg.ColorTemp + 0x80

	case "setColor":
		edit.LightType = 0x82
		cfg.Color[0] += 0x80
		cfg.Color[1] += 0x80
		cfg.Color[2] += 0x80
		cfg.Color[3] += 0x80
		edit.Color = cfg.Color

	case "setAuto":
		edit.LightType = 0x81
		edit.Auto = cfg.Auto + 0x80

	case "setSomebody":
		edit.LightType = 0x81
		edit.Somebody = cfg.Somebody + 0x80

	case "setNobody":
		edit.LightType = 0x81
		edit.Nobody = cfg.Nobody + 0x80

	case "setChained":
		edit.LightType = 0x81
		edit.Chained = cfg.Chained + 0x80

	case "setTransition":
		edit.LightType = 0x81
		edit.Transition = cfg.Transition + 0x80

	case "setDelay":
		edit.LightType = 0x81
		edit.Delay = cfg.Delay + 0x80

	case "setLightType":
		edit.LightType = cfg.LightType + 0x80

	default:
		x := "invalid channel id"
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", x)
	}

	//sending out data to edit
	err = WritePort(light.Port, edit)
	if err != nil {
		return "", err
	}

	//receiving data for feed back
	chunks, err := ReadPort(light.Port)
	if err != nil {
		return "", err
	}

	if chunks == nil {
		return "", fmt.Errorf("Did not received any data, check your connection")
	}

	value, err := receive(edit.Channel[3:], chunks)
	fmt.Println("channel: ", edit.Channel[3:])
	if err != nil {
		return "", err
	}

	return value, nil
}

func receive(option string, chunks []byte) (string, error) {
	if len(chunks) < 30 {
		x := "index out of range"
		return "", fmt.Errorf("[%v]", x)
	}

	var value string
	switch option {
	case "Power":
		value = strconv.Itoa(int(chunks[16]))

	case "Brightness":
		value = strconv.Itoa(int(chunks[17]))

	case "ColorTemp":
		value = strconv.Itoa(int(chunks[18]))

	case "Color":
		value = strconv.Itoa(int(chunks[19])) + strconv.Itoa(int(chunks[20])) + strconv.Itoa(int(chunks[21])) + strconv.Itoa(int(chunks[22]))

	case "Auto":
		value = strconv.Itoa(int(chunks[23]))

	case "Somebody":
		value = strconv.Itoa(int(chunks[24]))

	case "Nobody":
		value = strconv.Itoa(int(chunks[25]))

	case "Chained":
		value = strconv.Itoa(int(chunks[26]))

	case "Transition":
		value = strconv.Itoa(int(chunks[27]))

	case "Delay":
		value = strconv.Itoa(int(chunks[28]))

	case "LightType":
		value = strconv.Itoa(int(chunks[29]))

	default:
		x := "invalid channel id"
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", x)
	}

	return value, nil
}

//WritePort writes data to the port
func WritePort(port *serial.Port, cfg *public.SmartLEDOperationPayload) error {
	//getting the data frame
	head := []byte{0xa5, 0xa5, 0xa5, 0xa5, 0x03}
	version := byte(0x01)
	cfg.Sender = []byte{0x00, 0x00, 0x01}

	frame := DataFrame(head, version, cfg)

	_, err := port.Write(frame)
	if err != nil {
		return err
	}
	//fmt.Printf("out: %x\n", frame)
	return nil
}

//ReadPort reads data from the port
func ReadPort(port *serial.Port) ([]byte, error) {
	var chunks []byte
	buf := make([]byte, 128)

	for len(chunks) < 33 {
		n, err := port.Read(buf)
		if err != nil && err != io.EOF {
			return chunks, err
		}

		if n == 0 {
			x := "no returning data"
			return nil, fmt.Errorf("decode payload failed, errmsg [%v]", x)
		}

		for _, v := range buf[:n] {
			chunks = append(chunks, v)
		}
		//fmt.Println("here: ", chunks, " and ", n)
	}

	if len(chunks) > 33 {
		x := "invalid data frame received"
		return chunks, fmt.Errorf("decode payload failed, errmsg [%v]", x)
	}
	//fmt.Printf("In: %x", chunks)

	return chunks, nil
}

// DecodeSmartLEDBindingPayload decode binding payload
func DecodeSmartLEDBindingPayload(payload string) (public.CommonSerialBindingPayload, error) {
	var p public.CommonSerialBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

//DecodeSmartLEDPayload decodes the payload
func DecodeSmartLEDPayload(payload string) (public.SmartLEDOperationPayload, error) {
	var p public.SmartLEDOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}
	return p, nil
}

//SendFrmae returns a cfg for reading status
func SendFrmae(c *public.SmartLEDOperationPayload) *public.SmartLEDOperationPayload {
	cfg := &public.SmartLEDOperationPayload{
		Port:       c.Port,
		BaudRate:   9600,
		Length:     0x00,
		Sender:     []byte{0x00, 0x00, 0x01},
		Reciever:   c.Reciever,
		Number:     0x00,
		FnCode:     0x1f,
		ControlLn:  c.ControlLn,
		Power:      0x00,
		Brightness: 0x00,
		ColorTemp:  0x00,
		Color:      []byte{0x00, 0x00, 0x00, 0x00},
		Auto:       0x00,
		Somebody:   0x00,
		Nobody:     0x00,
		Chained:    0x00,
		Transition: 0x00,
		Delay:      0x00,
		LightType:  0x00,
		Addition:   nil,
	}

	return cfg
}

//DataFrame creates the data frame
func DataFrame(head []byte, ver byte, cfg *public.SmartLEDOperationPayload) []byte {
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
