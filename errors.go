package main

// Error codes defined here
const (
    ErCommaSeparatorIsNeeded = 2000
    ErOpeningBracketIsNeeded = 2001
    ErClosingBracketIsNeeded = 2002
    ErSemicolonSeparatorNeeded = 2003
    ErUnknownValueStackItem = 2004
    ErVariableNameIsNeeded = 2005
    ErArrayVariableIsWrong = 2006
    ErNumericVariableNeeded = 2007
    ErTrueFalseValueNeeded = 2008
    ErNumberOrStringNeeded = 2009
    ErInvalidExpressionFound = 2010
    ErNumericExpressionNeeded = 2011
    ErStringExpressionNeeded = 2012
    ErExpressionTooComplicated = 2013
    ErWrongNumberOfSubscripts = 2014
    ErArrayAlreadyDimensioned = 2015
    ErArrayNeededNoSubscript = 2016
    ErSpecifiedLineNotFound = 2017
    ErLineNumberLabelNeeded = 2018
    ErUnknownCommandProcedure = 2019
    ErFunctionArrayNotFound = 2020
    ErBadArgumentForFunction = 2021
    ErWidthSpecificationWrong = 2022
    ErTooManyGotosToRenumber = 2023
    ErGotoOrGosubIsMissing = 2024
    ErToIsNeededBeforeValue = 2025
    ErThenWithNoMatchingIf = 2026
    ErReturnWithoutAnyGosub = 2027
    ErNoErrorToContinueFrom = 2028
    ErUntilWithoutAnyRepeat = 2029
    ErNextWithoutMatchingFor = 2030
    ErTooManyForRepeatGosub = 2031
    ErWrongSetAskAttribute = 2032
    ErUnknownSetAskAttribute = 2033
    ErDirectoryCannotBeFound = 2034
    ErUnableToRemoveDirectory = 2035
    ErUnableToCreateDirectory = 2036
    ErExactFilenameIsNeeded = 2037
    ErUnableToRenameTheFile = 2038
    ErUnableToEraseTheFile = 2039
    ErNamedChannelIsAlreadyInUse = 2040
    ErUnableToOpenNamedFile = 2041
    ErChannelNotOpenForInput = 2042
    ErChannelNotOpenForOutput = 2043
    ErWrongChannelNumberUsed = 2044
    ErNamedFileAlreadyExists = 2045
    ErReadingPastEndOfFile = 2046
    ErExcessInputDataIgnored = 2047
    ErUnableToReadExcessData = 2048
    ErNoMoreDataToBeRead = 2049
    ErWritingAreaInappropriate = 2050
    ErEnvelopeNumberIsWrong = 2051
    ErPointCoordinatesNeeded = 2052
    ErACoordinatePairNeeded = 2053
    ErCoordinateIsOutOfRange = 2054
    ErFillAreaTooComplicated = 2055
    ErBadFloodStartCoordinates = 2056
    ErMouseOrJoystickProblem = 2057
    ErGeneralGraphicsFailure = 2058
    ErParameterNotUnderstood = 2059
    ErCoordinateArrayIsWrong = 2060
    ErAllAvailableMemoryUsed = 2061
    ErEnteredLineIsTooLong = 2062
    ErNumberNotInAllowedRange = 2063
    ErLineNumberOutOfRange = 2064
    ErLengthOfStringTooGreat = 2065
    ErStringOffsetOutOfRange = 2066
    ErStepValueNotLargeEnough = 2067
    ErPositiveValueRequired = 2068
    ErTryingToDivideByZero = 2069
    ErEmptyStringNotAllowed = 2070
    ErProcedureExitWithoutCall = 2071
    ErArraySubscriptIsWrong = 2072
    ErVariableUsedAsALocal = 2073
    ErVariableWithoutAnyValue = 2074
    ErInterruptedByBreakKey = 2075
    ErEndOfInstructionExpected = 2076
    ErMaximumArgumentsExeceeded = 2077
    ErReceiveVariablesNeeded = 2078
    ErNotEnoughParameters = 2079
    ErFunctionExitWithoutCall = 2080
    ErNeedResultToExitFunction = 2081
    ErCannotExecuteDefinition = 2082
    ErObjectVariableExpected = 2083
    ErArraysNotOfSameType = 2084
    ErNameOfDefinitionRequired = 2085
    ErEndOfDefinitionExpected = 2086
    ErFunctionNestingTooDeep = 2087
    ErTooManyFilesOpen = 2088
    ErFileSpecifiedInvalidOrDoesnTExist = 2089
    ErNoMoreSpaceAvailableOnDisk = 2090
    ErFileInUseOrNameIsADirVol = 2091
    ErFileOperationFailure = 2092
    ErDirectoryFull = 2093
    ErUnableToCloseFile = 2094
    ErNoSuchLevel = 2095
    ErConfigurationCommandUnknown = 2096
    ErErrorInLoadingExtension = 2097
    ErNotAnExtensionFile = 2098
    ErFilenameIsADirectory = 2099
    ErDiskIsReadOnly = 2100
    ErDriveNotReady = 2101
    ErWarning = 2102
    ErInLine = 2103
    ErInCommand = 2104
)

