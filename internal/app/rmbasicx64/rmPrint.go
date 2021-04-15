package rmbasicx64

import (
	"fmt"
	"math"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
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
	// PRINT with no args
	if len(i.TokenStack) == 1 {
		//fmt.Println("")
		i.g.Print("")
		return true
	}
	if len(i.TokenStack) > 1 {
		i.TokenPointer++
		if i.TokenStack[i.TokenPointer].TokenType == token.EndOfLine {
			// Also PRINT with no args
			// fmt.Println("")
			i.g.Print("")
			return true
		}
		// Handle channel number or writing area option
		tildeOrHash, ok := i.AcceptAnyOfTheseTokens([]int{token.Tilde, token.Hash})
		writingArea := 0
		printChannel := 0
		if ok {
			var optionVal float64
			optionVal, ok = i.AcceptAnyNumber()
			if !ok {
				// excepted a number
				return false
			} else {
				switch tildeOrHash.TokenType {
				case token.Tilde:
					// select writing area
					writingArea = int(math.Round(optionVal))
				case token.Hash:
					// select print channel
					printChannel = int(math.Round(optionVal))
				}
			}
			// Handle no further args
			if i.EndOfTokens() {
				//fmt.Println("")
				i.g.Print("")
				return true
			}
		}
		_ = writingArea // TODO: implement
		_ = printChannel
		// Handle expression list
		if i.IsAnyOfTheseTokens([]int{token.StringLiteral, token.IdentifierLiteral, token.NumericalLiteral, token.Exclamation}) {
			// Evaluate all expressions in list, concatenate their results and send to print
			printString := ""
			for !i.EndOfTokens() {
				// Handle list punctuation
				switch i.TokenStack[i.TokenPointer].TokenType {
				case token.Semicolon:
					// no space, go to next token
					i.TokenPointer++
				case token.Comma:
					// should jump to next print zone but for now add a tab
					printString += "\t"
					i.TokenPointer++
				case token.Exclamation:
					// add new line
					printString += "\n"
					i.TokenPointer++
				default:
					toPrint, ok := i.EvaluateExpression()
					if !ok {
						//i.BadTokenIndex++
						return false
					} else {
						switch GetType(toPrint) {
						case "string":
							printString += toPrint.(string)
							continue
						case "float64":
							printString += RenderNumberAsString(toPrint.(float64))
							continue
						}
					}
				}
			}
			// Got all expressions so print the final string
			//fmt.Println(printString)
			i.g.Print(printString)
			return true
		} else {
			i.ErrorCode = syntaxerror.EndOfInstructionExpected
			i.Message = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
			i.BadTokenIndex = i.TokenPointer
			return false
		}
	}
	i.ErrorCode = syntaxerror.EndOfInstructionExpected
	i.Message = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
	i.BadTokenIndex = i.TokenPointer
	return false
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
	// remove trailing zeros
	return strings.TrimRight(strings.TrimRight(result, "0"), ".")
}
