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

		assertEquals(t, stmt.TokenLiteral(), "let")
		letStmt := assertStatementType[*ast.LetStatement](t, stmt)
		assertEquals(t, letStmt.Name.Value, tt.expectedIdentifier)
		assertEquals(t, letStmt.Name.TokenLiteral(), tt.expectedIdentifier)
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
		assertEquals(t, returnStmt.TokenLiteral(), "return")
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
	assertLiteralExpression(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p)
	assertStatementsLen(t, program.Statements, 1)
	stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
	assertLiteralExpression(t, stmt.Expression, 5)
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"-15", "-", 15},
		{"!5", "!", 5},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)
		assertStatementsLen(t, program.Statements, 1)
		stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
		exp := assertExpressionType[*ast.PrefixExpression](t, stmt.Expression)
		assertEquals(t, exp.Operator, tt.operator)
		assertLiteralExpression(t, exp.Right, tt.value)
	}
}

func TestInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)
		assertStatementsLen(t, program.Statements, 1)
		stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
		assertInfixExpresssion(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)
		assertStatementsLen(t, program.Statements, 1)
		stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
		assertLiteralExpression(t, stmt.Expression, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p)
	assertStatementsLen(t, program.Statements, 1)

	stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
	exp := assertExpressionType[*ast.IfExpression](t, stmt.Expression)
	assertInfixExpresssion(t, exp.Condition, "x", "<", "y")

	assertStatementsLen(t, exp.Consequence.Statements, 1)
	consequence := assertStatementType[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
	assertIdentifier(t, consequence.Expression, "x")

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p)
	assertStatementsLen(t, program.Statements, 1)

	stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
	exp := assertExpressionType[*ast.IfExpression](t, stmt.Expression)
	assertInfixExpresssion(t, exp.Condition, "x", "<", "y")

	assertStatementsLen(t, exp.Consequence.Statements, 1)
	consequence := assertStatementType[*ast.ExpressionStatement](t, exp.Consequence.Statements[0])
	assertIdentifier(t, consequence.Expression, "x")

	assertStatementsLen(t, exp.Alternative.Statements, 1)
	alternative := assertStatementType[*ast.ExpressionStatement](t, exp.Alternative.Statements[0])
	assertIdentifier(t, alternative.Expression, "y")
}

// helpers

func assertParseErrors(t testing.TB, p *Parser) {
	t.Helper()

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

func assertEquals(t testing.TB, got any, want any) {
	t.Helper()

	if got != want {
		t.Errorf("not equal. got=%v want=%v", got, want)
	}
}

func assertIntegerLiteral(t testing.TB, exp ast.Expression, val int64) {
	t.Helper()

	integ := assertExpressionType[*ast.IntegerLiteral](t, exp)
	assertEquals(t, integ.Value, val)
	assertEquals(t, integ.TokenLiteral(), fmt.Sprintf("%d", val))
}

func assertIdentifier(t testing.TB, exp ast.Expression, value string) {
	t.Helper()

	ident := assertExpressionType[*ast.Identifier](t, exp)
	assertEquals(t, ident.Value, value)
	assertEquals(t, ident.TokenLiteral(), value)
}

func assertLiteralExpression(t testing.TB, exp ast.Expression, expected any) {
	t.Helper()

	switch v := expected.(type) {
	case int:
		assertIntegerLiteral(t, exp, int64(v))
	case int64:
		assertIntegerLiteral(t, exp, v)
	case string:
		assertIdentifier(t, exp, v)
	case bool:
		assertBooleanLiteral(t, exp, v)
	default:
		t.Errorf("type of exp not handled. got %T", exp)
	}
}

func assertInfixExpresssion(t testing.TB, exp ast.Expression, left any, operator string, right any) {
	t.Helper()

	opExp := assertExpressionType[*ast.InfixExpression](t, exp)
	assertLiteralExpression(t, opExp.Left, left)
	assertEquals(t, opExp.Operator, operator)
	assertLiteralExpression(t, opExp.Right, right)
}

func assertBooleanLiteral(t testing.TB, exp ast.Expression, expected bool) {
	t.Helper()

	boolExp := assertExpressionType[*ast.Boolean](t, exp)
	assertEquals(t, boolExp.Value, expected)
	assertEquals(t, boolExp.TokenLiteral(), fmt.Sprintf("%v", expected))
}
