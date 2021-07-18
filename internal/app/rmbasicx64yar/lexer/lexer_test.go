package lexer

import (
	"testing"

	"github.com/adamstimb/rmbasicx64yar/internal/app/rmbasicx64yar/token"
)

func TestLexer(t *testing.T) {

	// test data
	type test struct {
		Source         string
		ExpectedTokens []token.Token
	}
	tests := []test{
		{
			Source: "",
			ExpectedTokens: []token.Token{
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "\n",
			ExpectedTokens: []token.Token{
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "print \"Hello!\"",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.StringLiteral, Literal: "Hello!"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "Print \"So-called \"\"test\"\" this is\"",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.StringLiteral, Literal: "So-called \"\"test\"\" this is"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "Rem This is a comment",
			ExpectedTokens: []token.Token{
				{TokenType: token.REM, Literal: "REM"},
				{TokenType: token.Comment, Literal: "This is a comment"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "let x = 5",
			ExpectedTokens: []token.Token{
				{TokenType: token.LET, Literal: "LET"},
				{TokenType: token.IdentifierLiteral, Literal: "X"},
				{TokenType: token.Equal, Literal: "="},
				{TokenType: token.NumericLiteral, Literal: "5"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "&554a3d2c",
			ExpectedTokens: []token.Token{
				{TokenType: token.HexLiteral, Literal: "&554A3D2C"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "let y = &5d",
			ExpectedTokens: []token.Token{
				{TokenType: token.LET, Literal: "LET"},
				{TokenType: token.IdentifierLiteral, Literal: "Y"},
				{TokenType: token.Equal, Literal: "="},
				{TokenType: token.HexLiteral, Literal: "&5D"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "rM_baSic_HAD_thiS_Weird_camel_case_tHING_GoInG_On$",
			ExpectedTokens: []token.Token{
				{TokenType: token.IdentifierLiteral, Literal: "Rm_Basic_Had_This_Weird_Camel_Case_Thing_Going_On$"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "this% :=that$+ meh - 5.1234",
			ExpectedTokens: []token.Token{
				{TokenType: token.IdentifierLiteral, Literal: "This%"},
				{TokenType: token.Assign, Literal: ":="},
				{TokenType: token.IdentifierLiteral, Literal: "That$"},
				{TokenType: token.Plus, Literal: "+"},
				{TokenType: token.IdentifierLiteral, Literal: "Meh"},
				{TokenType: token.Minus, Literal: "-"},
				{TokenType: token.NumericLiteral, Literal: "5.1234"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "my_var% := yet_more_var% + foo_var",
			ExpectedTokens: []token.Token{
				{TokenType: token.IdentifierLiteral, Literal: "My_Var%"},
				{TokenType: token.Assign, Literal: ":="},
				{TokenType: token.IdentifierLiteral, Literal: "Yet_More_Var%"},
				{TokenType: token.Plus, Literal: "+"},
				{TokenType: token.IdentifierLiteral, Literal: "Foo_Var"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "50 > 50",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "50"},
				{TokenType: token.GreaterThan, Literal: ">"},
				{TokenType: token.NumericLiteral, Literal: "50"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "50 <= 50",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "50"},
				{TokenType: token.LessThanEqualTo1, Literal: "<="},
				{TokenType: token.NumericLiteral, Literal: "50"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "50 >= 50",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "50"},
				{TokenType: token.GreaterThanEqualTo1, Literal: ">="},
				{TokenType: token.NumericLiteral, Literal: "50"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "-1 AND 1",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "-1"},
				{TokenType: token.AND, Literal: "AND"},
				{TokenType: token.NumericLiteral, Literal: "1"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "-1.0",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "-1.0"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "NOT -1",
			ExpectedTokens: []token.Token{
				{TokenType: token.NOT, Literal: "NOT"},
				{TokenType: token.NumericLiteral, Literal: "-1"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "(-1)",
			ExpectedTokens: []token.Token{
				{TokenType: token.LeftParen, Literal: "("},
				{TokenType: token.NumericLiteral, Literal: "-1"},
				{TokenType: token.RightParen, Literal: ")"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "PRINT #1 !",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.Hash, Literal: "#"},
				{TokenType: token.NumericLiteral, Literal: "1"},
				{TokenType: token.Exclamation, Literal: "!"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "PRINT 1,2",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.NumericLiteral, Literal: "1"},
				{TokenType: token.Comma, Literal: ","},
				{TokenType: token.NumericLiteral, Literal: "2"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "PRINT ~1 !",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.Tilde, Literal: "~"},
				{TokenType: token.NumericLiteral, Literal: "1"},
				{TokenType: token.Exclamation, Literal: "!"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "PRINT 2.34E+4",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.NumericLiteral, Literal: "2.34E+4"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: "PRINT 1.344E-4.32",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.NumericLiteral, Literal: "1.344E-4.32"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: `LEN("HELLO")`,
			ExpectedTokens: []token.Token{
				{TokenType: token.IdentifierLiteral, Literal: "LEN"},
				{TokenType: token.LeftParen, Literal: "("},
				{TokenType: token.StringLiteral, Literal: "HELLO"},
				{TokenType: token.RightParen, Literal: ")"},
				{TokenType: token.EOF, Literal: token.EOF},
			},
		},
		{
			Source: `set mode 80
cls : print "Blergh"
print "Hello"
`,
			ExpectedTokens: []token.Token{
				{TokenType: token.SET, Literal: "SET"},
				{TokenType: token.MODE, Literal: "MODE"},
				{TokenType: token.NumericLiteral, Literal: "80"},
				{TokenType: token.NewLine, Literal: "\n"},
				{TokenType: token.CLS, Literal: "CLS"},
				{TokenType: token.Colon, Literal: ":"},
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.StringLiteral, Literal: "Blergh"},
				{TokenType: token.NewLine, Literal: "\n"},
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.StringLiteral, Literal: "Hello"},
				{TokenType: token.NewLine, Literal: "\n"},
				{TokenType: token.EOF, Literal: "EOF"},
			},
		},
		{
			Source: `10 set mode 80
20 cls : print "Blergh"
30 print "Hello"
`,
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "10"},
				{TokenType: token.SET, Literal: "SET"},
				{TokenType: token.MODE, Literal: "MODE"},
				{TokenType: token.NumericLiteral, Literal: "80"},
				{TokenType: token.NewLine, Literal: "\n"},
				{TokenType: token.NumericLiteral, Literal: "20"},
				{TokenType: token.CLS, Literal: "CLS"},
				{TokenType: token.Colon, Literal: ":"},
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.StringLiteral, Literal: "Blergh"},
				{TokenType: token.NewLine, Literal: "\n"},
				{TokenType: token.NumericLiteral, Literal: "30"},
				{TokenType: token.PRINT, Literal: "PRINT"},
				{TokenType: token.StringLiteral, Literal: "Hello"},
				{TokenType: token.NewLine, Literal: "\n"},
				{TokenType: token.EOF, Literal: "EOF"},
			},
		},
	}

	// test that we always get expected tokens
	l := &Lexer{}
	for i, test := range tests {
		tokens := l.Scan(test.Source)
		for j, tok := range tokens {
			if tok.TokenType != test.ExpectedTokens[j].TokenType {
				token.PrintToken(tok)
				t.Fatalf("Token [%d]: TokenType [%s] expected, got [%s] from source [%q]", i, test.ExpectedTokens[j].TokenType, tok.TokenType, test.Source)
			}
			if tok.Literal != test.ExpectedTokens[j].Literal {
				token.PrintToken(tok)
				t.Fatalf("Token [%d]: Literal [%q] expected, got [%q] from source [%q]", i, test.ExpectedTokens[j].Literal, tok.Literal, test.Source)
			}
		}
	}
}
