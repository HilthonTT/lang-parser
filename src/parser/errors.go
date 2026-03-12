package parser

import (
	"fmt"
	"lang-parser/src/lexer"
)

// SyntaxError describes a parse failure at a known source location.
type SyntaxError struct {
	Message string
	Token   lexer.Token
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf(
		"syntax error at line %d, col %d: %s (got %q)",
		e.Token.Line,
		e.Token.Col,
		e.Message,
		e.Token.Value,
	)
}

func syntaxError(tok lexer.Token, format string, args ...any) *SyntaxError {
	return &SyntaxError{
		Message: fmt.Sprintf(format, args...),
		Token:   tok,
	}
}
