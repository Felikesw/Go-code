package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"github.com/tarm/serial"
)

//Cfg is the configeration of the data frame
type Cfg struct {
	Port       string `json:"port"`
	BaudRate   int    `json:"baud_rate"`
	Length     byte   `json:"length"`
	Version    byte   `json:"version"`
	Sender     []byte `json:"sender"`
	Reciever   []byte `json:"reciever"`
	Number     byte   `json:"number"`
	FnCode     byte   `json:"fn_code"`
	ControlLn  byte   `json:"control_ln"`
	Power      byte   `json:"power"`
	Brightness byte   `json:"brightness"`
	ColorTemp  byte   `json:"color_temp"`
	Color      []byte `json:"color"`
	Auto       byte   `json:"auto"`
	Somebody   byte   `json:"somebody"`
	Nobody     byte   `json:"nobody"`
	Chained    byte   `json:"chained"`
	Transition byte   `json:"transition"`
	Delay      byte   `json:"delay"`
	LightType  byte   `json:"light_type"`
	Addition   []byte `json:"addition"`
	Check      []byte `json:"check"`
}

//DataFrame creates the data frame
func DataFrame(head []byte, ver byte, cfg *Cfg) ([]byte, error) {
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

	return data, nil

}

//example command
func example(ver byte) *Cfg {
	cfg := &Cfg{
		BaudRate:   9600,
		Version:    ver,
		Length:     0x1c,
		Sender:     []byte{0x06, 0x3f, 0x5e},
		Reciever:   []byte{0x00, 0x00, 0x01},
		Number:     0x07,
		FnCode:     0x40,
		ControlLn:  0x01,
		Power:      0x01,
		Brightness: 0x03,
		ColorTemp:  0x03,
		Color:      []byte{0x0f, 0x00, 0x00, 0x05},
		Auto:       0x01,
		Somebody:   0x0f,
		Nobody:     0x03,
		Chained:    0x0f,
		Transition: 0x04,
		Delay:      0x01,
		LightType:  0x01,
		Addition:   []byte{0x33},
		Check:      []byte{0x8e, 0x06},
	}

	return cfg
}

//File creates a cfg
func File() (*Cfg, error) {
	// test := example(version)
	// writer, _ := os.OpenFile("cfg.json", os.O_RDWR|os.O_TRUNC, 0644)
	// defer writer.Close()
	// je := json.NewEncoder(writer)
	// je.Encode(test)
	// fmt.Println(test)

	data, err := ioutil.ReadFile("./cfg.json")
	if err != nil {
		return &Cfg{}, err
	}
	cfg := &Cfg{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return &Cfg{}, err
	}

	return cfg, nil
}

//Light is the actual main
func Light(cfg *Cfg) ([]byte, error) {

	head := []byte{0xa5, 0xa5, 0xa5, 0xa5, 0x03}
	//version := byte(0x01)

	switch cfg.FnCode {
	case 0x1f: //editing parameters. Format for the editing parameter value(s): 0x80 + parameter value
		fmt.Println("Editing parameters")
		frame, _ := DataFrame(head, cfg.Version, cfg)
		fmt.Printf("Data frame sent: %x\n", frame)
		return frame, nil
	default:
		return nil, fmt.Errorf("Errors: invalid function code")
	}
}
