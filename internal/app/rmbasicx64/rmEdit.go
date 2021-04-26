package rmbasicx64

import (
	"fmt"
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// RmEdit represents the EDIT command
func (i *Interpreter) RmEdit() (ok bool) {
	i.TokenPointer++
	if i.OnSegmentEnd() {
		// No linenumber passed so get line of last error
		// TODO: implement
		return true
	}
	// Get optional line number
	val, ok := i.OnExpression("numeric")
	if !ok {
		i.ErrorCode = syntaxerror.LineNumberExpected
		return false
	}
	lineNumber := int(math.Round(val.(float64)))
	// Execute
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
}
