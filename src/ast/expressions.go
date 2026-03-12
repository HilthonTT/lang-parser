package ast

import "lang-parser/src/lexer"

// ── Literals ─────────────────────────────────────────────────────────────────

type NumberLiteral struct {
	Value float64
}

func (NumberLiteral) exprNode() {}

type StringLiteral struct {
	Value string
}

func (StringLiteral) exprNode() {}

// SymbolExpr represents an identifier (variable or function name).
type SymbolExpr struct {
	Value string
}

func (SymbolExpr) exprNode() {}

// ── Compound expressions ──────────────────────────────────────────────────────

// BinaryExpr represents an infix operation between two expressions.
// Example: 10 + 5 * 2
type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (BinaryExpr) exprNode() {}
