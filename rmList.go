package main

import (
	"fmt"
	"math"
)

// rmList represents the LIST command
// From the manual:
// LIST [~e/#e] [range]
// This command lists the specified line or range of lines from the program
// currently in memory to either the specified writing area or to the
// specified channel (0 for screen, 2 for printer, 11-127 for file).
// LIST - lists the whole program
// LIST line1 - lists line1
// LIST line1 TO - lists all lines from line1 onwards
// LIST line1 TO line2 - lists all lines between
// LIST TO line2 - lists all lines up to line2
func (i *Interpreter) rmList() (ok bool) {
	i.tokenPointer++
	line1 := -1
	line2 := -1
	if !i.EndOfTokens() {
		// Get line1 and line2
		if i.IsAnyOfTheseTokens([]int{NumericalLiteral, IdentifierLiteral, TO}) {
			switch i.tokenStack[i.tokenPointer].TokenType {
			case TO:
				i.tokenPointer++
				if i.EndOfTokens() {
					i.errorCode = LineNumberExpected
					i.message = errorMessage(LineNumberExpected)
					i.badTokenIndex = i.tokenPointer
					return false
				}
				val, ok := i.AcceptAnyNumber()
				if ok {
					// round val and store as line2
					line2 = int(math.Round(val))
				} else {
					return false
				}
			default:
				val, ok := i.AcceptAnyNumber()
				if ok {
					// round val and store as line1 and set line2 to the same value in case
					// user wants to list only one line
					line1 = int(math.Round(val))
					line2 = line1
				} else {
					return false
				}
				if !i.EndOfTokens() {
					// Must be followed by TO
					_, ok := i.AcceptAnyOfTheseTokens([]int{TO})
					if ok {
						val, ok := i.AcceptAnyNumber()
						if ok {
							// get line2 and accept no further parameters
							line2 = int(math.Round(val))
							if !i.EndOfTokens() {
								i.errorCode = EndOfInstructionExpected
								i.message = errorMessage(EndOfInstructionExpected)
								i.badTokenIndex = i.tokenPointer
								return false
							}
						} else {
							return false
						}
					} else {
						i.errorCode = EndOfInstructionExpected
						i.message = errorMessage(EndOfInstructionExpected)
						i.badTokenIndex = i.tokenPointer - 1
						return false
					}
				}
			}
		} else {
			i.errorCode = EndOfInstructionExpected
			i.message = errorMessage(EndOfInstructionExpected)
			i.badTokenIndex = i.tokenPointer
			return false
		}
	}
	// Pass through if no program
	if len(i.program) == 0 {
		return true
	}
	// Set default line1, line2
	lineOrder := i.GetLineOrder()
	if line1 < 0 {
		line1 = lineOrder[0]
	}
	if line2 < 0 {
		line2 = lineOrder[len(lineOrder)-1]
	}
	// Print the listing within range
	inRange := false
	for _, lineNumber := range lineOrder {
		if !inRange && lineNumber == line1 {
			inRange = true
		}
		if inRange {
			i.g.Print(fmt.Sprintf("%d %s", lineNumber, i.program[lineNumber]))
		}
		if inRange && lineNumber == line2 {
			inRange = false
		}
	}
	return true
}
