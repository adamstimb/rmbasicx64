package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PrintToken prints a token in the console
func PrintToken(thisToken Token) {
	out, err := json.Marshal(thisToken)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}

// PadString pads a string between spaces
func PadString(s string) string {
	return " " + s + " "
}

// Token defines the actual token returned from the tokenizer and sent to the parser
// and code-formatter
type Token struct {
	Type        int
	Location    []int
	Symbol      string
	ValueString string
	ValueFloat  float64
}

// TokenizeStringLiterals returns tokens of string literals
func TokenizeStringLiterals(code string) []Token {
	// convert string to bytes
	b := []byte(code)
	// tokens are collected in this slice
	tokens := []Token{}
	// use a regex to find a string between double quotes
	r, _ := regexp.Compile("\"(.*?)\"")
	// get all string values and put them in this list
	stringVals := []string{}
	for _, stringVal := range r.FindAll(b, -1) {
		stringVals = append(stringVals, string(stringVal))
	}
	// get all string value positions and add to tokens
	for i, location := range r.FindAllIndex(b, -1) {
		var thisToken Token
		thisToken.Type = LiString
		thisToken.Location = location
		thisToken.Symbol = stringVals[i]
		thisToken.ValueString = stringVals[i]
		tokens = append(tokens, thisToken)
	}
	return tokens
}

// TokenizeNumericalLiterals returns tokens of numerical literals
func TokenizeNumericalLiterals(code string) []Token {
	// convert string to bytes
	b := []byte(code)
	// tokens are collected in this slice
	tokens := []Token{}
	// use a regex to find a string between double quotes
	r, _ := regexp.Compile("[-+]?\\d*\\.\\d+|\\d+")
	// get all numerical values and put them in this list
	// for simplicity we'll store everything to float but in the token we'll use the type ID
	// to remember if it should be a float or int
	numericalVals := []float64{}
	for _, numericalVal := range r.FindAll(b, -1) {
		// try to convert to float64 and add it
		if numericalValFloat, err := strconv.ParseFloat(string(numericalVal), 64); err == nil {
			numericalVals = append(numericalVals, float64(numericalValFloat))
		} else {
			panic("Internal error: Failed to convert numerical literal")
		}
	}
	// get all numerical value positions and add to tokens
	for i, location := range r.FindAllIndex(b, -1) {
		var thisToken Token
		// integer or float?
		if numericalVals[i] == float64(int64(numericalVals[i])) {
			thisToken.Type = LiInteger
		} else {
			thisToken.Type = LiFloat
		}
		thisToken.Location = location
		// covert back to string to get the symbol
		thisToken.Symbol = strconv.FormatFloat(numericalVals[i], 'f', -1, 64)
		thisToken.ValueFloat = numericalVals[i]
		tokens = append(tokens, thisToken)
	}
	return tokens
}

// StringIndexAll is an iterative version of String.Index that returns the indexes of
// all matching substrings in a string
func StringIndexAll(s, sep string) []int {
	indexes := []int{}
	for len(s) > 0 {
		matchingIndex := strings.Index(s, sep)
		if matchingIndex < 0 {
			// not matches so give up
			break
		} else {
			// add this match to list then mask the string up to the end
			// of this match so we don't catch it again
			indexes = append(indexes, matchingIndex)
			masked := []byte(s)
			for i := 0; i < matchingIndex+len(sep); i++ {
				masked[i] = ' '
			}
			s = string(masked)
		}
	}
	return indexes
}

// TokenizeKeywords returns tokens of keywords.  Note that lines divided by ":" need to be
// split before this function is called.
func TokenizeKeywords(code string) []Token {
	// tokens are collected in this slice
	tokens := []Token{}
	// find all matching keyword symbols and add tokens
	for symbol, typeID := range KeywordsToTokens() {
		for _, index := range StringIndexAll(strings.ToUpper(code), PadString(symbol)) {
			var thisToken Token
			thisToken.Type = typeID
			thisToken.Symbol = symbol
			thisToken.Location = []int{index + 1, index + len(symbol)} // take padding of symbol into account
			tokens = append(tokens, thisToken)
		}
	}
	return tokens
}

