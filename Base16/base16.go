package main

import (
    "fmt"
)

func main(){

	results := "0B8F"
    for j := 0; j < 12; j += 2 {
			d := int32(0)
			d |= int32(results[j]) << 8
			d |= int32(results[j+1])

			fmt.Printf("s:%d: data:%d\n", j/2, d)
		}
}