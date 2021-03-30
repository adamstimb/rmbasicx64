package main

// Error codes define here
const (
	Success = iota
	ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall
	CouldNotInterpretAsANumber
	DidNotExpectInExpression
	HasNotBeenDefined
	IsAKeywordAndCannotBeUsedAsAVariableName
)

// errorMessage returns the template error message for a given error code
func errorMessage(errorCode int) string {
	errorMessages := map[int]string{
		Success: "",
		ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall: "Expected a keyword, line number, expression, variable assignment or procedure call",
		CouldNotInterpretAsANumber:               " could not be interpreted as a number",
		DidNotExpectInExpression:                 " was not expected in expression",
		HasNotBeenDefined:                        " has not been defined",
		IsAKeywordAndCannotBeUsedAsAVariableName: " is a keyword and cannot be used as a variable name",
	}
	return errorMessages[errorCode]
}
