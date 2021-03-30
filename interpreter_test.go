package main

import (
	"testing"
)

func TestInterpreterTokenize(t *testing.T) {

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
	}

	// test that we always get expected tokens
	interp := &Interpreter{}
	for i, test := range tests {
		interp.Init()
		interp.Tokenize(test.Source)
		for j, token := range interp.currentTokens {
			if token.TokenType != test.ExpectedTokens[j].TokenType {
				t.Fatalf("Token [%d]: TokenType [%d] expected, got [%d] from source [%q]", i, test.ExpectedTokens[j].TokenType, token.TokenType, test.Source)
			}
			if token.Literal != test.ExpectedTokens[j].Literal {
				t.Fatalf("Token [%d]: Literal [%q] expected, got [%q] from source [%q]", i, test.ExpectedTokens[j].Literal, token.Literal, test.Source)
			}
		}
	}
}

func TestInterpreterEvaluate(t *testing.T) {

	// test data
	type test struct {
		Source         string
		ExpectedResult float64
	}
	tests := []test{
		{
			Source:         "1+2",
			ExpectedResult: float64(3),
		},
		{
			Source:         "4-2",
			ExpectedResult: float64(2),
		},
		{
			Source:         "3+6.55",
			ExpectedResult: float64(9.55),
		},
		{
			Source:         "9*10",
			ExpectedResult: float64(90),
		},
		{
			Source:         "10        *  10",
			ExpectedResult: float64(100),
		},
		{
			Source:         "0.1 * 9",
			ExpectedResult: float64(0.9),
		},
		{
			Source:         "5 / 2",
			ExpectedResult: float64(2.5),
		},
		{
			Source:         "5 + 3 + 10",
			ExpectedResult: float64(18),
		},
		{
			Source:         "(2+4) * 10",
			ExpectedResult: float64(60),
		},
		{
			Source:         "2^10",
			ExpectedResult: float64(1024),
		},
		{
			Source:         "2^(5+5)",
			ExpectedResult: float64(1024),
		},
		{
			Source:         "6.3 \\ 2.2",
			ExpectedResult: float64(3),
		},
	}

	// test that we always get expected result
	interp := &Interpreter{}
	for _, test := range tests {
		interp.Init()
		interp.Tokenize(test.Source)
		_, _, _, result := interp.Evaluate(interp.currentTokens)
		if result != test.ExpectedResult {
			t.Fatalf("Expected [%f] but got [%f] from source [%q]", test.ExpectedResult, result, test.Source)
		}
	}
}

func TestFormatCode(t *testing.T) {

	// test data
	type test struct {
		Source              string
		HighlightTokenIndex int
		ExpectedCode        string
	}
	tests := []test{
		{
			Source:              "xpOs% := 542     + 3223  +    hello$",
			HighlightTokenIndex: -1,
			ExpectedCode:        "Xpos% := 542 + 3223 + Hello$",
		},
		{
			Source:              "xpOs% := 542     + 3223  +    hello$",
			HighlightTokenIndex: 2,
			ExpectedCode:        "Xpos% := --> 542 + 3223 + Hello$",
		},
		{
			Source:              "xpOs% := 542     + 3223  +    hello$",
			HighlightTokenIndex: -1,
			ExpectedCode:        "Xpos% := 542 + 3223 + Hello$",
		},
		{
			Source:              "xpOs% := 542     + 3223  +    hello$",
			HighlightTokenIndex: 0,
			ExpectedCode:        "--> Xpos% := 542 + 3223 + Hello$",
		},
	}
	// test that we always get expected result
	interp := &Interpreter{}
	for _, test := range tests {
		interp.Init()
		formattedCode := interp.FormatCode(test.Source, test.HighlightTokenIndex)
		if formattedCode != test.ExpectedCode {
			t.Fatalf("Expected [%s] but got [%s]", test.ExpectedCode, formattedCode)
		}
	}
}

func TestEvaluateErrorHandling(t *testing.T) {

	// test data
	type test struct {
		Source            string
		ExpectedErrorCode int
	}
	tests := []test{
		{
			Source:            "foo = bar + 2",
			ExpectedErrorCode: HasNotBeenDefined,
		},
		{
			Source:            "foo = foo + 2",
			ExpectedErrorCode: HasNotBeenDefined,
		},
	}
	// test that we always get expected result
	interp := &Interpreter{}
	for _, test := range tests {
		interp.Init()
		errorCode, _, _ := interp.RunLine(test.Source)
		if errorCode != test.ExpectedErrorCode {
			t.Fatalf("Expected errorCode %d (%s) but got %d (%s)", test.ExpectedErrorCode, errorMessage(test.ExpectedErrorCode), errorCode, errorMessage(errorCode))
		}
	}
}

func TestInterpreterVariableAssignment(t *testing.T) {

	// test data
	type test struct {
		Source        string
		ExpectedName  string
		ExpectedValue float64
	}
	tests := []test{
		{
			Source:        "one = 1",
			ExpectedName:  "One",
			ExpectedValue: float64(1),
		},
		{
			Source:        "two = 1+1",
			ExpectedName:  "Two",
			ExpectedValue: float64(2),
		},
		{
			Source:        "two% := 1+ 1",
			ExpectedName:  "Two%",
			ExpectedValue: float64(2),
		},
		{
			Source:        "x := 1.2 + 0.5",
			ExpectedName:  "X",
			ExpectedValue: float64(1.7),
		},
		{
			Source:        "x% := 1.6",
			ExpectedName:  "X%",
			ExpectedValue: float64(2),
		},
		{
			Source:        "x% := 1.2 + 0.5",
			ExpectedName:  "X%",
			ExpectedValue: float64(2),
		},
	}

	// test that we always get expected result
	interp := &Interpreter{}
	for _, test := range tests {
		interp.Init()
		interp.RunLine(test.Source)
		// Can variable be found?
		if _, ok := interp.store[test.ExpectedName]; ok {
			valfloat64, ok := interp.store[test.ExpectedName].(float64)
			// Can the value be parsed?
			if !ok {
				t.Fatalf("Could not interpret stored value for [%q] as a number", test.ExpectedName)
			} else {
				// Is the value correct?
				if valfloat64 != test.ExpectedValue {
					t.Fatalf("Expected [%f] but got [%f] for [%q]", test.ExpectedValue, valfloat64, test.ExpectedName)
				}
			}
		} else {
			t.Fatalf("Did not find [%q] in the store", test.ExpectedName)
		}
	}
}
