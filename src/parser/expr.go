package parser

import (
	"fmt"
	"lang-parser/src/ast"
	"lang-parser/src/helpers"
	"lang-parser/src/lexer"
	"strconv"
)

func parsePrimaryExpr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		token := p.advance()
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			panic(fmt.Sprintf("Invalid number literal %q: %v", token.Value, err))
		}
		return ast.NumberLiteral{Value: value}

	case lexer.STRING:
		return ast.StringLiteral{Value: p.advance().Value}

	case lexer.IDENTIFIER:
		return ast.SymbolExpr{Value: p.advance().Value}

	default:
		panic(fmt.Sprintf("Expected a primary expression, got %s", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parseBinaryExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	operatorToken := p.advance()
	right := parseExpr(p, bp)

	return ast.BinaryExpr{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}

// parseExpr parses an expression using Pratt (top-down operator precedence) parsing.
//
// The algorithm works in two phases:
//  1. NUD (null denotation): parse the token at the start of an expression —
//     a literal, identifier, or prefix operator like `-`.
//  2. LED (left denotation): as long as the next operator binds tighter than
//     our current binding power, consume it as an infix/postfix operator,
//     folding `left` into a larger expression.
//
// This naturally handles operator precedence without explicit precedence climbing.
func parseExpr(p *parser, bp bindingPower) ast.Expr {
	// Phase 1 — prefix/primary
	nudFn, exists := nudTable[p.currentTokenKind()]
	if !exists {
		panic(syntaxError(p.currentToken(), "expected an expression"))
	}

	left := nudFn(p)

	// Phase 2 — infix: keep consuming operators that bind tighter than bp
	for bindingPowers[p.currentTokenKind()] > bp {
		ledFn, exists := ledTable[p.currentTokenKind()]
		if !exists {
			panic(syntaxError(p.currentToken(), "no infix handler for %s", lexer.TokenKindString(p.currentTokenKind())))
		}
		left = ledFn(p, left, bindingPowers[p.currentTokenKind()])
	}

	return left
}

func parseAssignmentExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	operatorToken := p.advance()
	rhs := parseExpr(p, bp)

	return ast.AssignmentExpr{
		Operator: operatorToken,
		Value:    rhs,
		Assignee: left,
	}
}

func parsePrefixExpr(p *parser) ast.Expr {
	operatorToken := p.advance()
	rhs := parseExpr(p, defaultBP)

	return ast.PrefixExpr{
		Operator:  operatorToken,
		RightExpr: rhs,
	}
}

func parseGroupingExpr(p *parser) ast.Expr {
	p.advance() // Advance pass grouping start
	expr := parseExpr(p, defaultBP)
	p.expect(lexer.CLOSE_PAREN) // advance past close
	return expr
}

func parseStructInstatiationExpr(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	structName := helpers.ExpectType[ast.SymbolExpr](left).Value
	properties := map[string]ast.Expr{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		propertyName := p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)

		expr := parseExpr(p, logical)
		properties[propertyName] = expr

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructInstatiationExpr{
		StructName: structName,
		Properties: properties,
	}
}

func parseArrayLiteral(p *parser) ast.Expr {
	var underlyingType ast.Type
	contents := []ast.Expr{}

	p.expect(lexer.OPEN_BRACKET)
	p.expect(lexer.CLOSE_BRACKET)
	// If you want to have sizes here then you will need to handle that instead.

	underlyingType = parseType(p, defaultBP)

	p.expect(lexer.OPEN_CURLY)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		contents = append(contents, parseExpr(p, logical))

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}
	p.expect(lexer.CLOSE_CURLY)

	return ast.ArrayLiteralExpr{
		Underlying: underlyingType,
		Contents:   contents,
	}
}
