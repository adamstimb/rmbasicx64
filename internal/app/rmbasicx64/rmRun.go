package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// rmRun represents the RUN command
func (i *Interpreter) RmRun() (ok bool) {
	i.TokenPointer++
	i.GetData() // Scan for DATA statements and programData stack
	i.ProgramPointer = 0
	lineOrder := i.GetLineOrder()
	// Check for optional startFrom parameter
	startFrom := lineOrder[i.ProgramPointer]
	if !i.OnSegmentEnd() {
		// Get startFrom
		val, ok := i.OnExpression("numeric")
		if !ok {
			return false
		}
		startFrom = int(math.Round(val.(float64)))
	}
	// No more params
	if !i.OnSegmentEnd() {
		return false
	}
	// Execute
	// pass thru if no program
	if len(i.Program) == 0 {
		return true
	}
	// Run the program and if startFrom was passed scan ahead to it before executing
	scanAhead := true
	for i.ProgramPointer < len(lineOrder) {
		i.LineNumber = lineOrder[i.ProgramPointer]
		if scanAhead && (i.LineNumber != startFrom) {
			i.ProgramPointer++
			continue
		}
		if scanAhead && (i.LineNumber == startFrom) {
			scanAhead = false
		}
		ok := i.RunLine(i.Program[lineOrder[i.ProgramPointer]])
		if !ok {
			i.ProgramPointer = 0
			return false
		}
	}
	// Catch line number not found
	if scanAhead {
		i.ErrorCode = syntaxerror.SpecifiedLineNotFound
		i.TokenPointer = 1
		i.LineNumber = -1
		return false
	}
	// Otherwise ok
	i.ProgramPointer = 0
	return true
}
