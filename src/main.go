package main

import (
	"lang-parser/src/lexer"
	"lang-parser/src/parser"
	"os"

	"github.com/sanity-io/litter"
)

func main() {
	bytes, err := os.ReadFile("./examples/02.lang")
	if err != nil {
		panic(err)
	}

	source := string(bytes)
	tokens := lexer.Tokenize(source)

	ast := parser.Parse(tokens)
	litter.Dump(ast)
}
