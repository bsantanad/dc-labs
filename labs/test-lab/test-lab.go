package main

import (
	"fmt"
	"os"
)

func main() {

	var name string
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Printf("Error!\n")
		return
	}

	for i, word := range args {
		name += word
		if i != (len(args) - 1) {
			name += " "
		}
	}
	fmt.Printf("Hello %s, Welcome to the jungle\n", name)
}
