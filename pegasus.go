package main

import (
	"fmt"
	"pegasus/scanner"
)

func main() {
	s := "5 5 5"

	scan := scanner.NewScanner()

	scan.Tokenize(s)

	i := 0

	for {
		tok := scan.Advance()

		if tok.TType == scanner.TOK_EOF {
			break
		}

		tokStr := ""

		if int(tok.TType) < len(scanner.TokStrings) &&
			scanner.TokStrings[tok.TType] != "" {

			tokStr = scanner.TokStrings[tok.TType]
		} else {
			tokStr = tok.Text
		}

		fmt.Printf("%d: %s\n", i, tokStr)

		i++
	}
}
