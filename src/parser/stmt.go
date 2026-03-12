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
