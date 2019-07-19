/*
This gets the contents from local host and unmarshal it
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//ReadContent reads the content at the url and converts it from json to Addresses
func ReadContent() {
	response, err := http.Get("http://localhost:8080/")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		var iface Addresses

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		if err := json.Unmarshal([]byte(string(contents)), &iface); err != nil {
			fmt.Println("ugh: ", err)
		}
		fmt.Println(iface)
	}
}
