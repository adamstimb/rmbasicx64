package rmbasicx64

import (
	"math"
)

// RmSetBorder represents the SET BORDER command
func (i *Interpreter) RmSetBorder() (ok bool) {
	i.TokenPointer += 2
	// Get required colour
	val, ok := i.OnExpression("numeric")
	if !ok {
		return false
	}
	// TODO: ensure within range --> have a function to do that
	colour := int(math.Round(val.(float64)))
	if i.OnSegmentEnd() {
		i.g.SetBorder(colour)
		return true
	}
	return false
}
