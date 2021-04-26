package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// rmUpdateLine represents a line number being updated in the REPL
func (i *Interpreter) RmUpdateLine(code string) (ok bool) {
	// If the code begins with a round number then that's a line number
	if i.CurrentTokens[0].TokenType == token.NumericalLiteral {
		val, ok := i.GetValueFromToken(i.CurrentTokens[0], "float64")
		if ok {
			lineNumber := val.(float64)
			if lineNumber == math.Round(lineNumber) {
				// Is a line number so ...
				// If there are no more tokens delete the line if it exists
				// otherwise update it
				if i.CurrentTokens[1].TokenType == token.EndOfLine {
					if _, ok := i.Program[int(lineNumber)]; ok {
						delete(i.Program, int(lineNumber))
						return true
					} else {
						// line to delete does not exist so pass through
						return true
					}
				} else {
					i.Program[int(lineNumber)] = i.FormatCode(code, -1, true)
					return true
				}
			} else {
				i.ErrorCode = syntaxerror.LineNumberExpected
				i.TokenPointer = 0
				return false
			}
		}
	}
	i.ErrorCode = syntaxerror.LineNumberExpected
	i.TokenPointer = 0
	return false
}
