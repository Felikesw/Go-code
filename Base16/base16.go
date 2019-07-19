package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	fmt.Println("Please enter the lower two base16 numbers: ")

	scan := bufio.NewScanner(os.Stdin)

	out, err := Value(scan)

	for err == nil {
		fmt.Println("The equivalent base10 is: ", out)
		out, err = Value(scan)
	}
}

//Value takes in a scanner, scans the command and give back a int and err
func Value(s *bufio.Scanner) (int64, error) {
	s.Scan()
	word := s.Text()

	if word == "q" {
		os.Exit(1)
	}

	lower, err := strconv.ParseInt(word, 16, 64)
	if err != nil {
		return -1, err
	}

	return lower, nil
}
