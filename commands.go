package main

import "fmt"

// rmRun represents the RUN command
func (i *Interpreter) rmRun() (ok bool) {
	for lineNumber, code := range i.program {
		i.lineNumber = lineNumber
		isOk := i.RunLine(code)
		if !isOk {
			return false
		}
	}
	return true
}

// rmPrint represents the Print command
func (i *Interpreter) rmPrint(tokens []Token) (ok bool) {
	// PRINT with no args
	if len(tokens) == 1 {
		fmt.Println("")
		return true
	}
	if len(tokens) > 1 {
		if tokens[1].TokenType == EndOfLine {
			// Still PRINT with no args
			fmt.Println("")
			return true
		}
		if tokens[1].TokenType == StringLiteral {
			// PRINT "hello"
			fmt.Println(tokens[1].Literal)
			return true
		}
	}
	// set error status here
	return false
}
