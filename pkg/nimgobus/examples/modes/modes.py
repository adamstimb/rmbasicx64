package main

import (
	_ "image/png" // import only for side-effects
	"log"
	"time"

	"github.com/adamstimb/rmbasicx64/pkg/nimgobus"
	"github.com/hajimehoshi/ebiten/v2"
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

func App(g *Game) {
	// This is the Nimbus app itself
	Mode80(g)
	time.Sleep(4 * time.Second)
	Mode40(g)
}

func Mode80(g *Game) {
	// Show a colour swatch and info about mode 80
	g.SetMode(80)
	plotOpts := nimgobus.PlotOptions{
		Brush: 2,
		SizeX: 2,
		SizeY: 2,
	}
	g.Plot(plotOpts, "Mode 80", 30, 220)
	plotOpts.Brush = 3
	g.Plot(plotOpts, "Mode 80", 32, 221)
	g.SetWriting(1, 5, 7, 75, 11)
	g.SetWriting(1)
	g.Print("This is Mode 80.  The screen is 80 character columns wide and 25 columns tall, or 640 pixels wide and 250 pixels tall.  Pixels are doubled in length along the vertical so everything has this wacky stretched-out look!  4 colours are available.")
	// Draw colour swatch
	areaOpts := nimgobus.AreaOptions{}
	var x, y int
	width := 143
	for i := 0; i < 4; i++ {
		areaOpts.Brush = i
		y = 50
		x = 30 + (i * width)
		areaOpts.Brush = 3
		g.Area(areaOpts, x-1, y-1, x+width+1, y-1, x+width+1, y+81, x-1, y+81, x-1, y-1)
		areaOpts.Brush = i
		g.Area(areaOpts, x, y, x+width, y, x+width, y+80, x, y+80, x, y)
	}
}

func Mode40(g *Game) {
	// Show a colour swatch and info about mode 40
	g.SetMode(40)
	g.SetPaper(1)
	g.SetBorder(1)
	g.Cls()
	plotOpts := nimgobus.PlotOptions{
		Brush: 0,
		SizeX: 2,
		SizeY: 2,
	}
	g.Plot(plotOpts, "Mode 40", 15, 220)
	plotOpts.Brush = 14
	g.Plot(plotOpts, "Mode 40", 16, 221)
	g.SetWriting(1, 3, 6, 38, 13)
	g.SetWriting(1)
	g.Print("This is Mode 40.  The screen is 40 character columns wide and 25 columns tall, or 320 pixels wide and 250 pixels tall.  16 sumptious colours are available.")
	// Draw colour swatch
	areaOpts := nimgobus.AreaOptions{}
	var x, y int
	width := 18
	for i := 0; i < 16; i++ {
		areaOpts.Brush = i
		y = 50
		x = 15 + (i * width)
		areaOpts.Brush = 15
		g.Area(areaOpts, x-1, y-1, x+width+1, y-1, x+width+1, y+81, x-1, y+81, x-1, y-1)
		areaOpts.Brush = i
		g.Area(areaOpts, x, y, x+width, y, x+width, y+80, x, y+80, x, y)
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
	ebiten.SetWindowSize(1400, 1000)
	ebiten.SetWindowTitle("Nimgobus")
	ebiten.SetWindowResizable(true)

	// Create a new game and pass it to RunGame method
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
