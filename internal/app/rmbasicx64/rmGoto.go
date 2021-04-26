package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// rmGoto represents the GOTO command
func (i *Interpreter) RmGoto() (ok bool) {
	i.TokenPointer++
	// Get required gotoLine
	val, ok := i.OnExpression("numeric")
	if !ok {
		return false
	}
	gotoLine := int(math.Round(val.(float64)))
	// Execute
	// scan for line number to goto
	lineOrder := i.GetLineOrder()
	for l := 0; l < len(lineOrder); l++ {
		if lineOrder[l] == gotoLine {
			// found the line so set pointer to one behind because RunLine advances it and return
			i.ProgramPointer = l - 1
			return true
		}
	}
	// line does not exist
	i.ErrorCode = syntaxerror.SpecifiedLineNotFound
	i.TokenPointer = 1
	return false
}
