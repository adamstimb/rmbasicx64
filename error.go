package main

// Error codes define here
const (
	Success                                                               = iota //= 1001
	ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall        //= 1002
	CouldNotInterpretAsANumber                                                   //= 1003
	DidNotExpectInExpression                                                     //= 1004
	HasNotBeenDefined                                                            //= 1005
	IsAKeywordAndCannotBeUsedAsAVariableName                                     //= 1006
	InvalidExpression                                                            //= 1007
	CannotPerformBitwiseOperationsOnFloatValues                                  //= 1008
	CannotPerformBitwiseOperationsOnStringValues                                 //= 1009
)

// errorMessage returns the template error message for a given error code
func errorMessage(errorCode int) string {
	errorMessages := map[int]string{
		Success: "",
		ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall: "Expected a keyword, line number, expression, variable assignment or procedure call",
		CouldNotInterpretAsANumber:                   " could not be interpreted as a number",
		DidNotExpectInExpression:                     " was not expected in expression",
		HasNotBeenDefined:                            " has not been defined",
		IsAKeywordAndCannotBeUsedAsAVariableName:     " is a keyword and cannot be used as a variable name",
		InvalidExpression:                            " caused an invalid expression",
		CannotPerformBitwiseOperationsOnFloatValues:  "Cannot perform bitwise operations on float values",
		CannotPerformBitwiseOperationsOnStringValues: "Cannot perform bitwise operations on string values",
	}
	return errorMessages[errorCode]
}
