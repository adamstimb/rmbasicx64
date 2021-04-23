package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// RmFor represents the FOR command
// FOR v [:]= e1 TO e2 [STEP e3]
func (i *Interpreter) RmFor() (ok bool) {
	// Ensure a parameter is passed
	i.TokenPointer++
	if i.EndOfTokens() {
		i.ErrorCode = syntaxerror.NumericVariableNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get required controlVarName and make sure it isn't string
	t, ok := i.AcceptAnyOfTheseTokens([]int{token.IdentifierLiteral})
	if !ok || IsStringVar(t) {
		i.ErrorCode = syntaxerror.NumericVariableNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	controlVarName := t.Literal
	// Get assign (:= or =)
	_, ok = i.AcceptAnyOfTheseTokens([]int{token.Assign, token.Equal})
	if !ok {
		i.ErrorCode = syntaxerror.InvalidExpressionFound
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get controlVarValue
	controlVarValue, ok := i.AcceptAnyNumber()
	if !ok {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get required TO
	_, ok = i.AcceptAnyOfTheseTokens([]int{token.Assign, token.TO})
	if !ok {
		i.ErrorCode = syntaxerror.InvalidExpressionFound
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get required targetVal
	targetVal, ok := i.AcceptAnyNumber()
	if !ok {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get optional STEP
	step := float64(1) // default is 1
	if !i.EndOfTokens() {
		// Get STEP keyword
		_, ok = i.AcceptAnyOfTheseTokens([]int{token.Assign, token.STEP})
		if !ok {
			i.ErrorCode = syntaxerror.InvalidExpressionFound
			i.BadTokenIndex = i.TokenPointer
			return false
		}
		// Get step value
		if i.EndOfTokens() {
			i.ErrorCode = syntaxerror.NumericExpressionNeeded
			i.BadTokenIndex = i.TokenPointer
			return false
		}
		step, ok = i.AcceptAnyNumber()
		if !ok {
			return false
		}
	}
	// Execute
	// Pass through if this loop is already registered & iterating
	if len(i.forStack) > 0 {
		if i.forStack[0].startProgramPointer == i.ProgramPointer {
			// Is already registered so pass
			// This won't work if another FOR loop is declared on the same line, mind.
			// Also if the FOR statement is not in the first segment of a line.
			return true
		}
	}
	// Set the controlVar then push new forStackItem to forStack
	i.SetVar(controlVarName, controlVarValue)
	newItem := forStackItem{i.ProgramPointer, controlVarName, targetVal, step}
	i.forStack = append([]forStackItem{newItem}, i.forStack...)
	return true
}

// RmNext represents the NEXT command
// NEXT [v]
func (i *Interpreter) RmNext() (ok bool) {
	// TODO: Collect optional [v]
	i.TokenPointer++
	if i.EndOfTokens() {
		// Increment control var value and if it does not exceed target re-iterate
		// otherwise pop stack and pass through
		if len(i.forStack) == 0 {
			i.ErrorCode = syntaxerror.NextWithoutMatchingFor
			i.BadTokenIndex = i.TokenPointer
			return false
		}
		thisItem := i.forStack[0]
		val, _ := i.GetVar(thisItem.controlVarName)
		valfloat64 := val.(float64)
		valfloat64 += thisItem.step
		i.SetVar(thisItem.controlVarName, valfloat64)
		if math.Abs(valfloat64) > math.Abs(thisItem.targetVal) {
			// pop
			i.forStack = i.forStack[1:]
			return true
		} else {
			// return to opening FOR statement
			i.ProgramPointer = thisItem.startProgramPointer
			return true
		}
	}
	return true
}
