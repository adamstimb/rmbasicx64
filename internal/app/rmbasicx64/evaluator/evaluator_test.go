package evaluator

import (
	"log"
	"testing"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/game"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/lexer"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/object"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/parser"
)

func TestEvalNumericExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"--10", 10},
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
		evaluated := testEval(tt.input)
		testNumericObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := &lexer.Lexer{}
	l.Scan(input)
	p := parser.New(l, &game.Game{})
	program := p.ParseProgram()
	env := object.NewEnvironment(nil)
	g := &game.Game{}
	return Eval(g, program, env)
}

func testNumericObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Numeric)
	if !ok {
		t.Errorf("object is not float64, got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got %f, want %f", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"true", -1.0},
		{"false", 0},
		{"1 < 2", -1.0},
		{"1 > 2", 0},
		{"1 < 1", 0},
		{"1 > 1", 0},
		{"NOT TRUE", 0},
		{"NOT FALSE", -1.0},
		{"NOT 1 = 1", 0},
		{"1 = 1", -1.0},
		{"1 = 2", 0},
		{"true = true", -1.0},
		{"false = false", -1.0},
		{"true = false", 0},
		{"(1 < 2) = true", -1.0},
		{"(1 < 2) = false", 0},
		{"(1 > 2) = true", 0},
		{"(1 > 2) = false", -1.0},
		{"true and true", -1.0},
		{"true and false", 0},
		{"true or false", -1.0},
		{"true xor false", -1.0},
		{"\"this\" = \"that\"", 0},
		{"\"that\" = \"that\"", -1.0},
		{"\"this\" == \"THIS\"", -1.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumericObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got %t, want %t", result.Value, expected)
		return false
	}
	return true
}

func TestNotOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"not true", 0},
		{"not false", -1},
		{"not not true", -1},
		{"not not false", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumericObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"IF TRUE THEN TRUE", -1.0},
		{"IF -1.0 THEN TRUE", -1.0},
		{"IF 1 < 2 THEN TRUE", -1.0},
		{"IF 1 > 2 THEN TRUE ELSE FALSE", 0.0},
		{"IF 1 < 2 THEN TRUE ELSE FALSE", -1.0},
	}

	for _, tt := range tests {
		log.Printf("input: %s", tt.input)
		evaluated := testEval(tt.input)
		testNumericObject(t, evaluated, tt.expected)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got %T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"return 10;", 10}, // <-- this should fail
		{`return 10
9
`, 10},
		{`return 2 * 5 : 9
`, 10},
		{`
9
return 2 * 5 : 9
`, 10},
		{`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}
	return 1;
}
`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumericObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;", // <-- again, this should fail
			"type mismatch: NUMERIC + BOOLEAN",
		},
		{
			"5 + true : 5",
			"type mismatch: NUMERIC + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5 : true + false : 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}
	return 1;
}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"identifier not found: Foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got %T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message, expected %q, got %q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`let a = 5 : a`, 5},
		{`let a = 5 * 5
a
`, 25},
		{"let a = 5 : let b = a : b", 5},
		{`let a = 5
let b = a : let c = a + b + 5
c
`, 15},
	}

	for _, tt := range tests {
		testNumericObject(t, testEval(tt.input), tt.expected)
	}
}

func TestBindStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`a := 5 : a`, 5},
		{`a := 5 * 5
a
`, 25},
		{"a := 5 : b := a : b", 5},
		{`a := 5
b := a : c := a + b + 5
c
`, 15},
	}

	for _, tt := range tests {
		testNumericObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFloatToIntegerBind(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"let a = 1.23456 : a", 1.23456},
		{"let a% = 1.23456 : a%", 1.0},
		{"let a% = 1.99999 : a%", 1.0},
	}

	for _, tt := range tests {
		testNumericObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "function(x) : x + 2 : endfun"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function, got %T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters.Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "X" {
		t.Fatalf("parameter is not X, got %q", fn.Parameters[0])
	}
	expectedBody := "(X + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q, got %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"FUNCTION identity(x) : x : result 5 : endfun : Identity(5)", 5},
		//{"LET identity := function(x) : result x : endfun : Identity(5)", 5},
		//{"let double := function(x) : x * 2 : endfun : Double(5)", 10},
		//{"let add = function(x, y) : x + y : endfun : Add(5, 5)", 10},
		//{"let add = function(x, y) : x + y : endfun : Add(5 + 5, Add(5, 5))", 20},
		//{"function(x) : x : endfun : (5)", 5},
	}

	for _, tt := range tests {
		log.Printf("input: %s\n", tt.input)
		testNumericObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello world!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello world!" {
		t.Errorf("String has wrong value, expected %q, got %q", "Hello world!", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value, expected %q, got %q", "Hello World!", str.Value)
	}
}

func TestClosures(t *testing.T) {
	input := `let newAdder = function(x) : function(y) : x + y : endfun : endfun : let addTwo = newAdder(2) : addTwo(2)`
	testNumericObject(t, testEval(input), 4)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `LEN` not supported, got NUMERIC"},
		{`len("one", "two")`, "wrong number of arguments, got 2, want 1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testNumericObject(t, evaluated, float64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error, got %T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message, expected %q, got %q", expected, errObj.Message)
			}
		}
	}
}
