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
			Source: "",
			ExpectedTokens: []Token{
				{EndOfLine, ""},
			},
		},
		{
			Source: "\n",
			ExpectedTokens: []Token{
				{EndOfLine, ""},
			},
		},
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
			Source: "Print \"So-called \"\"test\"\" this is\"",
			ExpectedTokens: []Token{
				{PRINT, "PRINT"},
				{StringLiteral, "So-called \"\"test\"\" this is"},
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
			Source: "&554a3d2c",
			ExpectedTokens: []Token{
				{HexLiteral, "&554A3D2C"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "let y = &5d",
			ExpectedTokens: []Token{
				{LET, "LET"},
				{IdentifierLiteral, "Y"},
				{Equal, "="},
				{HexLiteral, "&5D"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "rM_baSic_HAD_thiS_Weird_camel_case_tHING_GoInG_On$",
			ExpectedTokens: []Token{
				{IdentifierLiteral, "Rm_Basic_Had_This_Weird_Camel_Case_Thing_Going_On$"},
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
				{IdentifierLiteral, "My_Var%"},
				{Assign, ":="},
				{IdentifierLiteral, "Yet_More_Var%"},
				{Plus, "+"},
				{IdentifierLiteral, "Foo_Var"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "50 > 50",
			ExpectedTokens: []Token{
				{NumericalLiteral, "50"},
				{GreaterThan, ">"},
				{NumericalLiteral, "50"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "50 <= 50",
			ExpectedTokens: []Token{
				{NumericalLiteral, "50"},
				{LessThanEqualTo1, "<="},
				{NumericalLiteral, "50"},
				{EndOfLine, ""},
			},
		},
		{
			Source: "50 >= 50",
			ExpectedTokens: []Token{
				{NumericalLiteral, "50"},
				{GreaterThanEqualTo1, ">="},
				{NumericalLiteral, "50"},
				{EndOfLine, ""},
			},
		},
	}

	// test that we always get expected tokens
	s := &Scanner{}
	for i, test := range tests {
		tokens := s.Scan(test.Source)
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
