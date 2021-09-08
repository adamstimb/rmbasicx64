package rmbasicx64

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/game"
	"github.com/adamstimb/rmbasicx64/pkg/nimgobus/resources/icon"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewGame() *game.Game {
	g := &game.Game{}
	g.Init()
	g.LoadConfig()
	g.EnsureWorkspace()
	go App(g)
	return g
}

func App(g *game.Game) {
	log.SetOutput(os.Stdout)
	StartUi(g)
}

// StartRepl sets up the application Window and initializes nimgobus/ebiten
func StartRepl() {
	// Set up resizeable window and icon
	ebiten.SetWindowSize(1260, 1000)
	ebiten.SetWindowTitle("RM BASICx64")
	ebiten.SetWindowResizable(true)
	iconImg, _, err := image.Decode(bytes.NewReader(icon.Rmbasicx64_ico_48_png))
	if err != nil {
		log.Printf("Failed to read application icon - using default GLFW icon instead")
	}
	ebiten.SetWindowIcon([]image.Image{iconImg})
	// Create a new game and pass it to RunGame method
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
