package syntaxerror

// Error codes defined here
const (
	Success = iota
	ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall
	CouldNotInterpretAsANumber
	DidNotExpectInExpression
	HasNotBeenDefined
	IsAKeywordAndCannotBeUsedAsAVariableName
	InvalidExpression
	CannotPerformBitwiseOperationsOnFloatValues
	CannotPerformBitwiseOperationsOnStringValues
	TooManyParametersFor
	NotEnoughParametersFor
	LineNumberExpected
	LineNumberDoesNotExist
	UnknownCommandProcedure
	TryingToDivideByZero
	EndOfInstructionExpected
	NumericExpressionNeeded
	InvalidExpressionFound
	SpecifiedLineNotFound
	LineNumberLabelNeeded
	StringExpressionNeeded
	ExactFilenameIsNeeded
	FilenameIsADirectory
	FileOperationFailure
	UnableToOpenNamedFile
	CommaSeparatorIsNeeded
	InterruptedByBreakKey
	VariableNameIsNeeded
)

// ErrorMessage returns the template error message for a given error code
func ErrorMessage(errorCode int) string {
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
		TooManyParametersFor:                         "Too many parameters for ",
		NotEnoughParametersFor:                       "Not enough parameters for ",
		LineNumberExpected:                           "Line number expected",
		LineNumberDoesNotExist:                       "Line number does not exist",
		TryingToDivideByZero:                         "Trying to divide by zero",    //70
		UnknownCommandProcedure:                      "Unknown command/procedure",   //20
		EndOfInstructionExpected:                     "End of instruction expected", //77
		NumericExpressionNeeded:                      "Numeric expression needed",   //12
		InvalidExpressionFound:                       "Invalid expression found",    //11
		SpecifiedLineNotFound:                        "Specified line not found",    //18
		LineNumberLabelNeeded:                        "Line number/label needed",    //19
		StringExpressionNeeded:                       "String expression needed",    //13
		ExactFilenameIsNeeded:                        "Exact filename is needed",    //38
		FilenameIsADirectory:                         "Filename is a directory",     //100
		FileOperationFailure:                         "File operation failure",      //93
		UnableToOpenNamedFile:                        "Unable to open named file",   //42
		CommaSeparatorIsNeeded:                       "Comma separator is needed",   //1
		InterruptedByBreakKey:                        "Interrupted by BREAK key",
		VariableNameIsNeeded:                         "Variable name is needed", // 6
	}
	return errorMessages[errorCode]
}
