package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(1) // correct
	run()          // correct
	os.Exit(0)     // want "avoid using os.Exit directly in main function"
}

func run() {
	os.Exit(0)
}
