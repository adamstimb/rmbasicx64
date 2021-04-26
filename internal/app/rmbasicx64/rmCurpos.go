package rmbasicx64

import (
	"math"
)

// RmSetCurpos represents the SET CURPOS command
func (i *Interpreter) RmSetCurpos() (ok bool) {
	i.TokenPointer += 2
	// Get the col
	val, ok := i.OnExpression("numeric")
	if !ok {
		return false
	}
	// TODO: ensure within range --> have a function to do that
	col := int(math.Round(val.(float64)))
	// Get the comma
	if !i.OnComma() {
		return false
	}
	// Get the row
	val, ok = i.OnExpression("numeric")
	if !ok {
		return false
	}
	// TODO: ensure within range --> have a function to do that
	row := int(math.Round(val.(float64)))
	// No more parameters
	if !i.OnSegmentEnd() {
		return false
	}
	// Execute
	i.g.SetCurpos(col, row)
	return true
}
