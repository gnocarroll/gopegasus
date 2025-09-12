package main

import (
	"flag"
	"pegasus/scanner"
)

func main() {
	flag.Parse()

	filename := ""

	args := flag.Args()
	nargs := len(args)

	if nargs > 0 {
		filename = args[0]
	}

	scan := scanner.NewScanner()

	scan.TokenizeFile(filename)
}
