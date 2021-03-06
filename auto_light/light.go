package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	Frequency  byte   `json:"frequency"`
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
func DataFrame(head []byte, ver byte, cfg *Cfg) []byte {
	data := head
	var temp0 int
	var temp1 int
	if cfg.FnCode == 0x40 {
		data = append(data, 0x1d)
		temp0 = 0x1d
		temp1 = 0x1d
	} else {
		data = append(data, 0x1c)
		temp0 = 0x1c
		temp1 = 0x1c
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

	data = append(data, cfg.Frequency)
	temp0 += int(cfg.Frequency)
	temp1 ^= int(cfg.Frequency)

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

//Save saves the return
func Save(c []byte) error {

	cfg := &Cfg{
		Port:       "COM6",
		BaudRate:   9600,
		Version:    0x01,
		Length:     c[5],
		Sender:     []byte{0x00, 0x00, 0x01},
		Reciever:   []byte{0x06, 0x3F, 0x25},
		Number:     0x01,
		Frequency:  c[14],
		FnCode:     0x1f,
		ControlLn:  c[16],
		Power:      c[17],
		Brightness: c[18],
		ColorTemp:  c[19],
		Color:      c[20:24],
		Auto:       c[24],
		Somebody:   c[25],
		Nobody:     c[26],
		Chained:    c[27],
		Transition: c[28],
		Delay:      c[29],
		LightType:  c[30],
		Addition:   nil,
		Check:      nil,
	}

	writer, err := os.OpenFile("cfg.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()

	je := json.NewEncoder(writer)
	je.SetIndent("", "	")
	err = je.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

//example command
func example() *Cfg {
	cfg := &Cfg{
		Port:       "COM6",
		BaudRate:   9600,
		Version:    0x01,
		Length:     0x1c,
		Sender:     []byte{0x00, 0x00, 0x01},
		Reciever:   []byte{0x06, 0x3F, 0x25},
		Number:     0x01,
		Frequency:  0x00,
		FnCode:     0x1f,
		ControlLn:  0x01,
		Power:      0x80,
		Brightness: 0x81,
		ColorTemp:  0x00,
		Color:      []byte{0xff, 0xff, 0xff, 0xff},
		Auto:       0x01,
		Somebody:   0x0f,
		Nobody:     0x81,
		Chained:    0x0f,
		Transition: 0x03,
		Delay:      0x03,
		LightType:  0x00,
		Addition:   nil,
		Check:      nil,
	}

	return cfg
}

//Read returns a cfg for reading status
func Read() *Cfg {
	cfg := &Cfg{
		Port:       "COM6",
		BaudRate:   9600,
		Version:    0x01,
		Length:     0x1c,
		Sender:     []byte{0x00, 0x00, 0x01},
		Reciever:   []byte{0x06, 0x3f, 0x25},
		Number:     0x00,
		Frequency:  0x00,
		FnCode:     0x1f,
		ControlLn:  0x01,
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
		Check:      nil,
	}

	return cfg
}

//Edit allows users to edit the parameters
func Edit() error {

	log.Print("Editing cfg.json...")
	test := example()

	writer, err := os.OpenFile("cfg.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()

	je := json.NewEncoder(writer)
	je.SetIndent("", "	")
	err = je.Encode(test)
	if err != nil {
		return err
	}

	log.Print("cfg.json updated\n")
	return nil
}

//File creates a cfg
func File() (*Cfg, error) {

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

//Light returns a frame
func Light(cfg *Cfg) ([]byte, error) {

	head := []byte{0xa5, 0xa5, 0xa5, 0xa5, 0x03}
	//version := byte(0x01)

	switch cfg.FnCode {
	case 0x1f: //editing parameters. Format for the editing parameter value(s): 0x80 + parameter value
		frame := DataFrame(head, cfg.Version, cfg)
		log.Printf("Sent data frame: %x\n", frame)
		return frame, nil
	default:
		return nil, fmt.Errorf("Errors: invalid function code")
	}
}
