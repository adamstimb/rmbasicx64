package rmbasicx64

import (
	"fmt"
	"math"
	"os"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// rmUpdateLine represents a line number being updated in the REPL
func (i *Interpreter) RmUpdateLine(code string) (ok bool) {
	// If the code begins with a round number then that's a line number
	if i.CurrentTokens[0].TokenType == token.NumericalLiteral {
		val, ok := i.GetValueFromToken(i.CurrentTokens[0], "float64")
		if ok {
			lineNumber := val.(float64)
			if lineNumber == math.Round(lineNumber) {
				// Is a line number so ...
				// If there are no more tokens delete the line if it exists
				// otherwise update it
				if i.CurrentTokens[1].TokenType == token.EndOfLine {
					if _, ok := i.Program[int(lineNumber)]; ok {
						delete(i.Program, int(lineNumber))
						return true
					} else {
						// line to delete does not exist so pass through
						return true
					}
				} else {
					i.Program[int(lineNumber)] = i.FormatCode(code, -1, true)
					return true
				}
			} else {
				i.ErrorCode = syntaxerror.LineNumberExpected
				i.BadTokenIndex = 0
				return false
			}
		}
	}
	i.ErrorCode = syntaxerror.LineNumberExpected
	i.BadTokenIndex = 0
	return false
}

// RmEdit represents the EDIT command
func (i *Interpreter) RmEdit() (ok bool) {
	i.TokenPointer++
	if i.EndOfTokens() {
		// No linenumber passed so get line of last error
		// TODO: implement
		return true
	}
	_, ok = i.AcceptAnyOfTheseTokens([]int{token.NumericalLiteral, token.IdentifierLiteral})
	if ok {
		i.TokenPointer--
		lineNumber, ok := i.AcceptAnyNumber()
		if ok {
			if lineNumber == math.Round(lineNumber) {
				if _, ok := i.Program[int(lineNumber)]; ok {
					// edit existing line
					edited := i.g.Input("", fmt.Sprintf("%d %s", int(lineNumber), i.Program[int(lineNumber)]))
					_ = i.ImmediateInput(edited)
					return true
				} else {
					// create new line
					edited := i.g.Input("", fmt.Sprintf("%d %s", int(lineNumber), ""))
					_ = i.ImmediateInput(edited)
					return true
				}
			} else {
				i.ErrorCode = syntaxerror.LineNumberExpected
				return false
			}
		}
	}
	return false
}

// rmAuto represents the AUTO command
// From the manual:
// AUTO [line][, increment]
// Generate line numbers automatically. If you enter AUTO
// on its own, line number 10 will be displayed.  After you
// type instruction(s) for that line and <ENTER>, 20 will be
// displayed at the beginning of the next line.  Use line
// when you want a different starting line number.  Use
// increment to specify the increment between lines.  Before
// a line number is generate, any existing lines with numbers
// between the new line and the last auto-generated line are
// listed.  Consequently you can see what you are skipping
// over or what you are about to replace.  If you do not want
// to overwrite existing lines, you must break out of automatic
// generation.
func (i *Interpreter) RmAuto() (ok bool) {
	i.TokenPointer++
	// Default parameters
	startLine := 10
	increment := 10
	// Try to collect parameters, if any
	if !i.EndOfTokens() {
		if i.TokenStack[1].TokenType == token.NumericalLiteral || i.TokenStack[1].TokenType == token.IdentifierLiteral {
			// Consume startLine
			val, ok := i.AcceptAnyNumber()
			if !ok {
				return false
			}
			startLine = int(math.Round(val))
			if !i.EndOfTokens() {
				// Consume comma
				_, ok := i.AcceptAnyOfTheseTokens([]int{token.Comma})
				if !ok {
					i.ErrorCode = syntaxerror.CommaSeparatorIsNeeded
					i.BadTokenIndex = i.TokenPointer
					return false
				}
				// Consume increment
				val, ok := i.AcceptAnyNumber()
				if !ok {
					return false
				}
				increment = int(math.Round(val))
				// End of expression
				if !i.EndOfTokens() {
					i.ErrorCode = syntaxerror.EndOfInstructionExpected
					i.BadTokenIndex = i.TokenPointer
					return false
				}
			}
		} else {
			i.ErrorCode = syntaxerror.LineNumberExpected
			i.BadTokenIndex = i.TokenPointer
			return false
		}
	}
	// Execute
	// TODO: Implement CTRL-ScrollLock interrupt in nimgobus
	autoLineNumber := startLine
	previousAutoLineNumber := autoLineNumber
	for {
		// list any existing lines between this line and the previous autogenerated line
		lineOrder := i.GetLineOrder()
		for _, lineNumber := range lineOrder {
			if lineNumber > previousAutoLineNumber && lineNumber <= autoLineNumber {
				i.g.Print(fmt.Sprintf("%d %s", lineNumber, i.Program[lineNumber]))
			}
		}
		// Get instructions for this line then increment lineNumber
		edited := i.g.Input("", fmt.Sprintf("%d %s", int(autoLineNumber), ""))
		_ = i.ImmediateInput(edited)
		previousAutoLineNumber = autoLineNumber
		autoLineNumber += increment
	}
}

// rmBye represents the BYE command
func (i *Interpreter) RmBye() (ok bool) {
	// Ensure no parameters
	i.TokenPointer++
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Exit application with success code
	os.Exit(0)
	// Return never actually happens
	return true
}

// rmAssign represents a variable assignment (var = expr or var := expr)
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

// rmRun represents the RUN command
func (i *Interpreter) RmRun() (ok bool) {
	// pass thru is no program
	if len(i.Program) == 0 {
		return true
	}
	i.ProgramPointer = 0
	lineOrder := i.GetLineOrder()
	// Check for optional startFrom parameter
	startFrom := lineOrder[i.ProgramPointer]
	if len(i.TokenStack) > 2 {
		i.TokenPointer++
		_, ok := i.AcceptAnyOfTheseTokens([]int{token.NumericalLiteral, token.IdentifierLiteral})
		if ok {
			i.TokenPointer--
			val, ok := i.AcceptAnyNumber()
			if ok {
				startFrom = int(math.Round(val))
			} else {
				return false
			}
		} else {
			return false
		}
	}
	// Run the program and if startFrom was passed scan ahead to it before executing
	scanAhead := true
	for i.ProgramPointer < len(lineOrder) {
		i.LineNumber = lineOrder[i.ProgramPointer]
		if scanAhead && (i.LineNumber != startFrom) {
			i.ProgramPointer++
			continue
		}
		if scanAhead && (i.LineNumber == startFrom) {
			scanAhead = false
		}
		ok := i.RunLine(i.Program[lineOrder[i.ProgramPointer]])
		if !ok {
			i.ProgramPointer = 0
			return false
		}
	}
	// Catch line number not found
	if scanAhead {
		i.ErrorCode = syntaxerror.SpecifiedLineNotFound
		i.BadTokenIndex = 1
		i.LineNumber = -1
		return false
	}
	// Otherwise ok
	i.ProgramPointer = 0
	return true
}