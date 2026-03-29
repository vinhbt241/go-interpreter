package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := prepareObject(tt.input)
		assertIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := prepareObject(tt.input)
		assertBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := prepareObject(tt.input)
		assertBooleanObject(t, evaluated, tt.expected)
	}
}

//helpers

func prepareObject(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func assertObjectType[T object.Object](t testing.TB, got object.Object) T {
	t.Helper()

	obj, ok := got.(T)
	if !ok {
		t.Errorf("wrong object type. got=%T want %T", *new(T), got)
	}

	return obj
}

func assertIntegerObject(t testing.TB, got object.Object, want int64) {
	t.Helper()

	result := assertObjectType[*object.Integer](t, got)

	if result.Value != want {
		t.Errorf("wrong value. got=%d want=%d", result.Value, want)
	}
}

func assertBooleanObject(t testing.TB, got object.Object, want bool) {
	t.Helper()

	result := assertObjectType[*object.Boolean](t, got)

	if result.Value != want {
		t.Errorf("wrong value. got=%t want=%t", result.Value, want)
	}
}
