package main

import (
	"fmt"
	"os"
)

func notMain() {
	fmt.Println("Test for os.Exit")
	os.Exit(1)
}
