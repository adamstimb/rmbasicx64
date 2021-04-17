package rmbasicx64

import (
	"fmt"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
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

// repl is the REPL that handles input
func repl(g *Game, i *Interpreter) {
	i.Init(g)
	for {
		rawInput := g.Input(":", "")
		code := strings.TrimSpace(rawInput)
		if !i.g.BreakInterruptDetected {
			// Don't execute if break detected
			_ = i.ImmediateInput(code)
		} else {
			// Might still have to print a message if <BREAK> occured while interpreter was at rest
			i.g.Print(syntaxerror.ErrorMessage(syntaxerror.InterruptedByBreakKey))
		}
		// Reset break flag
		i.g.BreakInterruptDetected = false
	}
}

// StartUi is called by the ebiten App.  It draws the welcome screen then starts the
// the REPL
func StartUi(g *Game) {
	welcomeScreen(g)
	repl(g, &Interpreter{})
}
