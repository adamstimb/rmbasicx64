package rmbasicx64

import (
	"math"
)

// RmSetCursor represents the SET CURSOR command
// SET CURSOR e1[,e2[,e3][,e4]]]
// e1 === cursorMode, e2 === cursorChar, e3 === cursorCharSet
func (i *Interpreter) RmSetCursor() (ok bool) {
	paramCount := 0
	i.TokenPointer += 2
	// Get required cursorMode
	val, ok := i.OnExpression("numeric")
	if !ok {
		return false
	}
	// TODO: ensure within range --> have a function to do that
	cursorMode := int(math.Round(val.(float64)))
	var cursorChar int
	var cursorCharSet int
	paramCount++
	// Get optional cursorChar
	if !i.OnSegmentEnd() {
		// Get comma
		if !i.OnComma() {
			return false
		}
		// Get cursorChar
		val, ok = i.OnExpression("numeric")
		if !ok {
			return false
		} else {
			cursorChar = int(math.Round(val.(float64)))
			paramCount++
		}
	}
	// Get optional cursorCharSet
	if !i.OnSegmentEnd() {
		// Get comma
		if !i.OnComma() {
			return false
		}
		// Get cursorChar
		val, ok = i.OnExpression("numeric")
		if !ok {
			return false
		} else {
			cursorCharSet = int(math.Round(val.(float64)))
			paramCount++
		}
	}
	// No more params
	if !i.OnSegmentEnd() {
		return false
	}
	// Execute
	switch paramCount {
	case 1:
		i.g.SetCursor(cursorMode)
	case 2:
		i.g.SetCursor(cursorMode, cursorChar)
	case 3:
		i.g.SetCursor(cursorMode, cursorChar, cursorCharSet)
	}
	return true
}
