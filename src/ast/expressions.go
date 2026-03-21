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

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expr
}

func (PrefixExpr) exprNode() {}

// a = a + 5;
// a += 5;
// foo.bar += 10;
type AssignmentExpr struct {
	Assignee Expr
	Operator lexer.Token
	Value    Expr
}

func (AssignmentExpr) exprNode() {}

type StructInstatiationExpr struct {
	StructName string
	Properties map[string]Expr
}

func (StructInstatiationExpr) exprNode() {}

type ArrayLiteralExpr struct {
	Underlying Type
	Contents   []Expr
}

func (ArrayLiteralExpr) exprNode() {}
