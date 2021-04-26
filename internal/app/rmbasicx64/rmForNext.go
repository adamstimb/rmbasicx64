package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// RmFor represents the FOR command
// FOR v [:]= e1 TO e2 [STEP e3]
func (i *Interpreter) RmFor() (ok bool) {
	i.TokenPointer++
	// Get required controlVarName
	t, ok := i.OnVariableName()
	if !ok {
		return false
	}
	controlVarName := t.Literal
	// controlVarName cannot be a string variable
	if IsStringVar(t) {
		i.ErrorCode = syntaxerror.NumericVariableNeeded
		return false
	}
	// Get assign (:= or =)
	_, ok = i.OnToken([]int{token.Assign, token.Equal})
	if !ok {
		i.ErrorCode = syntaxerror.InvalidExpressionFound
		return false
	}
	// Get controlVarValue
	controlVarValue, ok := i.OnExpression("numeric")
	if !ok {
		return false
	}
	// Get required TO
	_, ok = i.OnToken([]int{token.Assign, token.TO})
	if !ok {
		i.ErrorCode = syntaxerror.InvalidExpressionFound
		return false
	}
	// Get required targetVal
	targetVal, ok := i.OnExpression("numeric")
	if !ok {
		return false
	}
	// Get optional STEP
	step := float64(1) // default is 1
	if !i.OnSegmentEnd() {
		// Get STEP keyword
		_, ok = i.OnToken([]int{token.STEP})
		if !ok {
			i.ErrorCode = syntaxerror.InvalidExpressionFound
			return false
		}
		// Get step value
		val, ok := i.OnExpression("numeric")
		if !ok {
			return false
		}
		step = val.(float64)
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
	newItem := forStackItem{i.ProgramPointer, controlVarName, targetVal.(float64), step}
	i.forStack = append([]forStackItem{newItem}, i.forStack...)
	return true
}

// RmNext represents the NEXT command
// NEXT [v]
func (i *Interpreter) RmNext() (ok bool) {
	// TODO: Collect optional [v]
	i.TokenPointer++
	if !i.OnSegmentEnd() {
		return false
	}
	// Execute
	// Increment control var value and if it does not exceed target re-iterate
	// otherwise pop stack and pass through
	if len(i.forStack) == 0 {
		i.ErrorCode = syntaxerror.NextWithoutMatchingFor
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
