package rmbasicx64

import (
	"fmt"
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
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
func (i *Interpreter) RmList() (ok bool) {
	i.TokenPointer++
	line1 := -1
	line2 := -1
	if !i.EndOfTokens() {
		// Get line1 and line2
		if i.IsAnyOfTheseTokens([]int{token.NumericalLiteral, token.IdentifierLiteral, token.TO}) {
			switch i.TokenStack[i.TokenPointer].TokenType {
			case token.TO:
				i.TokenPointer++
				if i.EndOfTokens() {
					i.ErrorCode = syntaxerror.LineNumberExpected
					i.Message = syntaxerror.ErrorMessage(syntaxerror.LineNumberExpected)
					i.BadTokenIndex = i.TokenPointer
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
					_, ok := i.AcceptAnyOfTheseTokens([]int{token.TO})
					if ok {
						val, ok := i.AcceptAnyNumber()
						if ok {
							// get line2 and accept no further parameters
							line2 = int(math.Round(val))
							if !i.EndOfTokens() {
								i.ErrorCode = syntaxerror.EndOfInstructionExpected
								i.Message = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
								i.BadTokenIndex = i.TokenPointer
								return false
							}
						} else {
							return false
						}
					} else {
						i.ErrorCode = syntaxerror.EndOfInstructionExpected
						i.Message = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
						i.BadTokenIndex = i.TokenPointer - 1
						return false
					}
				}
			}
		} else {
			i.ErrorCode = syntaxerror.EndOfInstructionExpected
			i.Message = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
			i.BadTokenIndex = i.TokenPointer
			return false
		}
	}
	// Pass through if no program
	if len(i.Program) == 0 {
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
			i.g.Print(fmt.Sprintf("%d %s", lineNumber, i.Program[lineNumber]))
		}
		if inRange && lineNumber == line2 {
			inRange = false
		}
	}
	return true
}
