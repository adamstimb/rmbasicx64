package lexer

import (
	"testing"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
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
				{TokenType: token.EOF, Literal: token.EOF, Index: 0},
			},
		},
		{
			Source: "\n",
			ExpectedTokens: []token.Token{
				{TokenType: token.EOF, Literal: token.EOF, Index: 0},
			},
		},
		{
			Source: "print \"Hello!\"",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT", Index: 0},
				{TokenType: token.StringLiteral, Literal: "Hello!", Index: 1},
				{TokenType: token.EOF, Literal: token.EOF, Index: 2},
			},
		},
		{
			Source: "Print \"So-called \"\"test\"\" this is\"",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT", Index: 0},
				{TokenType: token.StringLiteral, Literal: "So-called \"\"test\"\" this is", Index: 1},
				{TokenType: token.EOF, Literal: token.EOF, Index: 2},
			},
		},
		{
			Source: "Rem This is a comment",
			ExpectedTokens: []token.Token{
				{TokenType: token.REM, Literal: "REM", Index: 0},
				{TokenType: token.Comment, Literal: "This is a comment", Index: 1},
				{TokenType: token.EOF, Literal: token.EOF, Index: 2},
			},
		},
		{
			Source: "let x = 5",
			ExpectedTokens: []token.Token{
				{TokenType: token.LET, Literal: "LET", Index: 0},
				{TokenType: token.IdentifierLiteral, Literal: "X", Index: 1},
				{TokenType: token.Equal, Literal: "=", Index: 2},
				{TokenType: token.NumericLiteral, Literal: "5", Index: 3},
				{TokenType: token.EOF, Literal: token.EOF, Index: 4},
			},
		},
		{
			Source: "&554a3d2c",
			ExpectedTokens: []token.Token{
				{TokenType: token.HexLiteral, Literal: "&554A3D2C", Index: 0},
				{TokenType: token.EOF, Literal: token.EOF, Index: 1},
			},
		},
		{
			Source: "let y = &5d",
			ExpectedTokens: []token.Token{
				{TokenType: token.LET, Literal: "LET", Index: 0},
				{TokenType: token.IdentifierLiteral, Literal: "Y", Index: 1},
				{TokenType: token.Equal, Literal: "=", Index: 2},
				{TokenType: token.HexLiteral, Literal: "&5D", Index: 3},
				{TokenType: token.EOF, Literal: token.EOF, Index: 4},
			},
		},
		{
			Source: "rM_baSic_HAD_thiS_Weird_camel_case_tHING_GoInG_On$",
			ExpectedTokens: []token.Token{
				{TokenType: token.IdentifierLiteral, Literal: "Rm_Basic_Had_This_Weird_Camel_Case_Thing_Going_On$", Index: 0},
				{TokenType: token.EOF, Literal: token.EOF, Index: 1},
			},
		},
		{
			Source: "this% :=that$+ meh - 5.1234",
			ExpectedTokens: []token.Token{
				{TokenType: token.IdentifierLiteral, Literal: "This%", Index: 0},
				{TokenType: token.Assign, Literal: ":=", Index: 1},
				{TokenType: token.IdentifierLiteral, Literal: "That$", Index: 2},
				{TokenType: token.Plus, Literal: "+", Index: 3},
				{TokenType: token.IdentifierLiteral, Literal: "Meh", Index: 4},
				{TokenType: token.Minus, Literal: "-", Index: 5},
				{TokenType: token.NumericLiteral, Literal: "5.1234", Index: 6},
				{TokenType: token.EOF, Literal: token.EOF, Index: 7},
			},
		},
		{
			Source: "my_var% := yet_more_var% + foo_var",
			ExpectedTokens: []token.Token{
				{TokenType: token.IdentifierLiteral, Literal: "My_Var%", Index: 0},
				{TokenType: token.Assign, Literal: ":=", Index: 1},
				{TokenType: token.IdentifierLiteral, Literal: "Yet_More_Var%", Index: 2},
				{TokenType: token.Plus, Literal: "+", Index: 3},
				{TokenType: token.IdentifierLiteral, Literal: "Foo_Var", Index: 4},
				{TokenType: token.EOF, Literal: token.EOF, Index: 5},
			},
		},
		{
			Source: "50 > 50",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "50", Index: 0},
				{TokenType: token.GreaterThan, Literal: ">", Index: 1},
				{TokenType: token.NumericLiteral, Literal: "50", Index: 2},
				{TokenType: token.EOF, Literal: token.EOF, Index: 3},
			},
		},
		{
			Source: "50 <= 50",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "50", Index: 0},
				{TokenType: token.LessThanEqualTo1, Literal: "<=", Index: 1},
				{TokenType: token.NumericLiteral, Literal: "50", Index: 2},
				{TokenType: token.EOF, Literal: token.EOF, Index: 3},
			},
		},
		{
			Source: "50 >= 50",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "50", Index: 0},
				{TokenType: token.GreaterThanEqualTo1, Literal: ">=", Index: 1},
				{TokenType: token.NumericLiteral, Literal: "50", Index: 2},
				{TokenType: token.EOF, Literal: token.EOF, Index: 3},
			},
		},
		{
			Source: "-1 AND 1",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "-1", Index: 0},
				{TokenType: token.AND, Literal: "AND", Index: 1},
				{TokenType: token.NumericLiteral, Literal: "1", Index: 2},
				{TokenType: token.EOF, Literal: token.EOF, Index: 3},
			},
		},
		{
			Source: "-1.0",
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "-1.0", Index: 0},
				{TokenType: token.EOF, Literal: token.EOF, Index: 1},
			},
		},
		{
			Source: "NOT -1",
			ExpectedTokens: []token.Token{
				{TokenType: token.NOT, Literal: "NOT", Index: 0},
				{TokenType: token.NumericLiteral, Literal: "-1", Index: 1},
				{TokenType: token.EOF, Literal: token.EOF, Index: 2},
			},
		},
		{
			Source: "(-1)",
			ExpectedTokens: []token.Token{
				{TokenType: token.LeftParen, Literal: "(", Index: 0},
				{TokenType: token.NumericLiteral, Literal: "-1", Index: 1},
				{TokenType: token.RightParen, Literal: ")", Index: 2},
				{TokenType: token.EOF, Literal: token.EOF, Index: 3},
			},
		},
		{
			Source: "PRINT #1 !",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT", Index: 0},
				{TokenType: token.Hash, Literal: "#", Index: 1},
				{TokenType: token.NumericLiteral, Literal: "1", Index: 2},
				{TokenType: token.Exclamation, Literal: "!", Index: 3},
				{TokenType: token.EOF, Literal: token.EOF, Index: 4},
			},
		},
		{
			Source: "PRINT 1,2",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT", Index: 0},
				{TokenType: token.NumericLiteral, Literal: "1", Index: 1},
				{TokenType: token.Comma, Literal: ",", Index: 2},
				{TokenType: token.NumericLiteral, Literal: "2", Index: 3},
				{TokenType: token.EOF, Literal: token.EOF, Index: 4},
			},
		},
		{
			Source: "PRINT ~1 !",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT", Index: 0},
				{TokenType: token.Tilde, Literal: "~", Index: 1},
				{TokenType: token.NumericLiteral, Literal: "1", Index: 2},
				{TokenType: token.Exclamation, Literal: "!", Index: 3},
				{TokenType: token.EOF, Literal: token.EOF, Index: 4},
			},
		},
		{
			Source: "PRINT 2.34E+4",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT", Index: 0},
				{TokenType: token.NumericLiteral, Literal: "2.34E+4", Index: 1},
				{TokenType: token.EOF, Literal: token.EOF, Index: 2},
			},
		},
		{
			Source: "PRINT 1.344E-4.32",
			ExpectedTokens: []token.Token{
				{TokenType: token.PRINT, Literal: "PRINT", Index: 0},
				{TokenType: token.NumericLiteral, Literal: "1.344E-4.32", Index: 1},
				{TokenType: token.EOF, Literal: token.EOF, Index: 2},
			},
		},
		{
			Source: `LEN("HELLO")`,
			ExpectedTokens: []token.Token{
				{TokenType: token.IdentifierLiteral, Literal: "LEN", Index: 0},
				{TokenType: token.LeftParen, Literal: "(", Index: 1},
				{TokenType: token.StringLiteral, Literal: "HELLO", Index: 2},
				{TokenType: token.RightParen, Literal: ")", Index: 3},
				{TokenType: token.EOF, Literal: token.EOF, Index: 4},
			},
		},
		{
			Source: `set mode 80
cls : print "Blergh"
print "Hello"
`,
			ExpectedTokens: []token.Token{
				{TokenType: token.SET, Literal: "SET", Index: 0},
				{TokenType: token.MODE, Literal: "MODE", Index: 1},
				{TokenType: token.NumericLiteral, Literal: "80", Index: 2},
				{TokenType: token.NewLine, Literal: "\n", Index: 3},
				{TokenType: token.CLS, Literal: "CLS", Index: 4},
				{TokenType: token.Colon, Literal: ":", Index: 5},
				{TokenType: token.PRINT, Literal: "PRINT", Index: 6},
				{TokenType: token.StringLiteral, Literal: "Blergh", Index: 7},
				{TokenType: token.NewLine, Literal: "\n", Index: 8},
				{TokenType: token.PRINT, Literal: "PRINT", Index: 9},
				{TokenType: token.StringLiteral, Literal: "Hello", Index: 10},
				{TokenType: token.NewLine, Literal: "\n", Index: 11},
				{TokenType: token.EOF, Literal: "EOF", Index: 12},
			},
		},
		{
			Source: `10 set mode 80
20 cls : print "Blergh"
30 print "Hello"
`,
			ExpectedTokens: []token.Token{
				{TokenType: token.NumericLiteral, Literal: "10", Index: 0},
				{TokenType: token.SET, Literal: "SET", Index: 1},
				{TokenType: token.MODE, Literal: "MODE", Index: 2},
				{TokenType: token.NumericLiteral, Literal: "80", Index: 3},
				{TokenType: token.NewLine, Literal: "\n", Index: 4},
				{TokenType: token.NumericLiteral, Literal: "20", Index: 5},
				{TokenType: token.CLS, Literal: "CLS", Index: 6},
				{TokenType: token.Colon, Literal: ":", Index: 7},
				{TokenType: token.PRINT, Literal: "PRINT", Index: 8},
				{TokenType: token.StringLiteral, Literal: "Blergh", Index: 9},
				{TokenType: token.NewLine, Literal: "\n", Index: 10},
				{TokenType: token.NumericLiteral, Literal: "30", Index: 11},
				{TokenType: token.PRINT, Literal: "PRINT", Index: 12},
				{TokenType: token.StringLiteral, Literal: "Hello", Index: 13},
				{TokenType: token.NewLine, Literal: "\n", Index: 14},
				{TokenType: token.EOF, Literal: "EOF", Index: 15},
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
			if tok.Index != test.ExpectedTokens[j].Index {
				token.PrintToken(tok)
				t.Fatalf("Token [%d]: Index [%d] expected, got [%d] from source [%q]", i, test.ExpectedTokens[j].Index, tok.Index, test.Source)

			}
		}
	}
}
