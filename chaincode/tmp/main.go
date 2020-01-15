package main

import "fmt"

func main() {
	subsCodes := []string{"aaaa", "vvvvv", "dddd", "eeeee", "gfgggg"}
	fmt.Println(subsCodes[0])
	fmt.Println(len(subsCodes))
	fmt.Println(subsCodes[1:len(subsCodes)-1])
	fmt.Println(subsCodes[len(subsCodes) -1])
}