/*
Description: This stores all IPv4 IP addresses which are not loopback and convert them to JSON, then procecced to write them onto the page
name: Johnson Zhuang
Date: 6/9/2019
*/

package main

import (
	"encoding/json"
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

	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var ips Addresses

	for _, iface := range interfaces {

		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}

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
	b, err := json.Marshal(ips)
	if err != nil {
		panic(err)
	}

	//write it to http.ResponseWriter
	w.Write(b)
}

func main() {

	http.HandleFunc("/", ip)

	http.ListenAndServe(":8080", nil)
}
