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
	if nextToken.Type == MaVariableString {
		return internalPrint(g, nextToken.Symbol[1:len(nextToken.Symbol)-1])
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
