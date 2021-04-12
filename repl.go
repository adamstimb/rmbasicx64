package main

import (
	"fmt"
	"strings"

	"github.com/elastic/go-sysinfo"
)

// convert bytes to Gb
func bToGb(b uint64) uint64 {
	return b / 1024 / 1024 / 1024
}

// welcomeScreen draws the RM Basic welcome screen
func welcomeScreen(g *Game) {
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
	g.Print("RM BASICx64 Version 0.01 12th April 2021")
	// Generate and print workspace available notification
	workspaceAvailable := fmt.Sprintf("%dG bytes workspace available.", bToGb(memInfo.Available))
	g.Print(workspaceAvailable)
}

// editor is used to receive commands from the user in direct mode and edit
// BASIC programs
func repl(g *Game, i *Interpreter) {
	i.Init()
	for {
		rawInput := g.Input(":")
		code := strings.TrimSpace(rawInput)
		response := i.ImmediateInput(code)
		if response != "" {
			g.Print(response)
		}
	}
}

// StartUi is called when BASIC loads if an argument to immediately run a BASIC
// program has not been received.  It displays the welcome screen and starts
// the editor.
func StartUi(g *Game) {
	welcomeScreen(g)
	repl(g, &Interpreter{})
}
