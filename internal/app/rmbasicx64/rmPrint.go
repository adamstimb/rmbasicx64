package rmbasicx64

import (
	"fmt"
	"math"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

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
func (i *Interpreter) RmPrint() (ok bool) {
	i.TokenPointer++
	if i.OnSegmentEnd() {
		// PRINT with no args
		i.g.Print("")
		i.g.Put(13)
		return true
	}
	var printLines []string
	writingArea := 0
	printChannel := 0
	// Optional tilde or hash
	t, ok := i.OnToken([]int{token.Tilde, token.Hash})
	if ok {
		// Got tilde or hash so must now get a numeric value
		val, ok := i.OnExpression("numeric")
		if !ok {
			return false
		} else {
			switch t.TokenType {
			case token.Tilde:
				// select writing area
				writingArea = int(math.Round(val.(float64)))
			case token.Hash:
				// select print channel
				printChannel = int(math.Round(val.(float64)))
			}
		}
	}
	// Handle no further args
	if i.OnSegmentEnd() {
		i.g.Print("")
		i.g.Put(13)
		return true
	}
	_ = writingArea // TODO: implement
	_ = printChannel
	noCarriageReturn := false // This flag will be set to true if the final expression is ;
	// Handle expression list
	// Evaluate all expressions in list, concatenate their results and send to print
	printString := ""
	for !i.OnSegmentEnd() {
		// Handle list punctuation
		// ;
		if i.OnSemicolon() {
			// no space
			if i.OnSegmentEnd() {
				noCarriageReturn = true
			}
			continue
		}
		// ,
		if i.OnComma() {
			// should jump to next print zone but for now add 15 char space
			printString += "               "
			continue
		}
		// !
		_, ok = i.OnToken([]int{token.Exclamation})
		if ok {
			// add new line
			printLines = append(printLines, printString)
			printString = ""
			continue
		}
		// expression
		toPrint, ok := i.OnExpression("any")
		if !ok {
			return false
		} else {
			switch HasType(toPrint) {
			case "string":
				printString += toPrint.(string)
				continue
			case "numeric":
				printString += RenderNumberAsString(toPrint.(float64))
				continue
			}
		}
	}
	// Got all expressions so print the final string(s) with a cr between each line
	printLines = append(printLines, printString)
	for index, toPrint := range printLines {
		i.g.Print(toPrint)
		if index < len(printLines)-1 {
			i.g.Put(13)
		}
	}
	if !noCarriageReturn {
		i.g.Put(13)
	}
	return true
}

// RenderNumberAsString receives a float64 type ad applies RM Basic's print rules for numbers
// returning a string representing the number that was passed.
func RenderNumberAsString(value float64) (result string) {
	if math.Abs(value) > 9999999 || math.Abs(value) < 0.001 {
		// use scientific notation to 5 decimal places
		result = fmt.Sprintf("%.5e", value)
	} else {
		result = fmt.Sprintf("%f", value)
	}
	// Special case of 0
	if value == 0 {
		return "0"
	}
	// remove trailing zeros
	return strings.TrimRight(strings.TrimRight(result, "0"), ".")
}
