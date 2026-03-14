package parser

import (
	"lang-parser/src/ast"
	"lang-parser/src/lexer"
)

type bindingPower int

const (
	defaultBP bindingPower = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

// stmtHandler parses a full statement (e.g. let, return, if).
type stmtHandler func(p *parser) ast.Stmt

// nudHandler is the "null denotation" — handles tokens that appear at the
// start of an expression (literals, prefix operators, grouped expressions).
type nudHandler func(p *parser) ast.Expr

// ledHandler is the "left denotation" — handles tokens that appear in the
// middle of an expression (infix/postfix operators).
type ledHandler func(p *parser, left ast.Expr, bp bindingPower) ast.Expr

var (
	bindingPowers = map[lexer.TokenKind]bindingPower{}
	nudTable      = map[lexer.TokenKind]nudHandler{}
	ledTable      = map[lexer.TokenKind]ledHandler{}
	stmtTable     = map[lexer.TokenKind]stmtHandler{}
)

func led(kind lexer.TokenKind, bp bindingPower, handler ledHandler) {
	bindingPowers[kind] = bp
	ledTable[kind] = handler
}

func nud(kind lexer.TokenKind, handler nudHandler) {
	nudTable[kind] = handler
}

func stmt(kind lexer.TokenKind, handler stmtHandler) {
	bindingPowers[kind] = defaultBP
	stmtTable[kind] = handler
}

func init() {
	led(lexer.ASSIGNMENT, assignment, parseAssignmentExpr)
	led(lexer.PLUS_EQUALS, assignment, parseAssignmentExpr)
	led(lexer.MINUS_EQUALS, assignment, parseAssignmentExpr)
	// TODO: Add *= /= %=

	// Logical
	led(lexer.AND, logical, parseBinaryExpr)
	led(lexer.OR, logical, parseBinaryExpr)
	led(lexer.DOT_DOT, logical, parseBinaryExpr) // example: 10..math.random()

	// Relational
	led(lexer.LESS, relational, parseBinaryExpr)
	led(lexer.LESS_EQUALS, relational, parseBinaryExpr)
	led(lexer.GREATER, relational, parseBinaryExpr)
	led(lexer.GREATER_EQUALS, relational, parseBinaryExpr)
	led(lexer.EQUALS, relational, parseBinaryExpr)
	led(lexer.NOT_EQUALS, relational, parseBinaryExpr)

	// Additive & Multiplicative
	led(lexer.PLUS, additive, parseBinaryExpr)
	led(lexer.DASH, additive, parseBinaryExpr)
	led(lexer.STAR, multiplicative, parseBinaryExpr)
	led(lexer.SLASH, multiplicative, parseBinaryExpr)
	led(lexer.PERCENT, multiplicative, parseBinaryExpr)

	// Literals & identifiers
	nud(lexer.NUMBER, parsePrimaryExpr)
	nud(lexer.STRING, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, parsePrimaryExpr)
	nud(lexer.OPEN_PAREN, parseGroupingExpr)
	nud(lexer.DASH, parsePrefixExpr)

	// Statements
	stmt(lexer.CONST, parseVariableDeclarationStmt)
	stmt(lexer.LET, parseVariableDeclarationStmt)
}
