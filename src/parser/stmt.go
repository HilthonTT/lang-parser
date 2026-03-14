package parser

import (
	"lang-parser/src/ast"
	"lang-parser/src/lexer"
)

func parseStmt(p *parser) ast.Stmt {
	stmtFn, exists := stmtTable[p.currentTokenKind()]

	// example: let x = 10 + 5;
	if exists {
		return stmtFn(p)
	}

	// example: 10 + 5;
	expression := parseExpr(p, defaultBP)
	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStmt{
		Expression: expression,
	}
}

func parseVariableDeclarationStmt(p *parser) ast.Stmt {
	isConstant := p.advance().Kind == lexer.CONST
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name").Value

	p.expect(lexer.ASSIGNMENT)
	assignedValue := parseExpr(p, assignment)
	p.expect(lexer.SEMI_COLON)

	return ast.VariableDeclarationStmt{
		VariableName:  varName,
		IsConstant:    isConstant,
		AssignedValue: assignedValue,
	}
}
