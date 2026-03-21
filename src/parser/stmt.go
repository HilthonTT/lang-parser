package parser

import (
	"fmt"
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
	var explicitType ast.Type
	var assignedValue ast.Expr

	isConstant := p.advance().Kind == lexer.CONST
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name").Value

	// Explicit type could be present
	if p.currentTokenKind() == lexer.COLON {
		p.advance() // eat the colon
		explicitType = parseType(p, defaultBP)
	}

	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assignedValue = parseExpr(p, assignment)
	} else if explicitType == nil {
		panic("Missing either right-hand side in var declaration or explicit type.")
	}

	p.expect(lexer.SEMI_COLON)

	if isConstant && assignedValue == nil {
		panic("Cannot define constant without providing value")
	}

	return ast.VariableDeclarationStmt{
		ExplicitType:  explicitType,
		VariableName:  varName,
		IsConstant:    isConstant,
		AssignedValue: assignedValue,
	}
}

func parseStructDeclarationStmt(p *parser) ast.Stmt {
	p.expect(lexer.STRUCT) // advance past struct keyword
	structName := p.expect(lexer.IDENTIFIER).Value
	properties := map[string]ast.StructProperty{}
	methods := map[string]ast.StructMethod{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		var isStatic bool
		var propertyName string
		if p.currentTokenKind() == lexer.STATIC {
			isStatic = true
			p.expect(lexer.STATIC)
		}

		// Property
		if p.currentTokenKind() == lexer.IDENTIFIER {
			propertyName = p.expect(lexer.IDENTIFIER).Value
			p.expectError(lexer.COLON, "Expected to find colon following property name inside struct declaration")
			structType := parseType(p, defaultBP)
			p.expect(lexer.SEMI_COLON)

			_, exists := properties[propertyName]
			if exists {
				panic(fmt.Sprintf("Property %s has already been defined in struct declaration", propertyName))
			}

			properties[propertyName] = ast.StructProperty{
				IsStatic: isStatic,
				Type:     structType,
			}

			continue
		}

		panic("Cannot currently handle methods inside struct declaration")
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructDeclarationStmt{
		Properties: properties,
		Methods:    methods,
		StructName: structName,
	}
}
