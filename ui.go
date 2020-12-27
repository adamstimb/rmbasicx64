package main

import (
	"fmt"
	"log"

	"github.com/elastic/go-sysinfo"
)

// convert bytes to Gb
func bToGb(b uint64) uint64 {
	return b / 1024 / 1024 / 1024
}

const (
	// DEBUG will print log messages in the console if true
	DEBUG = true
)

func logMsg(msg string) {
	if DEBUG {
		log.Println(msg)
	}
}

// welcomeScreen draws the RM Basic welcome screen
func welcomeScreen(g *Game) {
	logMsg("welcomeScreen")
	// Collect system info
	host, err := sysinfo.Host()
	if err != nil {
		panic("Could not detect system information")
	}
	memInfo, err := host.Memory()
	// Draw welcome screen
	g.SetMode(80)
	g.PlonkLogo(0, 220)
	g.SetCurpos(1, 5)
	g.SetCursor(0)
	g.Print("This is a tribute project and is in no way linked to or endorsed by RM plc.")
	g.Print("")
	g.Print("RM BASICx64 Version 0.00 9th December 2020")
	// Generate and print workspace available notification
	workspaceAvailable := fmt.Sprintf("%dG bytes workspace available.", bToGb(memInfo.Available))
	g.Print(workspaceAvailable)
}

// raiseError prints an error message
func raiseError(g *Game, err int) {
	logMsg("RaiseError: " + errorMessages()[err])
	g.Print(errorMessages()[err])
}

// editor is used to receive commands from the user in direct mode and edit
// BASIC programs
func editor(g *Game) {
	logMsg("editor")
	// loop to prompt user for input and process that input
	for {
		// get raw console input for parsing
		rawInput := g.Input(":")
		logMsg("rawInput=" + rawInput)
		if rawInput != "" {
			s := &Scanner{}
			tokens := s.ScanTokens(rawInput)
			logMsg(rawInput)
			for _, token := range tokens {
				PrintToken(token)
			}
		}
	}
}

// StartUi is called when BASIC loads if an argument to immediately run a BASIC
// program has not been received.  It displays the welcome screen and starts
// the editor.
func StartUi(g *Game) {
	logMsg("StartUi")
	welcomeScreen(g)
	editor(g)
}
