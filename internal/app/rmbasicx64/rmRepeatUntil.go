package rmbasicx64

import (
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// RmRepeat represents the REPEAT command
// REPEAT
func (i *Interpreter) RmRepeat() (ok bool) {
	i.TokenPointer++
	// Ensure no parameters passed
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Execute
	// Pass through if this loop is already registered & iterating
	if len(i.repeatStack) > 0 {
		if i.repeatStack[0].startProgramPointer == i.ProgramPointer {
			return true
		}
	}
	// Push new forStackItem to forStack
	newItem := repeatStackItem{i.ProgramPointer}
	i.repeatStack = append([]repeatStackItem{newItem}, i.repeatStack...)
	return true
}

// RmUntil represents the UNTIL command
// UNTIL e
func (i *Interpreter) RmUntil() (ok bool) {
	i.TokenPointer++
	// Ensure there is a repeat to go back to
	if len(i.repeatStack) == 0 {
		i.ErrorCode = syntaxerror.UntilWithoutAnyRepeat
		i.BadTokenIndex = 0
		return false
	}
	// Get required expression e
	if i.EndOfTokens() {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	val, ok := i.AcceptAnyNumber()
	if !ok {
		return false
	}
	// If result is true then pop the stack and pass through, otherwise go back
	// to top of loop
	thisItem := i.repeatStack[0]
	if IsTrue(val) {
		// pop
		i.repeatStack = i.repeatStack[1:]
		return true
	}
	i.ProgramPointer = thisItem.startProgramPointer
	return true
}
