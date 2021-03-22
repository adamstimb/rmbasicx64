package main

import "testing"

func TestLetStatements(t *testing.T) {

	// test data
	input := `let x = 5;
let y = 10;
let foobar = 1234` // we'll change this to BASIC when we see it works

	// scan for tokens
	s := &Scanner{}
	tokens := s.ScanTokens(input)

	// new parser
	p := &Parser{}
	p.New(tokens)

	// get program
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	// test identifiers
	tests := []struct{ expectedIdentifier string }{
		{"X"},
		{"Y"},
		{"Foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		println(i)
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		println("next")
	}
}

func testLetStatement(t *testing.T, s Statement, name string) bool {
	println("testLetStatement")
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not LET, got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*LetStatement)
	if !ok {
		t.Errorf("s not *LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s, got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not %s, got %s", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}
