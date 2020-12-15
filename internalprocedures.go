package main

import "os"

func internalBye(g *Game) int {
	// Exit with success code
	logMsg("internalBye")
	os.Exit(0)
	return 0
}

func internalPrint(g *Game, text string) int {
	// PRINT command
	logMsg("internalPrint: " + text)
	g.Print(text)
	return 0
}
