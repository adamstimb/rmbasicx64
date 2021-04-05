package main

import "testing"

func TestScanner(t *testing.T) {

	// test data
	type test struct {
		Source         string
		ExpectedTokens []Token
	}
	tests := []test{
		{
			Source: "print \"Hello!\"",
			ExpectedTokens: []Token{
				{PRINT, "PRINT"},
				{StringLiteral, "Hello!"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "Print \"Illegal char\" {",
			ExpectedTokens: []Token{
				{PRINT, "PRINT"},
				{StringLiteral, "Illegal char"},
				{Illegal, "{"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "Rem This is a comment",
			ExpectedTokens: []Token{
				{REM, "REM"},
				{Comment, "This is a comment"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "let x = 5",
			ExpectedTokens: []Token{
				{LET, "LET"},
				{IdentifierLiteral, "X"},
				{Equal, "="},
				{NumericalLiteral, "5"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "this% :=that$+ meh - 5.1234",
			ExpectedTokens: []Token{
				{IdentifierLiteral, "This%"},
				{Assign, ":="},
				{IdentifierLiteral, "That$"},
				{Plus, "+"},
				{IdentifierLiteral, "Meh"},
				{Minus, "-"},
				{NumericalLiteral, "5.1234"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "my_var% := yet_more_var% + foo_var",
			ExpectedTokens: []Token{
				{IdentifierLiteral, "My_var%"},
				{Assign, ":="},
				{IdentifierLiteral, "Yet_more_var%"},
				{Plus, "+"},
				{IdentifierLiteral, "Foo_var"},
				{EndOfLine, ""},
			},
		},
	}

	// test that we always get expected tokens
	s := &Scanner{}
	for i, test := range tests {
		tokens := s.ScanTokens(test.Source)
		for j, token := range tokens {
			if token.TokenType != test.ExpectedTokens[j].TokenType {
				t.Fatalf("Token [%d]: TokenType [%d] expected, got [%d] from source [%q]", i, test.ExpectedTokens[j].TokenType, token.TokenType, test.Source)
			}
			if token.Literal != test.ExpectedTokens[j].Literal {
				t.Fatalf("Token [%d]: Literal [%q] expected, got [%q] from source [%q]", i, test.ExpectedTokens[j].Literal, token.Literal, test.Source)
			}
		}
	}
}