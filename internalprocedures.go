package main

import "os"

func parseBye(g *Game, tokens []Token) int {
	logMsg("ParseBye")
	if len(tokens) > 1 {
		return ErEndOfInstructionExpected
	}
	return internalBye(g)
}

func internalBye(g *Game) int {
	// Exit with success code
	logMsg("internalBye")
	os.Exit(0)
	return 0
}

func parsePrint(g *Game, tokens []Token) int {
	logMsg("ParsePrint")
	if len(tokens) == 1 {
		return ErNotEnoughParameters
	}
	nextToken := tokens[1]
	// trim the symbol as necessary based on type then call print
	if nextToken.Type == LiFloat ||
		nextToken.Type == LiInteger {
		return internalPrint(g, nextToken.Symbol)
	}
	if nextToken.Type == LiString {
		return internalPrint(g, nextToken.Symbol)
	}
	// if printing a string variable or string expression then try to evaluate it and print result
	if nextToken.Type == MaVariableString {
		value, err := parseStringExpression(g, tokens[1:])
		if err != 0 {
			return err
		} else {
			return internalPrint(g, value)
		}
	}
	// handle bad type for print here:
	return 0
}

func internalPrint(g *Game, text string) int {
	// PRINT command
	logMsg("internalPrint: " + text)
	g.Print(text)
	return 0
}
