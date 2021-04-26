package rmbasicx64

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
func (i *Interpreter) RmList() (ok bool) {
	i.TokenPointer++
	line1 := -1
	line2 := -1
	if !i.OnSegmentEnd() {
		// Get either "TO line2" or "line1"
		if i.OnTo() {
			// Get TO, must now get line 2
			val, ok := i.OnExpression("numeric")
			if !ok {
				return false
			}
			line2 = int(math.Round(val.(float64)))
		} else {
			// Didn't get TO, must now get line1
			val, ok := i.OnExpression("numeric")
			if !ok {
				return false
			}
			line1 = int(math.Round(val.(float64)))
			// If there're more parameters we need a TO
			if !i.OnSegmentEnd() {
				if !i.OnTo() {
					return false
				}
				// If there're more parameters it has to be line2
				if !i.OnSegmentEnd() {
					val, ok := i.OnExpression("numeric")
					if !ok {
						return false
					}
					line2 = int(math.Round(val.(float64)))
				}
			}
		}
	}
	// No more params
	if !i.OnSegmentEnd() {
		return false
	}
	// Execute
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
			i.g.Put(13)
		}
		if inRange && lineNumber == line2 {
			inRange = false
		}
	}
	return true
}
