package rmbasicx64

import (
	"fmt"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// GetData is an internal function that scans the program for DATA statements
// and populates the programData stack.  It is called automatically by the RUN
// and RESTORE commands.
func (i *Interpreter) GetData() (ok bool) {
	fmt.Println("GetData")
	// clear data stack and set up a new scanner
	s := &Scanner{}
	i.programData = make([]interface{}, 0)
	// search for data
	// 1. Tokenize and get list of segments
	lineOrder := i.GetLineOrder()
	segments := make([][]token.Token, 0)
	for _, lineNumber := range lineOrder {
		// tokenize
		tokens := s.Scan(i.Program[lineNumber])
		// split the tokens into executable segments for each : token found
		this_segment := make([]token.Token, 0)
		for _, t := range tokens {
			if t.TokenType != token.Colon {
				this_segment = append(this_segment, t)
			} else {
				// add EndOfLine to segment before adding to segments slice
				this_segment = append(this_segment, token.Token{token.EndOfLine, ""})
				segments = append(segments, this_segment)
				this_segment = make([]token.Token, 0)
			}
		}
		if len(this_segment) > 0 {
			segments = append(segments, this_segment)
		}
	}
	// 2. If any segment begins with data go ahead and execute it
	i.collectingData = true
	for _, segment := range segments {
		token.PrintToken(segment[0])
		if segment[0].TokenType == token.DATA {
			//i.TokenStack = segment
			//i.TokenPointer = 0
			ok = i.RmData(segment)
			if !ok {
				i.collectingData = false
				return false
			}
		}
	}
	i.collectingData = false
	fmt.Println("end of GetData")
	return true
}

// RmRestore represents the RESTORE command
// RESTORE [l]
func (i *Interpreter) RmRestore() (ok bool) {
	i.TokenPointer++
	// TODO: Accept optional line number
	if !i.OnSegmentEnd() {
		return false
	}
	// Pass thru if not in a program otherwise execute
	if i.LineNumber == -1 {
		return true
	}
	return i.GetData()
}

// RmData represents the DATA command
// DATA v1[,v2 ...]
func (i *Interpreter) RmData(tokens []token.Token) (ok bool) {

	endOfTokens := func(tokens []token.Token, tokenPointer int) bool {
		if tokenPointer > len(tokens) {
			return true
		}
		if tokens[tokenPointer].TokenType == token.EndOfLine {
			return true
		}
		return false
	}

	acceptAnyOfTheseTokens := func(tokens []token.Token, tokenPointer int, acceptableTokens []int) (acceptedToken token.Token, acceptOk bool) {
		for _, tokenType := range acceptableTokens {
			if tokenType == tokens[tokenPointer].TokenType {
				// found a match
				return tokens[tokenPointer], true
			}
		}
		// no matches
		return token.Token{}, false
	}

	tokenPointer := 1

	// Pass if no data
	if endOfTokens(tokens, tokenPointer) {
		return true
	}
	// Execute
	// Get first value
	t, ok := acceptAnyOfTheseTokens(tokens, tokenPointer, []int{token.NumericalLiteral, token.StringLiteral})
	tokenPointer++
	if ok {
		fmt.Println("Got first value")
		// TODO: Error handling for invalid tokens in DATA statement?
		if IsStringVar(t) {
			val, _ := i.GetValueFromToken(t, "string")
			if i.collectingData {
				i.programData = append([]interface{}{val}, i.programData...)
			}
		} else {
			val, _ := i.GetValueFromToken(t, "float64")
			if i.collectingData {
				i.programData = append([]interface{}{val}, i.programData...)
			}
		}
	}
	// Get optional following values
	for !endOfTokens(tokens, tokenPointer) {
		// Get required comma
		_, ok := acceptAnyOfTheseTokens(tokens, tokenPointer, []int{token.Comma})
		tokenPointer++
		if !ok {
			i.ErrorCode = syntaxerror.CommaSeparatorIsNeeded
			return false
		}
		// Get value if present --- can/should we handle this better?
		t, ok := acceptAnyOfTheseTokens(tokens, tokenPointer, []int{token.NumericalLiteral, token.StringLiteral})
		tokenPointer++
		if ok {
			if IsStringVar(t) {
				val, _ := i.GetValueFromToken(t, "string")
				if i.collectingData {
					i.programData = append([]interface{}{val}, i.programData...)
				}
			} else {
				val, _ := i.GetValueFromToken(t, "float64")
				if i.collectingData {
					i.programData = append([]interface{}{val}, i.programData...)
				}
			}
		}
	}
	fmt.Printf("len programData = %d\n", len(i.programData))
	return true
}
