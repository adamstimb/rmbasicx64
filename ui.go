package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-sysinfo"
)

// convert bytes to Gb
func bToGb(b uint64) uint64 {
	return b / 1024 / 1024 / 1024
}

// welcomeScreen draws the RM Basic welcome screen
func WelcomeScreen(g *Game) {
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
	g.Print("RM Basicx64 version 0.00 9th December 2020")
	// Generate and print workspace available notification
	workspaceAvailable := fmt.Sprintf("%dG bytes workspace available", bToGb(memInfo.Available))
	g.Print(workspaceAvailable)
}

// Editor is used to receive commands from the user in direct mode and edit
// BASIC programs
func Editor(g *Game) {
	// loop to prompt user for input and process that input
	for {
		// get raw console input
		rawInput := g.Input(":")
		if strings.ToUpper(rawInput) == "BYE" {
			// Exit with success code
			os.Exit(0)
		}
		if rawInput != "" {
			formatted := Format(rawInput)
			g.Print(formatted)
		}

	}
}

// StartUi is called when BASIC loads if an argument to immediately run a BASIC
// program has not been received.  It displays the welcome screen and starts
// the editor.
func StartUi(g *Game) {
	WelcomeScreen(g)
	Editor(g)
}
