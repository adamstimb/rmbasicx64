package main

import (
	"testing"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64"
)

func TestInterpreterEvaluateExpression(t *testing.T) {

	// test data
	type test struct {
		Source         string
		ExpectedResult interface{}
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
		{
			Source:         "100 = 100",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "100.00 == 100",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "100 = 10",
			ExpectedResult: float64(0),
		},
		{
			Source:         "10 < 100",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "100 > 10",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "100 > 1000",
			ExpectedResult: float64(0),
		},
		{
			Source:         "10000 < 10",
			ExpectedResult: float64(0),
		},
		{
			Source:         "100 <= 100",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "100 >= 100",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "50 <> 50",
			ExpectedResult: float64(0),
		},
		{
			Source:         "50 <> 55",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "\"Hey\"  + \" \"+ \"you\"",
			ExpectedResult: "Hey you",
		},
		{
			Source:         "\"Hey\"  + \" \"+ \"You \"   +\"Guys!\"",
			ExpectedResult: "Hey You Guys!",
		},
		{
			Source:         "\"Screaming\" + \"Lord\" + \"Sutch\"",
			ExpectedResult: "ScreamingLordSutch",
		},
		{
			Source:         "\"Front\" + 242",
			ExpectedResult: "Front242",
		},
		// Test some real examples from the original RM Basic book why not:
		{
			Source:         "\"Freda\" > \"Fred\"",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "\"banana\" > \"BANANA\"",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "\"Class A\" > \"Class 1\"",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "\"banana\" == \"BANANA\"",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "\"banana\" = \"BANANA\"",
			ExpectedResult: float64(0),
		},
		{
			Source:         "4 AND 2",
			ExpectedResult: float64(0),
		},
		{
			Source:         "-1.0 AND -1.0",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "0 OR -1",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "NOT -1",
			ExpectedResult: float64(0),
		},
		{
			Source:         "NOT 0",
			ExpectedResult: float64(-1),
		},
		{
			Source:         "1 + (-1 AND -1)",
			ExpectedResult: float64(0),
		},
		{
			Source:         "not (-1 AND -1)",
			ExpectedResult: float64(0),
		},
	}

	// test that we always get expected result
	interp := &rmbasicx64.Interpreter{}
	for _, test := range tests {
		game := rmbasicx64.NewGame()
		interp.Init(game)
		interp.Tokenize(test.Source)
		interp.TokenStack = interp.CurrentTokens
		interp.TokenPointer = 0
		result, _ := interp.EvaluateExpression()
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
	interp := &rmbasicx64.Interpreter{}
	for _, test := range tests {
		game := rmbasicx64.NewGame()
		interp.Init(game)
		formattedCode := interp.FormatCode(test.Source, test.HighlightTokenIndex, false)
		if formattedCode != test.ExpectedCode {
			t.Fatalf("Expected [%s] but got [%s]", test.ExpectedCode, formattedCode)
		}
	}
}

//func TestEvaluateErrorHandling(t *testing.T) {
//
//	// test data
//	type test struct {
//		Source            string
//		ExpectedErrorCode int
//	}
//	tests := []test{
//		{
//			Source:            "foo = bar + 2",
//			ExpectedErrorCode: HasNotBeenDefined,
//		},
//		{
//			Source:            "foo = foo + 2",
//			ExpectedErrorCode: HasNotBeenDefined,
//		},
//		{
//			Source:            "foo = \"foo\" * \"bar\"",
//			ExpectedErrorCode: InvalidExpression,
//		},
//	}
//	// test that we always get expected result
//	interp := &Interpreter{}
//	for _, test := range tests {
//		interp.Init()
//		_ = interp.RunLine(test.Source)
//		if interp.ErrorCode != test.ExpectedErrorCode {
//			t.Fatalf("Expected errorCode %d (%s) but got %d (%s)", test.ExpectedErrorCode, syntaxerror.ErrorMessage(test.ExpectedErrorCode), interp.ErrorCode, syntaxerror.ErrorMessage(interp.ErrorCode))
//		}
//	}
//}
//
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
	interp := &rmbasicx64.Interpreter{}
	for _, test := range tests {
		game := rmbasicx64.NewGame()
		interp.Init(game)
		interp.RunLine(test.Source)
		// Can variable be found?
		if _, ok := interp.Store[test.ExpectedName]; ok {
			valfloat64, ok := interp.Store[test.ExpectedName].(float64)
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

func TestInterpreterWeighString(t *testing.T) {
	w := rmbasicx64.WeighString("Ohhhh yeah")
	expected := 79 + (4 * 104) + 32 + 121 + 101 + 97 + 104
	if w != expected {
		t.Fatalf("Expected [%d] but got [%d]", expected, w)
	}
}

//func TestImmediateInput(t *testing.T) {
//
//	// test data
//	type test struct {
//		Source          string
//		ExpectedProgram map[int]string
//	}
//	tests := []test{
//		{
//			Source: "10 set  mode  40",
//			ExpectedProgram: map[int]string{
//				10: "SET MODE 40",
//			},
//		},
//		{
//			Source: "20 print \"Just testing\"",
//			ExpectedProgram: map[int]string{
//				10: "SET MODE 40",
//				20: "PRINT \"Just testing\"",
//			},
//		},
//		{
//			Source: "5 cls",
//			ExpectedProgram: map[int]string{
//				5:  "CLS",
//				10: "SET MODE 40",
//				20: "PRINT \"Just testing\"",
//			},
//		},
//		{
//			Source: "run",
//			ExpectedProgram: map[int]string{
//				5:  "CLS",
//				10: "SET MODE 40",
//				20: "PRINT \"Just testing\"",
//			},
//		},
//	}
//
//	// This test simulates a user manually keying in a program
//	interp := &rmbasicx64.Interpreter{}
//	game := rmbasicx64.NewGame()
//	interp.Init(game)
//	for _, test := range tests {
//		interp.ImmediateInput(test.Source)
//		for lineNumber, expectedCode := range test.ExpectedProgram {
//			actualCode, ok := interp.Program[lineNumber]
//			if !ok {
//				t.Fatalf("Could not find line %d in program", lineNumber)
//			}
//			if actualCode != expectedCode {
//				t.Fatalf("Expected [%s] but got [%s] in line %d", expectedCode, actualCode, lineNumber)
//			}
//		}
//	}
//}
