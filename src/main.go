package main

import (
	"lang-parser/src/lexer"
	"os"
)

func main() {
	bytes, err := os.ReadFile("./examples/01.lang")
	if err != nil {
		panic(err)
	}

	source := string(bytes)
	tokens := lexer.Tokenize(source)

	for _, token := range tokens {
		token.Debug()
	}
}
