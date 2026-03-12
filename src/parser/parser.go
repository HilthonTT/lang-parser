package parser

import (
	"fmt"
	"lang-parser/src/ast"
	"lang-parser/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokens []lexer.Token) *parser {
	return &parser{
		tokens: tokens,
		pos:    0,
	}
}

func Parse(tokens []lexer.Token) (result ast.BlockStmt, err error) {
	// Catch any SyntaxError panics and return them as normal errors.
	defer func() {
		if r := recover(); r != nil {
			if syntaxErr, ok := r.(*SyntaxError); ok {
				err = syntaxErr
			} else {
				// Re-panic for unexpected bugs — don't swallow real crashes.
				panic(r)
			}
		}
	}()

	p := createParser(tokens)
	body := make([]ast.Stmt, 0)

	for p.hasTokens() {
		body = append(body, parseStmt(p))
	}

	return ast.BlockStmt{Body: body}, nil
}

// HELPER METHODS
func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	p.pos++
	return tk
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) expectError(expectedKind lexer.TokenKind, msg string) lexer.Token {
	token := p.currentToken()
	if token.Kind != expectedKind {
		if msg == "" {
			msg = fmt.Sprintf("expected %s", lexer.TokenKindString(expectedKind))
		}
		panic(syntaxError(token, msg))
	}
	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, "")
}
