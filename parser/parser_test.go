package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)
		assertStatementsLen(t, program.Statements, 1)
		stmt := program.Statements[0]
		assertLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedValue)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"return 5;", 5},
		{"return foo;", "foo"},
		{"return true", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)
		assertStatementsLen(t, program.Statements, 1)

		returnStmt := assertStatementType[*ast.ReturnStatement](t, program.Statements[0])
		assertEquals(t, returnStmt.TokenLiteral(), "return")

		val := returnStmt.ReturnValue
		assertLiteralExpression(t, val, tt.expectedValue)
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
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
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

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p)
	assertStatementsLen(t, program.Statements, 1)
	stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])

	function := assertExpressionType[*ast.FunctionLiteral](t, stmt.Expression)
	assertParametersLen(t, function.Parameters, 2)

	assertLiteralExpression(t, function.Parameters[0], "x")
	assertLiteralExpression(t, function.Parameters[1], "y")

	assertStatementsLen(t, function.Body.Statements, 1)
	bodyStmt := assertStatementType[*ast.ExpressionStatement](t, function.Body.Statements[0])
	assertInfixExpresssion(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParamaterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)

		stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
		function := assertExpressionType[*ast.FunctionLiteral](t, stmt.Expression)
		assertParametersLen(t, function.Parameters, len(tt.expectedParams))
		for i, ident := range tt.expectedParams {
			assertLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p)
	assertStatementsLen(t, program.Statements, 1)

	stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
	exp := assertExpressionType[*ast.CallExpression](t, stmt.Expression)

	assertIdentifier(t, exp.Function, "add")
	assertExpressionsLen(t, exp.Arguments, 3)

	assertLiteralExpression(t, exp.Arguments[0], 1)
	assertInfixExpresssion(t, exp.Arguments[1], 2, "*", 3)
	assertInfixExpresssion(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedIdent  string
		expectedParams []string
	}{
		{input: "add();", expectedIdent: "add", expectedParams: []string{}},
		{input: "add(1);", expectedIdent: "add", expectedParams: []string{"1"}},
		{input: "add(1, 2 * 3, 4 + 5);", expectedIdent: "add", expectedParams: []string{"1", "(2 * 3)", "(4 + 5)"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p)
		stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
		exp := assertExpressionType[*ast.CallExpression](t, stmt.Expression)

		assertIdentifier(t, exp.Function, tt.expectedIdent)

		assertExpressionsLen(t, exp.Arguments, len(tt.expectedParams))
		for i, ident := range tt.expectedParams {
			assertEquals(t, exp.Arguments[i].String(), ident)
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello, world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p)
	stmt := assertStatementType[*ast.ExpressionStatement](t, program.Statements[0])
	literal := assertExpressionType[*ast.StringLiteral](t, stmt.Expression)
	assertEquals(t, literal.Value, "hello, world")
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

func assertParametersLen(t testing.TB, parameters []*ast.Identifier, want int) {
	t.Helper()

	got := len(parameters)
	if got != want {
		t.Fatalf("wrong number of parameters. got=%d want=%d", got, want)
	}
}

func assertExpressionsLen(t testing.TB, arguments []ast.Expression, want int) {
	t.Helper()

	got := len(arguments)
	if got != want {
		t.Fatalf("wrong number of expressions. got=%d want=%d", got, want)
	}
}

func assertLetStatement(t testing.TB, stmt ast.Statement, expectedIdent string, expectedValue any) {
	t.Helper()

	assertEquals(t, stmt.TokenLiteral(), "let")
	letStmt := assertStatementType[*ast.LetStatement](t, stmt)

	assertEquals(t, letStmt.Name.Value, expectedIdent)
	assertEquals(t, letStmt.Name.TokenLiteral(), expectedIdent)

	val := letStmt.Value
	assertLiteralExpression(t, val, expectedValue)
}
