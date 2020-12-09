package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/adamstimb/nimgobus"
	"github.com/elastic/go-sysinfo"
	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	count           int
	nimgobus.Nimbus // Embed the Nimbus in the Game struct
}

func NewGame() *Game {
	game := &Game{}
	game.Init() // Initialize Nimgobus
	return game
}

func (g *Game) Update() error {
	if g.count == 0 {
		go App(g) // Launch the Nimbus app on first iteration
	}
	g.count++
	g.Nimbus.Update() // Update the app on all subsequent iterations
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

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
	g.Print("RM BASICx64 version 0.00 9th December 2020")
	// Generate and print workspace available notification
	workspaceAvailable := fmt.Sprintf("%dG bytes workspace available", bToGb(memInfo.Available))
	g.Print(workspaceAvailable)
}

func App(g *Game) {
	// Draw welcome screen then start main interpreter loop
	welcomeScreen(g)
	for {
		// get raw console input
		rawInput := g.Input(":")
		if strings.ToUpper(rawInput) == "BYE" {
			// Exit with success code
			os.Exit(0)
		}
		if rawInput != "" {
			g.Print("Syntax error - well actually I don't know how to interpet BASIC yet!")
		}

	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the Nimbus monitor on the screen and scale to current window size.
	monitorWidth, monitorHeight := g.Monitor.Size()

	// Get ebiten window size so we can scale the Nimbus screen up or down
	// but if (0, 0) is returned we're not running on a desktop so don't do any scaling
	windowWidth, windowHeight := ebiten.WindowSize()

	// Calculate aspect ratios of Nimbus monitor and ebiten screen
	monitorRatio := float64(monitorWidth) / float64(monitorHeight)
	windowRatio := float64(windowWidth) / float64(windowHeight)

	// If windowRatio > monitorRatio then clamp monitorHeight to windowHeight otherwise
	// clamp monitorWidth to screenWidth
	var scale, offsetX, offsetY float64
	switch {
	case windowRatio > monitorRatio:
		scale = float64(windowHeight) / float64(monitorHeight)
		offsetX = (float64(windowWidth) - float64(monitorWidth)*scale) / 2
		offsetY = 0
	case windowRatio <= monitorRatio:
		scale = float64(windowWidth) / float64(monitorWidth)
		offsetX = 0
		offsetY = (float64(windowHeight) - float64(monitorHeight)*scale) / 2
	}

	// Apply scale and centre monitor on screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(offsetX, offsetY)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(g.Monitor, op)
}

func main() {
	// Set up resizeable window
	ebiten.SetWindowSize(1232, 1000)
	ebiten.SetWindowTitle("RM BASICx64")
	ebiten.SetWindowResizable(true)

	// Create a new game and pass it to RunGame method
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
