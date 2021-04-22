package rmbasicx64

import (
	"fmt"
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// RmEdit represents the EDIT command
func (i *Interpreter) RmEdit() (ok bool) {
	i.TokenPointer++
	if i.EndOfTokens() {
		// No linenumber passed so get line of last error
		// TODO: implement
		return true
	}
	_, ok = i.AcceptAnyOfTheseTokens([]int{token.NumericalLiteral, token.IdentifierLiteral})
	if ok {
		i.TokenPointer--
		lineNumber, ok := i.AcceptAnyNumber()
		if ok {
			if lineNumber == math.Round(lineNumber) {
				if _, ok := i.Program[int(lineNumber)]; ok {
					// edit existing line
					edited := i.g.Input(fmt.Sprintf("%d %s", int(lineNumber), i.Program[int(lineNumber)]))
					_ = i.ImmediateInput(edited)
					return true
				} else {
					// create new line
					edited := i.g.Input(fmt.Sprintf("%d %s", int(lineNumber), ""))
					_ = i.ImmediateInput(edited)
					return true
				}
			} else {
				i.ErrorCode = syntaxerror.LineNumberExpected
				return false
			}
		}
	}
	return false
}
