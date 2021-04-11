package main

import "fmt"

// rmPrint represents the Print command
// From the manual:
// PRINT [~e/#e,][print list]
// PRINT (or ? abbreviation) outputs data to a writing area or channel
// (0 for screen, 2 for printer, 11-127 for an open file).  When #e is
// omitted, output goes to the current writing area.
// [print list] is a list of numeric and/or string expressions. The
// expressions should be separated by a semicolon, comma, space or !
// A number is written according to the following rules:
// - a positive number (and zero) has a preceding space
// - a negative number has a preceding minus sign
// - numbers > 9999999 or < 0.001 are written in scientific notation,
//    where a number is followed by the letter e and a signed integer
//    (the exponent), e.g. 2.2e+06
// - each number has a space following it
// String expressions must be enclosed in " " marks. To output a " mark
// place 2 "" together, e.g. PRINT "The ""RM Nimbus"""
// ; does not introduce a space, e.g. PRINT 2;2 --> 22
// , puts following value in next available print zone. Print zones are
// 15 columns wide
// ! puts the following value on a new line
// If ; or , terminates [print list] then PRINT does not line feed
// and applies the above spacing rules
// See also WIDTH command.....
func (i *Interpreter) rmPrint() (ok bool) {
	// PRINT with no args
	if len(i.tokenStack) == 1 {
		fmt.Println("")
		return true
	}
	if len(i.tokenStack) > 1 {
		if i.tokenStack[1].TokenType == EndOfLine {
			// Also PRINT with no args
			fmt.Println("")
			return true
		}
		if i.tokenStack[1].TokenType == StringLiteral || i.tokenStack[1].TokenType == IdentifierLiteral || i.tokenStack[1].TokenType == NumericalLiteral {
			i.tokenPointer++
			toPrint, ok := i.EvaluateExpression(i.ExtractExpression())
			if !ok {
				i.badTokenIndex = 1
				return false
			} else {
				switch GetType(toPrint) {
				case "string":
					fmt.Println(toPrint.(string))
					return true
				case "float64":
					fmt.Println(toPrint.(float64))
					return true
				}
			}
		}
	}
	// set error status here
	return false
}
