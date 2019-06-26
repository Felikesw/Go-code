package main

import (
	"testing"
)

//TestopenFile test to see if openFile works correctly, if it does, it then proceed to test the data retrived
func TestOpenFile(t *testing.T) {
	pair, err := OpenFile("test.json")
	if err != nil {
		t.Fatal(err)
	}

	data, ok := pair["test1"]
	if !ok {
		t.Fatalf("Failed to retrived the correct key, supposed to be \"test1\"")
	}

	if data != "192.168.1.80" {
		t.Fatalf("The data retrived from the key is incorrect, supposed to be 192.168.1.80")
	}

}

//TestWriteFile test to see if the data is written onto the file or not, if yes, it will retrive it and check the contents
func TestWriteFile(t *testing.T) {
	testPair := map[string]string{
		"test1": "192.168.1.80",
		"test2": "123456",
	}
	WriteFile(testPair, "test.json")

	pair, err := OpenFile("test.json")
	if err != nil {
		t.Fatal(err)
	}

	data, ok := pair["test2"]
	if !ok {
		t.Error("Failed to retrived the correct key, supposed to be \"test2\"")
	}

	if data != "123456" {
		t.Error("The data retrived from the key is incorrect, supposed to be 123456")
	}

}
