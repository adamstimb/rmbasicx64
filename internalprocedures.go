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
	result, err := parseStringExpression(g, tokens[1:])
	if err != 0 {
		return err
	} else {
		return internalPrint(g, result)
	}
}

func internalPrint(g *Game, text string) int {
	// PRINT command
	logMsg("internalPrint: " + text)
	g.Print(text)
	return 0
}
