package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 242424;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	assertParseErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	assertStatementsLen(t, program.Statements, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		assertValue(t, stmt.TokenLiteral(), "let")
		letStmt := assertStatementType[*ast.LetStatement](t, stmt)
		assertValue(t, letStmt.Name.Value, tt.expectedIdentifier)
		assertValue(t, letStmt.Name.TokenLiteral(), tt.expectedIdentifier)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 224411;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	assertParseErrors(t, p)
	assertStatementsLen(t, program.Statements, 3)

	for _, stmt := range program.Statements {
		returnStmt := assertStatementType[*ast.ReturnStatement](t, stmt)
		assertValue(t, returnStmt.TokenLiteral(), "return")
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p)
	assertStatementsLen(t, program.Statements, 1)
	stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
	ident := assertExpressionType[*ast.Identifier](t, stmt.Expression)
	assertValue(t, ident.Value, "foobar")
	assertValue(t, ident.TokenLiteral(), "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p)
	assertStatementsLen(t, program.Statements, 1)
	stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
	assertIntegerLiteral(t, stmt.Expression, 5)
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"-15", "-", 15},
		{"!5", "!", 5},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)
		assertStatementsLen(t, program.Statements, 1)
		stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
		exp := assertExpressionType[*ast.PrefixExpression](t, stmt.Expression)
		assertValue(t, exp.Operator, tt.operator)
		assertIntegerLiteral(t, exp.Right, tt.integerValue)
	}
}

// helpers

func assertParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parse error: %q", msg)
	}
	t.FailNow()
}

func assertStatementsLen(t testing.TB, statements []ast.Statement, want int) {
	t.Helper()

	got := len(statements)
	if got != want {
		t.Fatalf("wrong number of statements. got=%d want=%d", got, want)
	}
}

func assertStatementType[T ast.Statement](t testing.TB, got ast.Statement) T {
	t.Helper()

	stmt, ok := got.(T)
	if !ok {
		t.Fatalf("wrong statement type. got=%T want %T", *new(T), got)
	}

	return stmt
}

func assertExpressionType[T ast.Expression](t testing.TB, got ast.Expression) T {
	t.Helper()

	exp, ok := got.(T)

	if !ok {
		t.Fatalf("wrong expression type. got=%T want %T", *new(T), got)
	}

	return exp
}

func assertValue[T string | int64](t testing.TB, got T, want T) {
	t.Helper()

	if got != want {
		t.Errorf("wrong value. got=%v want=%v", got, want)
	}
}

func assertIntegerLiteral(t testing.TB, exp ast.Expression, val int64) {
	t.Helper()

	integ := assertExpressionType[*ast.IntegerLiteral](t, exp)
	assertValue(t, integ.Value, val)
	assertValue(t, integ.TokenLiteral(), fmt.Sprintf("%d", val))
}
