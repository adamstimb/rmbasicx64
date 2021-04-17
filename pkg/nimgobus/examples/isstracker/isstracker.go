package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png" // import only for side-effects
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adamstimb/rmbasicx64/pkg/nimgobus"
	"github.com/adamstimb/rmbasicx64/pkg/nimgobus/examples/isstracker/issImages"
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
	SplashScreen(g)
	Track(g)
}

func SplashScreen(g *Game) {
	// Display a funky splash screen
	// Load images from string vars as recommended in ebiten docs
	img, _, err := image.Decode(bytes.NewReader(issImages.Iss))
	if err != nil {
		log.Fatal(err)
	}
	issImg := ebiten.NewImageFromImage(img)
	g.SetMode(40)
	g.Fetch(issImg, 1)
	g.Writeblock(1, 0, 0)
	plotOpts := nimgobus.PlotOptions{
		Brush: 14, Font: 1, SizeX: 1, SizeY: 2,
	}
	g.Plot(plotOpts, "ISS TRACKER", 230, 30)
	plotOpts.SizeY = 1
	g.Plot(plotOpts, "copyright (c) P.P. Bottoms-Farts 1986", 22, 20)
	g.Plot(plotOpts, "Fegg-Heyes Primary School, North Staffs", 8, 10)
	time.Sleep(3 * time.Second)
}

func getPosition() (float64, float64) {
	// Get current latitude and logitude of ISS
	r, _ := http.Get("http://api.open-notify.org/iss-now.json")
	type Position struct {
		Latitude  string
		Longitude string
	}
	type ApiBody struct {
		Timestamp    int
		Message      string
		Iss_position Position
	}
	var body ApiBody
	rawBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(string(rawBody)), &body)
	long, _ := strconv.ParseFloat(body.Iss_position.Longitude, 64)
	lat, _ := strconv.ParseFloat(body.Iss_position.Latitude, 64)
	return long, lat
}

func Track(g *Game) {
	// Display and update the tracking screen
	g.SetMode(80)
	g.SetColour(0, 9)
	g.SetColour(1, 1)
	g.SetColour(2, 2)
	g.SetBorder(0)
	g.SetPaper(0)
	g.SetCharset(1)
	g.SetPen(3)
	g.Cls()
	img, _, err := image.Decode(bytes.NewReader(issImages.World500x250))
	if err != nil {
		log.Fatal(err)
	}
	worldImg := ebiten.NewImageFromImage(img)
	g.Fetch(worldImg, 1)
	g.Writeblock(1, 0, 0)
	longScale := 500.0 / 360.0
	latScale := 250.0 / 180.0
	for {
		long, lat := getPosition()
		circleOpts := nimgobus.CircleOptions{
			Brush: 2,
		}
		x := int(long*longScale) + 250
		if x > 500 {
			x -= 500
		}
		y := 125 + int(lat*latScale)
		g.Circle(circleOpts, 6, x, y)
		circleOpts.Brush = 3
		g.Circle(circleOpts, 4, x, y)
		g.SetCurpos(66, 1)
		g.Print("Longitude:")
		g.SetCurpos(66, 2)
		g.Print(fmt.Sprintf("%f", long))
		g.SetCurpos(66, 4)
		g.Print("Latitude:")
		g.SetCurpos(66, 5)
		g.Print(fmt.Sprintf("%f", lat))
		time.Sleep(1 * time.Second)
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
