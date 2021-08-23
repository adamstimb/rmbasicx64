package rmbasicx64

import (
	"fmt"
	"strings"
	"time"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/evaluator"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/game"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/lexer"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/object"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/parser"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/elastic/go-sysinfo"
)

// convert bytes to Gb
func bToGb(b uint64) uint64 {
	return b / 1024 / 1024 / 1024
}

// welcomeScreen draws the RM Basic welcome screen
func welcomeScreen(g *game.Game) {
	// Collect system info
	host, err := sysinfo.Host()
	if err != nil {
		panic("Could not detect system information")
	}
	memInfo, _ := host.Memory()
	// Draw welcome screen
	g.SetMode(80)
	g.PlonkLogo(0, 220)
	g.SetCurpos(1, 5)
	g.Print("This is a tribute project and is in no way linked to or endorsed by RM plc.")
	g.Put(13)
	g.Put(13)
	g.Print("RM BASICx64 Version 0.01B 21st July 2021")
	g.Put(13)
	// Generate and print workspace available notification
	workspaceAvailable := fmt.Sprintf("%dG bytes workspace available.", bToGb(memInfo.Available))
	g.Print(workspaceAvailable)
	g.Put(13)
}

// repl is the REPL that handles input
func repl(g *game.Game) {
	l := &lexer.Lexer{}
	env := object.NewEnvironment()
	//opt := nimgobus.PlotOptions{Brush: 2, SizeX: 1, SizeY: 1, Over: -1, Font: 1}
	//lastTPS := 0
	for {
		//opt.Brush = 0
		//g.Plot(opt, fmt.Sprintf("TPS: %d", lastTPS), 0, 239)
		//opt.Brush = 2
		//lastTPS := g.GetTPS()
		//g.Plot(opt, fmt.Sprintf("TPS: %d", lastTPS), 0, 239)
		g.Print(":")
		rawInput := g.Input("")
		code := strings.TrimSpace(rawInput)
		if !g.BreakInterruptDetected {
			// Don't execute if break detected
			l.Scan(code)
			p := parser.New(l, g)
			line := p.ParseLine()
			// Check for parser errors here.  Parser errors are handled just like evaluation errors but
			// obviously we'll skip evaluation if parsing already failed.
			if errorMsg, hasError := p.GetError(); hasError {
				g.Print(errorMsg)
				g.Put(13)
				p.JumpToToken(0)
				g.Print(p.PrettyPrint())
				g.Put(13)
				continue
			}
			// And this is temporary while we're still migrating from Monkey to RM Basic
			if len(p.Errors()) > 0 {
				g.Print("Oops! Some random parsing error occurred. These will be handled properly downstream by for now here's some spewage:")
				g.Put(13)
				p.JumpToToken(0)
				g.Print(p.PrettyPrint())
				g.Put(13)
				for _, msg := range p.Errors() {
					g.Print(msg)
					g.Put(13)
				}
				continue
			}
			// Add new line to stored program
			if line.Statements == nil {
				env.Program.AddLine(line.LineNumber, line.LineString)
				continue
			}
			// Execute each statement in the inputted line.  If an error occurs, print the
			// error message and stop.
			for statementNumber, stmt := range line.Statements {
				env.Program.CurrentStatementNumber = statementNumber
				obj := evaluator.Eval(g, stmt, env)
				if errorMsg, ok := obj.(*object.Error); ok {
					if errorMsg.ErrorTokenIndex != 0 {
						p.ErrorTokenIndex = errorMsg.ErrorTokenIndex
					}
					g.Print(fmt.Sprintf(errorMsg.Message))
					g.Put(13)
					p.JumpToToken(0)
					g.Print(p.PrettyPrint())
					g.Put(13)
					break
				}
			}
		} else {
			// Might still have to print a message if <BREAK> occurred while interpreter was at rest
			g.Print(syntaxerror.ErrorMessage(syntaxerror.InterruptedByBreakKey))
			g.Put(13)
			time.Sleep(150 * time.Millisecond)
		}
		// Reset break flag
		g.BreakInterruptDetected = false
	}
}

// StartUi is called by the ebiten App.  It draws the welcome screen then starts the
// the REPL.
func StartUi(g *game.Game) {
	if g.Config.Boot {
		g.Boot()
	}
	g.PrettyPrintIndent = ""
	welcomeScreen(g)
	repl(g)
}
