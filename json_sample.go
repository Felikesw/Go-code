package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

//Address contains the name and the ip of the address
type Address struct {
	Name string
	IP   string
}

//Addresses contains all the address
type Addresses struct {
	Addrs []Address
}

func ip(w http.ResponseWriter, r *http.Request) {

	//getting all the interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var ips Addresses

	//loop through the interfaces
	for _, iface := range interfaces {

		//getting the address of the interfaces
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}

		//loop through the address(es)
		for _, addr := range addrs {

			//check if the address is a loopback or not
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

				//check if this is a IPv4 or not
				if ipnet.IP.To4() != nil {
					ips.Addrs = append(ips.Addrs, Address{Name: iface.Name, IP: addr.String()})
				}
			}
		}
	}

	//convert it to JSON
	IPJSON, err := json.MarshalIndent(ips, "", "	")
	if err != nil {
		panic(err)
	}

	//write it to http.ResponseWriter
	w.Write(IPJSON)
}

//start the server
func server() {
	http.HandleFunc("/", ip)
	http.ListenAndServe(":8080", nil)
}

//the client
func client() {

	//get a response from the server
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//read the response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	//convert and output the response body
	log.Println(string(body))
}

func main() {
	go server()
	client()
}
