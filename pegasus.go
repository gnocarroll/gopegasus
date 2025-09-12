package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	filename := ""

	args := flag.Args()
	nargs := len(args)

	if nargs > 0 {
		filename = args[0]
	}

	fmt.Printf("FILE: %s\n", filename)
}
