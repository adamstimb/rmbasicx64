package rmbasicx64

import (
	"os"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// rmBye represents the BYE command
func (i *Interpreter) RmBye() (ok bool) {
	// Ensure no parameters
	i.TokenPointer++
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Exit application with success code
	os.Exit(0)
	// Return never actually happens
	return true
}
