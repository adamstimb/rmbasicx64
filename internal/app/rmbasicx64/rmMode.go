package rmbasicx64

import (
	"math"
)

// RmSetMode represents the SET MODE command
func (i *Interpreter) RmSetMode() (ok bool) {
	i.TokenPointer += 2
	// Get required mode
	val, ok := i.OnExpression("numeric")
	if !ok {
		return false
	}
	// TODO: ensure within range --> have a function to do that
	mode := int(math.Round(val.(float64)))
	// Ensure no more parameters
	if !i.OnSegmentEnd() {
		return false
	}
	// Execute
	i.g.SetMode(mode)
	return true
}
