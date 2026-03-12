package ast

// Stmt is a node that performs an action but produces no value.
// Examples: variable declarations, if/else blocks, return statements.
type Stmt interface {
	stmtNode()
}

// Expr is a node that evaluates to a value.
// Examples: literals, binary operations, function calls.
type Expr interface {
	exprNode()
}
