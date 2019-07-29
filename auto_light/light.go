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

//Save saves the return
func Save(c []byte) error {

	cfg := &Cfg{
		Port:       "COM4",
		BaudRate:   9600,
		Version:    0x01,
		Length:     c[5],
		Sender:     c[7:10],
		Reciever:   c[10:13],
		Number:     c[13],
		FnCode:     c[14],
		ControlLn:  c[15],
		Power:      c[16],
		Brightness: c[17],
		ColorTemp:  c[18],
		Color:      c[19:23],
		Auto:       c[23],
		Somebody:   c[24],
		Nobody:     c[25],
		Chained:    c[26],
		Transition: c[27],
		Delay:      c[28],
		LightType:  c[29],
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
		Port:       "COM4",
		BaudRate:   9600,
		Version:    0x01,
		Length:     0x1b,
		Sender:     []byte{0x00, 0x00, 0x01},
		Reciever:   []byte{0x06, 0x3F, 0x58},
		Number:     0x01,
		FnCode:     0x1f,
		ControlLn:  0x00,
		Power:      0x01,
		Brightness: 0x89,
		ColorTemp:  0x00,
		Color:      []byte{0xff, 0xff, 0xff, 0xff},
		Auto:       0x01,
		Somebody:   0x0f,
		Nobody:     0x81,
		Chained:    0x0f,
		Transition: 0x03,
		Delay:      0x03,
		LightType:  0x01,
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
		log.Print("\nEditing parameters...")
		frame := DataFrame(head, cfg.Version, cfg)
		log.Printf("Sent data frame: %x\n", frame)
		return frame, nil
	default:
		return nil, fmt.Errorf("Errors: invalid function code")
	}
}
