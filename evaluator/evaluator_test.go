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
	}

	for _, tt := range tests {
		evaluated := prepareObject(tt.input)
		assertIntegerObject(t, evaluated, tt.expected)
	}
}

func prepareObject(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

//helpers

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
