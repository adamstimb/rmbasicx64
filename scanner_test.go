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
				{StringLiteral, "Hello\""},
			},
		},
	}

	// run tests
	s := &Scanner{}
	for i, test := range tests {
		tokens := s.ScanTokens(test.Source)
		for j, token := range tokens {
			if token.TokenType != test.ExpectedTokens[j].TokenType {
				t.Fatalf("Test [%d]: TokenType [%d] expected, got [%d] from source [%q]", i, test.ExpectedTokens[j].TokenType, token.TokenType, test.Source)
			}
		}
	}
}
