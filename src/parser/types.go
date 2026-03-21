package parser

import (
	"lang-parser/src/ast"
	"lang-parser/src/lexer"
)

type typeNudHandler func(p *parser) ast.Type
type typeLedHandler func(p *parser, left ast.Type, bp bindingPower) ast.Type

var (
	typeBindingPowers = map[lexer.TokenKind]bindingPower{}
	typeNudTable      = map[lexer.TokenKind]typeNudHandler{}
	typeLedTable      = map[lexer.TokenKind]typeLedHandler{}
)

func typeLed(kind lexer.TokenKind, bp bindingPower, handler typeLedHandler) {
	typeBindingPowers[kind] = bp
	typeLedTable[kind] = handler
}

func typeNud(kind lexer.TokenKind, handler typeNudHandler) {
	typeNudTable[kind] = handler
}

func init() {
	typeNud(lexer.IDENTIFIER, parseSymbolType)
	typeLed(lexer.OPEN_BRACKET, call, parseArrayType)
}

func parseSymbolType(p *parser) ast.Type {
	return ast.SymbolType{
		Name: p.expect(lexer.IDENTIFIER).Value,
	}
}

func parseArrayType(p *parser, left ast.Type, bp bindingPower) ast.Type {
	p.expect(lexer.OPEN_BRACKET)
	p.expect(lexer.CLOSE_BRACKET)

	return ast.ArrayType{
		Underlying: left,
	}
}

func parseType(p *parser, bp bindingPower) ast.Type {
	// Phase 1 — prefix/primary
	nudFn, exists := typeNudTable[p.currentTokenKind()]
	if !exists {
		panic(syntaxError(p.currentToken(), "expected an expression"))
	}

	left := nudFn(p)

	// Phase 2 — infix: keep consuming operators that bind tighter than bp
	for typeBindingPowers[p.currentTokenKind()] > bp {
		ledFn, exists := typeLedTable[p.currentTokenKind()]
		if !exists {
			panic(syntaxError(p.currentToken(), "no infix type handler for %s", lexer.TokenKindString(p.currentTokenKind())))
		}
		left = ledFn(p, left, typeBindingPowers[p.currentTokenKind()])
	}

	return left
}