// ErrorMessages returns a map of error codes and their messages
func ErrorMessages() map[int]string {
    return map[int]string{
        ErCommaSeparatorIsNeeded: "Comma separator is needed",
        ErOpeningBracketIsNeeded: "Opening bracket is needed",
        ErClosingBracketIsNeeded: "Closing bracket is needed",
        ErSemicolonSeparatorNeeded: "Semicolon separator needed",
        ErUnknownValueStackItem: "Unknown value stack item",
        ErVariableNameIsNeeded: "Variable name is needed",
        ErArrayVariableIsWrong: "Array variable is wrong",
        ErNumericVariableNeeded: "Numeric variable needed",
        ErTrueFalseValueNeeded: "TRUE/FALSE value needed",
        ErNumberOrStringNeeded: "Number or string needed",
        ErInvalidExpressionFound: "Invalid expression found",
        ErNumericExpressionNeeded: "Numeric expression needed",
        ErStringExpressionNeeded: "String expression needed",
        ErExpressionTooComplicated: "Expression too complicated",
        ErWrongNumberOfSubscripts: "Wrong number of subscripts",
        ErArrayAlreadyDimensioned: "Array already dimensioned",
        ErArrayNeededNoSubscript: "Array needed (no subscript)",
        ErSpecifiedLineNotFound: "Specified line not found",
        ErLineNumberLabelNeeded: "Line number/label needed",
        ErUnknownCommandProcedure: "Unknown command/procedure",
        ErFunctionArrayNotFound: "Function/array not found",
        ErBadArgumentForFunction: "Bad argument for function",
        ErWidthSpecificationWrong: "WIDTH specification wrong",
        ErTooManyGotosToRenumber: "Too many GOTOs to renumber",
        ErGotoOrGosubIsMissing: "GOTO or GOSUB is missing",
        ErToIsNeededBeforeValue: "TO is needed before value",
        ErThenWithNoMatchingIf: "THEN with no matching IF",
        ErReturnWithoutAnyGosub: "RETURN without any GOSUB",
        ErNoErrorToContinueFrom: "No error to continue from",
        ErUntilWithoutAnyRepeat: "UNTIL without any REPEAT",
        ErNextWithoutMatchingFor: "NEXT without matching FOR",
        ErTooManyForRepeatGosub: "Too many FOR/REPEAT/GOSUB",
        ErWrongSetAskAttribute: "Wrong SET/ASK attribute",
        ErUnknownSetAskAttribute: "Unknown SET/ASK attribute",
        ErDirectoryCannotBeFound: "Directory cannot be found",
        ErUnableToRemoveDirectory: "Unable to remove directory",
        ErUnableToCreateDirectory: "Unable to create directory",
        ErExactFilenameIsNeeded: "Exact filename is needed",
        ErUnableToRenameTheFile: "Unable to rename the file",
        ErUnableToEraseTheFile: "Unable to erase the file",
        ErNamedChannelIsAlreadyInUse: "Named channel is already in use",
        ErUnableToOpenNamedFile: "Unable to open named file",
        ErChannelNotOpenForInput: "Channel not open for input",
        ErChannelNotOpenForOutput: "Channel not open for output",
        ErWrongChannelNumberUsed: "Wrong channel number used",
        ErNamedFileAlreadyExists: "Named file already exists",
        ErReadingPastEndOfFile: "Reading past end of file",
        ErExcessInputDataIgnored: "Excess input data ignored",
        ErUnableToReadExcessData: "Unable to read excess DATA",
        ErNoMoreDataToBeRead: "No more DATA to be read",
        ErWritingAreaInappropriate: "Writing area inappropriate",
        ErEnvelopeNumberIsWrong: "Envelope number is wrong",
        ErPointCoordinatesNeeded: "Point coordinates needed",
        ErACoordinatePairNeeded: "A coordinate pair needed",
        ErCoordinateIsOutOfRange: "Coordinate is out of range",
        ErFillAreaTooComplicated: "Fill area too complicated",
        ErBadFloodStartCoordinates: "Bad FLOOD start coordinates",
        ErMouseOrJoystickProblem: "Mouse or joystick problem",
        ErGeneralGraphicsFailure: "General graphics failure",
        ErParameterNotUnderstood: "Parameter not understood",
        ErCoordinateArrayIsWrong: "Coordinate array is wrong",
        ErAllAvailableMemoryUsed: "All available memory used",
        ErEnteredLineIsTooLong: "Entered line is too long",
        ErNumberNotInAllowedRange: "Number not in allowed range",
        ErLineNumberOutOfRange: "Line number out of range",
        ErLengthOfStringTooGreat: "Length of string too great",
        ErStringOffsetOutOfRange: "String offset out of range",
        ErStepValueNotLargeEnough: "Step value not large enough",
        ErPositiveValueRequired: "Positive value required",
        ErTryingToDivideByZero: "Trying to divide by zero",
        ErEmptyStringNotAllowed: "Empty string not allowed",
        ErProcedureExitWithoutCall: "Procedure exit without call",
        ErArraySubscriptIsWrong: "Array subscript is wrong",
        ErVariableUsedAsALocal: "Variable used as a local",
        ErVariableWithoutAnyValue: "Variable without any value",
        ErInterruptedByBreakKey: "Interrupted by BREAK key",
        ErEndOfInstructionExpected: "End of instruction expected",
        ErMaximumArgumentsExeceeded: "Maximum arguments execeeded",
        ErReceiveVariablesNeeded: "RECEIVE variables needed",
        ErNotEnoughParameters: "Not enough parameters",
        ErFunctionExitWithoutCall: "Function exit without call",
        ErNeedResultToExitFunction: "Need RESULT to exit function",
        ErCannotExecuteDefinition: "Cannot execute definition",
        ErObjectVariableExpected: "Object variable expected",
        ErArraysNotOfSameType: "Arrays not of same type",
        ErNameOfDefinitionRequired: "Name of definition required",
        ErEndOfDefinitionExpected: "End of definition expected",
        ErFunctionNestingTooDeep: "Function nesting too deep",
        ErTooManyFilesOpen: "Too many files open",
        ErFileSpecifiedInvalidOrDoesnTExist: "File specified invalid or doesn't exist",
        ErNoMoreSpaceAvailableOnDisk: "No more space available on disk",
        ErFileInUseOrNameIsADirVol: "File in use or name is a Dir/Vol",
        ErFileOperationFailure: "File operation failure",
        ErDirectoryFull: "Directory full",
        ErUnableToCloseFile: "Unable to close file",
        ErNoSuchLevel: "No such level",
        ErConfigurationCommandUnknown: "Configuration command unknown",
        ErErrorInLoadingExtension: "Error in loading extension",
        ErNotAnExtensionFile: "Not an extension file",
        ErFilenameIsADirectory: "Filename is a directory",
        ErDiskIsReadOnly: "Disk is read only",
        ErDriveNotReady: "Drive not ready",
        ErWarning: "Warning : ",
        ErInLine: " in line ",
        ErInCommand: " in command",
    }
}
