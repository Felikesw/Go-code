package main

import (
	"fmt"
)

func main() {

	x := uint16(75)

	fmt.Printf("%T\n", x)
	fmt.Printf("%T", fmt.Sprintf("%x", x))
	fmt.Printf("%v", x)

}
