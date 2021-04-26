package rmbasicx64

import (
	"fmt"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// RmIfThenElse represents an IF ... THEN ... ELSE ... statement
// TODO: Allow multiple segments in the THEN and ELSE blocks
func (i *Interpreter) RmIfThenElse() (ok bool) {
	// Ensure parameter is passed
	i.TokenPointer++
	if i.EndOfTokens() {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		return false
	}
	// Get required expression
	val, ok := i.AcceptAnyNumber()
	if !ok {
		return false
	}
	// Try to evaluate result
	// TODO: Does IF only look for TRUE or does it reject anything that is neither TRUE nor FALSE?
	var result bool
	if IsTrue(val) {
		result = true
	} else {
		result = false
	}
	// Get required THEN token
	_, ok = i.AcceptAnyOfTheseTokens([]int{token.THEN})
	if !ok {
		i.ErrorCode = syntaxerror.ThenExpected
		return false
	}
	// Collect thenSegment until end of tokens or ELSE token
	thenSegment := make([]token.Token, 0)
	fmt.Println("Collect thenSegment")
	for !i.EndOfTokens() {
		fmt.Println("iter")
		if i.TokenStack[i.TokenPointer].TokenType == token.ELSE {
			break
		}
		thenSegment = append(thenSegment, i.TokenStack[i.TokenPointer])
		token.PrintToken(i.TokenStack[i.TokenPointer])
		i.TokenPointer++
	}
	thenSegment = append(thenSegment, token.Token{token.EndOfLine, ""})
	// Collect optional elseSegment
	elseSegment := make([]token.Token, 0)
	fmt.Println("Collect elseSegment")
	if !i.EndOfTokens() {
		// Get required ELSE statement
		_, ok = i.AcceptAnyOfTheseTokens([]int{token.ELSE})
		if !ok {
			i.ErrorCode = syntaxerror.EndOfInstructionExpected
			return false
		}
		// Collect elseSegment until end of tokens
		for !i.EndOfTokens() {
			elseSegment = append(elseSegment, i.TokenStack[i.TokenPointer])
			i.TokenPointer++
		}
		elseSegment = append(elseSegment, token.Token{token.EndOfLine, ""})
	}
	// Depending on result, execute the thenSegment or the elseSegment
	if result {
		fmt.Println("True - execute THEN")
		return i.RunSegment(thenSegment)
	} else {
		if len(elseSegment) > 0 {
			fmt.Println("True - execute THEN")
			return i.RunSegment(elseSegment)
		}
	}
	return true
}
