/*
Description: this will read in the file and the internet user inputted int, 
			 then output the Receive value and Trasmit face value
Date: 5/28/2019
Name: Johnson Zhuang
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main(){

	//read the files and put them into array len = 9 
	file, err := os.Open("test.txt")
 
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
 
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
 	
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
 
	file.Close()
 	var length = len(txtlines)

/* 
	for _, eachline := range txtlines {
		fmt.Println(eachline)
	}
*/
	//find ":" use it as delimeter and get the words in front
	var titles []string 				
	var receive []string
	var transmit []string

	for i:=1; i<length; i++{
		//x := txtlines[i]
		x := strings.Fields(txtlines[i])
		temp := strings.Replace(x[0], ":", "", -1)
		titles = append(titles, strings.Replace(temp, " ", "", -1))
		receive = append(receive, x[1])
		transmit = append(transmit, x[2])

		//fmt.Println(titles[i-1], receive[i-1], transmit[i-1])
	}

	//ask for inputs and output
	fmt.Print("Enter the internet: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
    input := scanner.Text()
	

	if scanner.Err() != nil {
    	fmt.Println("Scanner error")
	}

	check := 0
	for i, n := range titles {
        if input == n {
            fmt.Println("\nReceive: ", receive[i], "\nTransmit face : ", transmit[i])
            check = 1
        }else if i+1 == len(titles) && check==0{
        	fmt.Println("Invalid entry. exiting.")
        }
    }
}