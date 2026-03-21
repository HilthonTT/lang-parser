package main

import (
	"fmt"
	"lang-parser/src/lexer"
	"lang-parser/src/parser"
	"os"

	"github.com/sanity-io/litter"
)

func main() {
	bytes, err := os.ReadFile("./examples/05.lang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	source := string(bytes)
	tokens := lexer.Tokenize(source)

	tree, err := parser.Parse(tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	litter.Dump(tree)
}
