package rmbasicx64

import (
	"math"
	"strconv"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// rmInput represents the INPUT command
// From the manual:
// Read data form a specified input channel
// INPUT [LINE] [#e,][prompt] var[,var2...]
// prompt is a string to be displayed on the screen
// before you read data from the keyboard.  If you
// use a semicolon after the prompt, the string will
// be followed by a question mark.
// The INPUT command reads data from the keyboard or
// specified channel and assigns, appropriately, a
// number of string to each variable.
// These rules apply:
// - string data should be separated by commas, and
//   numeric data by spaces or commas
// - the number of entered values must be the same as
//   the number of variables
// - if you do not enter enough data a double question
//   mark prompt is displayed on the screen
// - if you press <ENTER> without entering data at the
//   keyboard, either an empty string or zero will be
//   assigned to the relevant variable
// - if there is not enough data entered from an input
//   channel, an end of file error will occur
// - if you enter too much data a warning is displayed
//   on the screen
// Take care when using INPUT to read strings that may
// contain spaces or commas.  Such data is best handled
// by the INPUT LINE (or LINE INPUT) command.  INPUT
// LINE reads an entire line of data into one string.
// Example:
// INPUT "Enter 3 numbers :", A, B, C
func (i *Interpreter) RmInput() (ok bool) {
	i.TokenPointer++
	// must have a variable if nothing else
	if !i.OnSegmentEnd() {
		i.ErrorCode = syntaxerror.VariableNameIsNeeded
		return false
	}
	// Optional LINE keyword to set INPUT LINE mode
	inputLineMode := false
	_, ok = i.OnToken([]int{token.LINE})
	if ok {
		inputLineMode = true
	}
	// Optional channel number - default is 0
	// var channel int
	_, ok = i.OnToken([]int{token.Hash})
	if ok {
		// # must be followed be channel number
		_, ok := i.OnExpression("numeric") // val, ok := ...
		if ok {
			// channel = int(math.Round(val.(float64)))
		} else {
			return false
		}
	}
	// Optional prompt - default is ""
	hasPrompt := false
	prompt := ""
	val, ok := i.OnExpression("string")
	if !ok {
		prompt = ""
	} else {
		prompt = val.(string)
		hasPrompt = true
	}
	// Optional semicolon
	promptWithQuestionMark := false
	_, ok = i.OnToken([]int{token.Semicolon})
	if ok {
		promptWithQuestionMark = true
	}
	// If a prompt was passed and no semicolon came after it then we need a
	// comma separator here
	if hasPrompt && !promptWithQuestionMark {
		if !i.OnComma() {
			return false
		}
	}
	// Get required first variable name
	variableName, ok := i.OnVariableName()
	if !ok {
		return false
	}
	variableNames := make([]string, 0)
	variableNames = append(variableNames, variableName.Literal)
	// INPUT LINE mode can't accept more variables
	if inputLineMode && !i.OnSegmentEnd() {
		return false
	}
	// Collect more variable names
	for !i.OnSegmentEnd() {
		variableName, ok := i.OnVariableName()
		if !ok {
			return false
		} else {
			variableNames = append(variableNames, variableName.Literal)
		}
		// If there are more variables then a comma separator is needed
		if !i.OnSegmentEnd() {
			if !i.OnComma() {
				return false
			}
		}
	}
	// Execute
	// Prompt and get input string
	i.g.Print(prompt)
	if promptWithQuestionMark {
		i.g.Print("?")
	}
	rawInput := i.g.Input("")
	// Parse input string and assign values to vars
	if inputLineMode {
		// assign var without parsing
		i.SetVar(variableNames[0], rawInput)
	} else {
		parseInput(i, rawInput, variableNames)
	}
	return true
}

// Applies the parsing rules given in the manual to raw input
func parseInput(i *Interpreter, rawInput string, variableNames []string) {
	// Make a sequence of data types we expect to find in the raw input
	var sequence []int
	numericType := 0
	stringType := 1
	for _, variableName := range variableNames {
		if IsStringVar(token.Token{token.IdentifierLiteral, variableName}) {
			sequence = append(sequence, stringType)
			continue
		}
		if IsFloatVar(token.Token{token.IdentifierLiteral, variableName}) || IsIntVar(token.Token{token.IdentifierLiteral, variableName}) {
			sequence = append(sequence, numericType)
			continue
		}
	}
	// QND parser: split the raw input by spaces, then re-assemble according
	// to the sequence
	segments := strings.Fields(rawInput)
	collectString := ""
	segmentIndex := 0
	maxIndex := 0
	for index, expectedType := range sequence {
		maxIndex++
		if segmentIndex >= len(segments) {
			break
		}
		switch expectedType {
		case stringType:
			for {
				collectString += segments[segmentIndex]
				segmentIndex++
				if collectString[len(collectString)-1:] == "," || segmentIndex == len(segments) {
					break
				}
			}
			// finished collecting string - remove trailing comma if present
			// then set var and wipe collectString
			if collectString[len(collectString)-1:] == "," {
				collectString = collectString[:len(collectString)-1]
			}
			i.SetVar(variableNames[index], collectString)
			collectString = ""
		case numericType:
			val := segments[segmentIndex]
			if val[len(val)-1:] == "," {
				// delete optional comma
				val = val[:len(val)-1]
				if valfloat64, err := strconv.ParseFloat(val, 64); err == nil {
					// Round if variable is integer type
					if IsIntVar(token.Token{token.IdentifierLiteral, variableNames[index]}) {
						valfloat64 = math.Round(valfloat64)
					}
					i.SetVar(variableNames[index], valfloat64)
				}
			}
			segmentIndex++
		}
	}
	// Print ?? if not enough data was entered
	if maxIndex < len(sequence)-1 {
		i.g.Print("??")
		i.g.Put(13)
	}
}
