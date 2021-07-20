package parser

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/ast"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/lexer"
)

func TestLine(t *testing.T) {

	tests := []struct {
		input             string
		expectedStmtCount int
		expectedString    string
	}{
		{"let x := 5\n", 1, "LET X := 5"},
		{"let x := 5 : let x := 5\n", 2, "LET X := 5 : LET X := 5"},
		{"let x := 5 : let x := 5 : \n", 3, "LET X := 5 : LET X := 5 : "},
		{"let x := 5 : let y := 5 + x: \n", 3, "LET X := 5 : LET Y := (5 + X) : "},
		{"let x := 5 : y := 5 + x: \n", 3, "LET X := 5 : Y := (5 + X) : "},
		{"let x = 5 : y =5 + x\n", 2, "LET X = 5 : Y = (5 + X)"},
	}

	for _, tt := range tests {
		l := &lexer.Lexer{}
		l.Scan(tt.input)
		p := New(l)

		line := p.ParseLine()
		if line == nil {
			t.Fatalf("parseLine() returned nil")
		}
		if len(line.Statements) != tt.expectedStmtCount {
			t.Fatalf("line.Statements does not contain %d statements, got %d", tt.expectedStmtCount, len(line.Statements))
		}
		if line.String() != tt.expectedString {
			t.Fatalf("line.String is not %q, got %q", tt.expectedString, line.String())
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "LET" {
		t.Errorf("s.TokenLiteral not 'LET', got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not ast.LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s' got '%s'", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s' got '%s'", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestBindStatements(t *testing.T) {

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x := 5\n", "X", 5},
		{"y = true\n", "Y", -1.0},
		{"foobar := y\n", "Foobar", "Y"},
	}

	for _, tt := range tests {
		l := &lexer.Lexer{}
		l.Scan(tt.input)
		p := New(l)

		stmt := p.parseStatement()

		bindStmt, ok := stmt.(*ast.BindStatement)
		if !ok {
			t.Fatalf("s not ast.BindStatement, got %T", stmt)
		}

		if tt.expectedIdentifier != bindStmt.Name.Value {
			t.Fatalf("Expected identifier name %q, got %q", tt.expectedIdentifier, bindStmt.Name.Value)
		}

		val := stmt.(*ast.BindStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestLetStatements(t *testing.T) {

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x := 5\n", "X", 5},
		{"let y = true\n", "Y", true},
		{"let foobar := y\n", "Foobar", "Y"},
	}

	for _, tt := range tests {
		l := &lexer.Lexer{}
		l.Scan(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements, got %d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

// -------------------------------------------------------------------------
// -- RETURN (TODO: This is Monkey implementation, not Basic)

func TestResultStatements(t *testing.T) {
	input := `result 5 : result 10
result 993322
`
	l := &lexer.Lexer{}
	l.Scan(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		resultStmt, ok := stmt.(*ast.ResultStatement)
		if !ok {
			t.Errorf("stmt not *ast.ResultStatement, got %T", stmt)
			continue
		}
		if resultStmt.TokenLiteral() != "RESULT" {
			t.Errorf("returnStmt.TokenLiteral is not 'RESULT', got %q", resultStmt.TokenLiteral())
		}
	}
}

// -------------------------------------------------------------------------
// -- Identifier Expression

func TestIdentifierExpression(t *testing.T) {
	input := `foobar
`

	l := &lexer.Lexer{}
	l.Scan(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "Foobar" {
		t.Errorf("ident.Value not %s, got %s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "Foobar" {
		t.Errorf("ident.TokenLiteral not %s, got %s", "foobar", ident.TokenLiteral())
	}
}

// -------------------------------------------------------------------------
// -- Numeric Literal

func TestNumericLiteralExpression(t *testing.T) {
	input := `5`

	l := &lexer.Lexer{}
	l.Scan(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.NumericLiteral)

	if !ok {
		t.Fatalf("expression not *ast.NumericLiteral, got %T", stmt.Expression)
	}

	if literal.Value != float64(5) {
		t.Errorf("literal.Value not %f, got %f", float64(5), float64(literal.Value))
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s, got %s", "5", literal.TokenLiteral())
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("expression not *ast.Identifier, got %T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s, got %s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s, got %s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testNumericLiteral(t *testing.T, il ast.Expression, value float64) bool {
	numeric, ok := il.(*ast.NumericLiteral)
	if !ok {
		t.Errorf("il not *ast.NumericLiteral, got %T", il)
		return false
	}
	if numeric.Value != value {
		t.Errorf("numeric.Value not %f, got %f", value, numeric.Value)
		return false
	}

	literalvalue, err := strconv.ParseFloat(numeric.TokenLiteral(), 64)
	if err != nil {
		t.Fatalf("could not parse %s as float", numeric.TokenLiteral())
	}

	if literalvalue != value {
		t.Errorf("numeric.TokenLiteral not %f, got %s", value, numeric.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, il ast.Expression, value string) bool {
	stringval, ok := il.(*ast.StringLiteral)
	if !ok {
		t.Errorf("il not *ast.StringLiteral, got %T", il)
		return false
	}
	if stringval.Value != value {
		t.Errorf("stringval.Value not %q, got %q", value, stringval.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("expression not *ast.Boolean, got %T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t, got %t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != strings.ToUpper(fmt.Sprintf("%t", value)) {
		t.Errorf("bo.TokenLiteral not %t, got %s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testNumericLiteral(t, exp, float64(v))
	case int64:
		return testNumericLiteral(t, exp, float64(v))
	case float64:
		return testNumericLiteral(t, exp, v)
	case string:
		// lol yuck
		_, ok := exp.(*ast.Identifier)
		if ok {
			return testIdentifier(t, exp, v)
		} else {
			return testStringLiteral(t, exp, v)
		}
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of expression not handled, got %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got %T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		log.Println("booyah left")
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %s, got %q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		log.Println("booyah right")
		return false
	}
	return true
}

// -------------------------------------------------------------------------
// -- Prefix Expression

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"NOT 5", "NOT", 5},
		//{"-15;", "-", 15},   // TODO: Clarify how we want to do - numbers
		{"not true", "NOT", -1.0},
		{"not false", "NOT", 0},
	}
	for _, tt := range prefixTests {
		l := &lexer.Lexer{}
		l.Scan(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got %d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression, got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s, got %s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

// -------------------------------------------------------------------------
// -- Infix Expression

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 = 5", 5, "=", 5},
		{"true = true", -1.0, "=", -1.0},
		{"false = false", 0, "=", 0},
		{"true and true", -1.0, "AND", -1.0},
		{"true or false", -1.0, "OR", 0},
		{"false xor true", 0, "XOR", -1},
		{"\"this\" = \"that\"", "this", "=", "that"},
		{"\"this\" == \"THIS\"", "this", "==", "THIS"},
	}

	for _, tt := range infixTests {
		l := &lexer.Lexer{}
		l.Scan(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
		}

		log.Println("got here")
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			log.Println("this failed")
			return
		}
		log.Println("got to end")
	}
}

// -------------------------------------------------------------------------
// -- Precedence

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-A) * B)",
		},
		{
			"a + b + c",
			"((A + B) + C)",
		},
		{
			"a + b - c",
			"((A + B) - C)",
		},
		{
			"a * b * c",
			"((A * B) * C)",
		},
		{
			"a + b / c",
			"(A + (B / C))",
		},
		{
			"a + b * c + d / e - f",
			"(((A + (B * C)) + (D / E)) - F)",
		},
		{
			"5 > 4 = 3 < 4",
			"((5 > 4) = (3 < 4))",
		},
		{
			"true",
			"TRUE",
		},
		{
			"false",
			"FALSE",
		},
		{
			"3 > 5 = false",
			"((3 > 5) = FALSE)",
		},
		{
			"3 < 5 = true",
			"((3 < 5) = TRUE)",
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
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		//{
		//	"-5(5 + 5)",
		//	"(-5(5 + 5))",
		//}
		{
			"NOT(true = true)",
			"(NOT(TRUE = TRUE))",
		},
		{
			"a + add(b * c) + d",
			"((A + Add((B * C))) + D)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"Add(A,B,1,(2 * 3),(4 + 5),Add(6,(7 * 8)))",
		},
	}

	for _, tt := range tests {
		l := &lexer.Lexer{}
		l.Scan(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected %q got %q", tt.expected, actual)
		}
	}
}

// -------------------------------------------------------------------------
// -- Boolean

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean float64
	}{
		{"true", -1.0},
		{"false", 0},
	}

	for _, tt := range tests {
		l := &lexer.Lexer{}
		l.Scan(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.NumericLiteral)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %g. got=%g", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

// -------------------------------------------------------------------------
// -- If

func TestIfExpression(t *testing.T) {
	input := "if x < y then x = 8\n"
	l := &lexer.Lexer{}
	l.Scan(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d\n", 1, len(program.Statements))
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if x < y then x = 0 : b = 3 else y = 2\n"
	l := &lexer.Lexer{}
	l.Scan(input)
	p := New(l)
	program := p.ParseLine()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		for _, stmt := range program.Statements {
			log.Println(stmt.String())
		}
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}
}

// -------------------------------------------------------------------------
// -- Function Definition

func TestFunctionDefinitionParsing(t *testing.T) {

	input := `function add(x, y) : x + y ENDFUN`
	l := &lexer.Lexer{}
	l.Scan(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionDefinition)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionDefinition, got %T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong, want 2, got %d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "X")
	testLiteralExpression(t, function.Parameters[1], "Y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements, got %d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement, got %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "X", "+", "Y")

	tests := []struct {
		input          string
		expectedParams []string
	}{
		{
			input:          "function foobar() : endfun",
			expectedParams: []string{},
		},
		{
			input: `function foobar(x)
endfun`,
			expectedParams: []string{"X"},
		},
		{
			input: `function foobar(x, y, z) : ENDFUN
`,
			expectedParams: []string{"X", "Y", "Z"},
		},
	}

	for _, tt := range tests {
		l := &lexer.Lexer{}
		l.Scan(tt.input)
		p := New(l)
		program := p.ParseProgram()
		log.Println(tt.input)
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionDefinition)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("wanted %d parameters, got %d\n", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

// -------------------------------------------------------------------------
// -- Call expression

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"
	l := &lexer.Lexer{}
	l.Scan(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression, got %T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "Add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong lengt of arguments, got %d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`
	l := &lexer.Lexer{}
	l.Scan(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expression not *ast.StringLiteral, got %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q, got %q", "hello world", literal.Value)
	}
}
