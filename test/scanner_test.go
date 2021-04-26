package main

import (
	"testing"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

func TestScanner(t *testing.T) {

	// test data
	type test struct {
		Source         string
		ExpectedTokens []token.Token
	}
	tests := []test{
		{
			Source: "",
			ExpectedTokens: []token.Token{
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "\n",
			ExpectedTokens: []token.Token{
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "print \"Hello!\"",
			ExpectedTokens: []token.Token{
				{token.PRINT, "PRINT"},
				{token.StringLiteral, "Hello!"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "Print \"Illegal char\" {",
			ExpectedTokens: []token.Token{
				{token.PRINT, "PRINT"},
				{token.StringLiteral, "Illegal char"},
				{token.Illegal, "{"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "Print \"So-called \"\"test\"\" this is\"",
			ExpectedTokens: []token.Token{
				{token.PRINT, "PRINT"},
				{token.StringLiteral, "So-called \"\"test\"\" this is"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "Rem This is a comment",
			ExpectedTokens: []token.Token{
				{token.REM, "REM"},
				{token.Comment, "This is a comment"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "let x = 5",
			ExpectedTokens: []token.Token{
				{token.LET, "LET"},
				{token.IdentifierLiteral, "X"},
				{token.Equal, "="},
				{token.NumericalLiteral, "5"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "&554a3d2c",
			ExpectedTokens: []token.Token{
				{token.HexLiteral, "&554A3D2C"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "let y = &5d",
			ExpectedTokens: []token.Token{
				{token.LET, "LET"},
				{token.IdentifierLiteral, "Y"},
				{token.Equal, "="},
				{token.HexLiteral, "&5D"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "rM_baSic_HAD_thiS_Weird_camel_case_tHING_GoInG_On$",
			ExpectedTokens: []token.Token{
				{token.IdentifierLiteral, "Rm_Basic_Had_This_Weird_Camel_Case_Thing_Going_On$"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "this% :=that$+ meh - 5.1234",
			ExpectedTokens: []token.Token{
				{token.IdentifierLiteral, "This%"},
				{token.Assign, ":="},
				{token.IdentifierLiteral, "That$"},
				{token.Plus, "+"},
				{token.IdentifierLiteral, "Meh"},
				{token.Minus, "-"},
				{token.NumericalLiteral, "5.1234"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "my_var% := yet_more_var% + foo_var",
			ExpectedTokens: []token.Token{
				{token.IdentifierLiteral, "My_Var%"},
				{token.Assign, ":="},
				{token.IdentifierLiteral, "Yet_More_Var%"},
				{token.Plus, "+"},
				{token.IdentifierLiteral, "Foo_Var"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "50 > 50",
			ExpectedTokens: []token.Token{
				{token.NumericalLiteral, "50"},
				{token.GreaterThan, ">"},
				{token.NumericalLiteral, "50"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "50 <= 50",
			ExpectedTokens: []token.Token{
				{token.NumericalLiteral, "50"},
				{token.LessThanEqualTo1, "<="},
				{token.NumericalLiteral, "50"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "50 >= 50",
			ExpectedTokens: []token.Token{
				{token.NumericalLiteral, "50"},
				{token.GreaterThanEqualTo1, ">="},
				{token.NumericalLiteral, "50"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "-1 AND 1",
			ExpectedTokens: []token.Token{
				{token.NumericalLiteral, "-1"},
				{token.AND, "AND"},
				{token.NumericalLiteral, "1"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "-1.0",
			ExpectedTokens: []token.Token{
				{token.NumericalLiteral, "-1.0"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "NOT -1",
			ExpectedTokens: []token.Token{
				{token.NOT, "NOT"},
				{token.NumericalLiteral, "-1"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "(-1)",
			ExpectedTokens: []token.Token{
				{token.LeftParen, "("},
				{token.NumericalLiteral, "-1"},
				{token.RightParen, ")"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "PRINT #1 !",
			ExpectedTokens: []token.Token{
				{token.PRINT, "PRINT"},
				{token.Hash, "#"},
				{token.NumericalLiteral, "1"},
				{token.Exclamation, "!"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "PRINT 1,2",
			ExpectedTokens: []token.Token{
				{token.PRINT, "PRINT"},
				{token.NumericalLiteral, "1"},
				{token.Comma, ","},
				{token.NumericalLiteral, "2"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "PRINT ~1 !",
			ExpectedTokens: []token.Token{
				{token.PRINT, "PRINT"},
				{token.Tilde, "~"},
				{token.NumericalLiteral, "1"},
				{token.Exclamation, "!"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "PRINT 2.34E+4",
			ExpectedTokens: []token.Token{
				{token.PRINT, "PRINT"},
				{token.NumericalLiteral, "2.34E+4"},
				{token.EndOfLine, ""},
			},
		},
		{
			Source: "PRINT 1.344E-4.32",
			ExpectedTokens: []token.Token{
				{token.PRINT, "PRINT"},
				{token.NumericalLiteral, "1.344E-4.32"},
				{token.EndOfLine, ""},
			},
		},
	}

	// test that we always get expected tokens
	s := &rmbasicx64.Scanner{}
	for i, test := range tests {
		tokens := s.Scan(test.Source)
		for j, tok := range tokens {
			if tok.TokenType != test.ExpectedTokens[j].TokenType {
				token.PrintToken(tok)
				t.Fatalf("Token [%d]: TokenType [%d] expected, got [%d] from source [%q]", i, test.ExpectedTokens[j].TokenType, tok.TokenType, test.Source)
			}
			if tok.Literal != test.ExpectedTokens[j].Literal {
				token.PrintToken(tok)
				t.Fatalf("Token [%d]: Literal [%q] expected, got [%q] from source [%q]", i, test.ExpectedTokens[j].Literal, tok.Literal, test.Source)
			}
		}
	}
}
