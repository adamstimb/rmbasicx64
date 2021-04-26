package rmbasicx64

import (
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// rmAssign represents a variable assignment (var = expr or var := expr)
func (i *Interpreter) RmAssign() (ok bool) {
	// Consume optional LET
	_, _ = i.OnToken([]int{token.LET})
	// Catch case where a keyword has been used as a variable name to assign to
	if IsKeyword(i.TokenStack[i.TokenPointer]) &&
		(i.TokenStack[i.TokenPointer+1].TokenType == token.Equal || i.TokenStack[i.TokenPointer+1].TokenType == token.Assign) {
		i.ErrorCode = syntaxerror.InvalidExpressionFound
		return false
	}
	varNamePosition := i.TokenPointer
	// advance token point, extract expression, evaluate result then store
	i.TokenPointer += 2
	result, ok := i.EvaluateExpression()
	if ok {
		// Evaluation was successful so check data type and store
		if i.SetVar(i.TokenStack[varNamePosition].Literal, result) {
			return true
		} else {
			return false
		}
	} else {
		// Something went wrong in the evaluation
		return false
	}
}
