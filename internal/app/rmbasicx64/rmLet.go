package rmbasicx64

import (
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// rmAssign represents a variable assignment (var = expr or var := expr)
// TODO: Also allow instruction to begin with optional LET
func (i *Interpreter) RmAssign() (ok bool) {
	// Catch case where a keyword has been used as a variable name to assign to
	if IsKeyword(i.TokenStack[0]) &&
		(i.TokenStack[1].TokenType == token.Equal || i.TokenStack[1].TokenType == token.Assign) {
		i.ErrorCode = syntaxerror.InvalidExpressionFound
		i.BadTokenIndex = 0
		return false
	}
	// advance token point, extract expression, evaluate result then store
	i.TokenPointer += 2
	result, ok := i.EvaluateExpression()
	if ok {
		// Evaluation was successful so check data type and store
		if i.SetVar(i.TokenStack[0].Literal, result) {
			return true
		} else {
			return false
		}
	} else {
		// Something went wrong in the evaluation
		return false
	}
}
