package main

import "fmt"

func main() {
	res, err := LoadXML("test")

	fmt.Printf("%q\n%q\n", res, err)
}
