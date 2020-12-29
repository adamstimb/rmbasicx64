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
			},
		},
		{
			Source: "Print \"Illegal char\" {",
			ExpectedTokens: []Token{
				{PRINT, "PRINT"},
				{StringLiteral, "Illegal char"},
				{Illegal, "{"},
			},
		},
		{
			Source: "Rem This is a comment",
			ExpectedTokens: []Token{
				{REM, "REM"},
				{Comment, "This is a comment"},
			},
		},
		{
			Source: "this% :=that$+ meh - 5.1234",
			ExpectedTokens: []Token{
				{Identifier, "This%"},
				{Assign, ":="},
				{Identifier, "That$"},
				{Plus, "+"},
				{Identifier, "Meh"},
				{Minus, "-"},
				{NumericalLiteral, "5.1234"},
			},
		},
		{
			Source: "my_var% := yet_more_var% + foo_var",
			ExpectedTokens: []Token{
				{Identifier, "My_var%"},
				{Assign, ":="},
				{Identifier, "Yet_more_var%"},
				{Plus, "+"},
				{Identifier, "Foo_var"},
			},
		},
	}

	// test that we always get expected tokens
	s := &Scanner{}
	for i, test := range tests {
		tokens := s.ScanTokens(test.Source)
		for j, token := range tokens {
			if token.TokenType != test.ExpectedTokens[j].TokenType {
				t.Fatalf("TestScanner [%d]: TokenType [%d] expected, got [%d] from source [%q]", i, test.ExpectedTokens[j].TokenType, token.TokenType, test.Source)
			}
			if token.Literal != test.ExpectedTokens[j].Literal {
				t.Fatalf("TestScanner [%d]: Literal [%q] expected, got [%q] from source [%q]", i, test.ExpectedTokens[j].Literal, token.Literal, test.Source)
			}
		}
	}
}