// TokenizePunctuation returns tokens of punctuation symbols
func TokenizePunctuation(code string) []Token {
	// tokens are collected in this slice
	tokens := []Token{}
	// find all matching punctuation symbols and add tokens
	for symbol, typeID := range PunctuationToTokens() {
		for _, index := range StringIndexAll(code, symbol) {
			var thisToken Token
			thisToken.Type = typeID
			thisToken.Symbol = symbol
			thisToken.Location = []int{index, index + len(symbol) - 1}
			// if a dividing : is followed by = then don't collect this token
			// because it's actually := symbol
			if thisToken.Type == PnDivideLine && thisToken.Location[1] < len(code) {
				if code[thisToken.Location[1]+1:thisToken.Location[1]+2] == "=" {
					// it's := so skip
					continue
				}
			}
			tokens = append(tokens, thisToken)
		}
	}
	return tokens
}

// TokenizeMathematical returns tokens of mathematical symbols
func TokenizeMathematical(code string) []Token {
	// tokens are collected in this slice
	tokens := []Token{}
	// find all matching punctuation symbols and add tokens
	for symbol, typeID := range MathematicalToTokens() {
		for _, index := range StringIndexAll(code, symbol) {
			var thisToken Token
			thisToken.Type = typeID
			thisToken.Symbol = symbol
			thisToken.Location = []int{index, index + len(symbol)}
			tokens = append(tokens, thisToken)
		}
	}
	// remove smaller overlapping symbols
	// e.g. >= will generate tokens for > and >= and = but the symbol is >=
	filteredTokens := []Token{}
	for ai, aToken := range tokens {
		// always add 2 char tokens
		if len(aToken.Symbol) == 2 {
			filteredTokens = append(filteredTokens, aToken)
			continue
		}
		// otherwise check if it overlaps before adding
		noOverlap := true
		for bi, bToken := range tokens {
			if ai != bi && aToken.Location[0] >= bToken.Location[0] && aToken.Location[1] <= bToken.Location[1] {
				// overlaps with something
				noOverlap = false
				break
			}
		}
		if noOverlap {
			filteredTokens = append(filteredTokens, aToken)
		}
	}
	return filteredTokens
}

// TokenizeVariables receives tokens collected so far and code and returns tokens for the variables
func TokenizeVariables(tokens []Token, code string) []Token {
	// new tokens are collected in this slice
	newTokens := []Token{}
	print("code: ")
	println(code)
	// first mask everything that was already tokenized
	bytes := []byte(code)
	for _, thisToken := range tokens {
		for i := thisToken.Location[0]; i <= thisToken.Location[1]; i++ {
			bytes[i] = ' '
		}
	}
	code = string(bytes)
	print("masked code: ")
	println(code)
	// split potential keywords using fields - upper case
	code = strings.ToUpper(code)
	fields := strings.Fields(code)
	// tokenize keywords
	for _, thisField := range fields {
		print(thisField)
		print(">")
		for _, index := range StringIndexAll(code, thisField) {
			var thisToken Token
			thisToken.Location = []int{index, index + len(thisField) - 1}
			// identify type by ending - first assume float because floats have no suffix in RM Basic
			thisToken.Type = MaVariableFloat
			if strings.HasSuffix(thisField, "$") {
				thisToken.Type = MaVariableString
			}
			if strings.HasSuffix(thisField, "%") {
				thisToken.Type = MaVariableInteger
			}
			// format the symbol, which is basically titlize (but _ notation complicates that...to do)
			formattedSymbol := strings.Title(strings.ToLower(thisField))
			thisToken.Symbol = formattedSymbol
			PrintToken(thisToken)
			newTokens = append(newTokens, thisToken)
		}
	}
	return newTokens
}
