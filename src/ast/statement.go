package ast

// ── Statements ────────────────────────────────────────────────────────────────

// BlockStmt is a sequence of statements wrapped in braces.
// Example: { let x = 1; return x }
type BlockStmt struct {
	Body []Stmt
}

func (BlockStmt) stmtNode() {}

// ExpressionStmt wraps an expression where a statement is expected.
// Example: a function call used as a standalone line — foo()
type ExpressionStmt struct {
	Expression Expr
}

func (ExpressionStmt) stmtNode() {}

type VariableDeclarationStmt struct {
	VariableName  string
	IsConstant    bool
	AssignedValue Expr
	ExplicitType  Type
}

func (VariableDeclarationStmt) stmtNode() {}

type StructProperty struct {
	IsStatic bool // determine whether the property is static
	Type     Type
}

type StructMethod struct {
	IsStatic bool
	// Type     FnType
}

type StructDeclarationStmt struct {
	StructName string
	// IsPublic   bool
	Properties map[string]StructProperty
	Methods    map[string]StructMethod
}

func (StructDeclarationStmt) stmtNode() {}
