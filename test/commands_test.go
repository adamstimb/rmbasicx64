package main

import (
	"testing"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

func TestRun(t *testing.T) {

	// test data
	type test struct {
		Program         map[int]string
		ExpectedTestVal float64
	}
	tests := []test{
		{
			Program: map[int]string{
				10: "Test := 1",
			},
			ExpectedTestVal: float64(1),
		},
		{
			Program: map[int]string{
				10: "Test := 1",
				20: "Test := 2",
			},
			ExpectedTestVal: float64(2),
		},
		{
			Program: map[int]string{
				10: "Test := 1",
				20: "Test := 2",
				5:  "Test := 5",
			},
			ExpectedTestVal: float64(2),
		},
	}

	interp := &rmbasicx64.Interpreter{}
	for _, test := range tests {
		interp.Init(&rmbasicx64.Game{})
		interp.Program = test.Program
		if !interp.RmRun() {
			t.Fatalf("%s\n", syntaxerror.ErrorMessage(interp.ErrorCode))
		}
		if _, ok := interp.Store["Test"]; ok {
			valfloat64, ok := interp.Store["Test"].(float64)
			if !ok {
				t.Fatalf("Could not extract value of variable Test from store")
			} else {
				if valfloat64 != test.ExpectedTestVal {
					t.Fatalf("Expected %f for value of variable Test but got %f instead", test.ExpectedTestVal, valfloat64)
				}
			}
		} else {
			t.Fatalf("Could not find variable Test in store")
		}
	}
}

func TestGoto(t *testing.T) {

	// test data
	type test struct {
		Program         map[int]string
		ExpectedTestVal float64
	}
	tests := []test{
		{
			Program: map[int]string{
				10: "Test := 1",
				20: "GOTO 30",
				30: "PRINT \"Hello\"",
				40: "Test := 100",
				50: "PRINT \"world\"",
			},
			ExpectedTestVal: float64(100),
		},
		{
			Program: map[int]string{
				10: "Test := 1",
				20: "GOTO 50",
				30: "PRINT \"Hello\"",
				40: "Test := 100",
				50: "PRINT \"world\"",
			},
			ExpectedTestVal: float64(1),
		},
		{
			Program: map[int]string{
				10: "Test := 50",
				20: "GOTO Test",
				30: "PRINT \"Hello\"",
				40: "Test := 100",
				50: "PRINT \"world\"",
			},
			ExpectedTestVal: float64(50),
		},
		{
			Program: map[int]string{
				10: "Test := 40",
				20: "GOTO Test + 10",
				30: "PRINT \"Hello\"",
				40: "Test := 100",
				50: "PRINT \"world\"",
			},
			ExpectedTestVal: float64(40),
		},
	}

	interp := &rmbasicx64.Interpreter{}
	for _, test := range tests {
		game := rmbasicx64.NewGame()
		interp.Init(game)
		interp.Program = test.Program
		if !interp.RmRun() {
			t.Fatalf("%s\n", syntaxerror.ErrorMessage(interp.ErrorCode))
		}
		if _, ok := interp.Store["Test"]; ok {
			valfloat64, ok := interp.Store["Test"].(float64)
			if !ok {
				t.Fatalf("Could not extract value of variable Test from store")
			} else {
				if valfloat64 != test.ExpectedTestVal {
					t.Fatalf("Expected %f for value of variable Test but got %f instead", test.ExpectedTestVal, valfloat64)
				}
			}
		} else {
			t.Fatalf("Could not find variable Test in store")
		}
	}
}