package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Test for os.Exit")
	os.Exit(1) // want `expression returns direct call of os.Exit`
}
