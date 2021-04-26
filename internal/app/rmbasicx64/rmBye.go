package rmbasicx64

import (
	"os"
)

// rmBye represents the BYE command
func (i *Interpreter) RmBye() (ok bool) {
	i.TokenPointer++
	if !i.OnSegmentEnd() {
		return false
	}
	// Exit application with success code
	os.Exit(0)
	// Return never actually happens
	return true
}
