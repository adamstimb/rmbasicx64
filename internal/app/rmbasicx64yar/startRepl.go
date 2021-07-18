package rmbasicx64yar

import (
	_ "image/png"
	"log"
	"os"

	"github.com/adamstimb/rmbasicx64yar/internal/app/rmbasicx64yar/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewGame() *game.Game {
	g := &game.Game{}
	g.Init()  // Initialize Nimgobus
	go App(g) // Start the app
	return g
}

func App(g *game.Game) {
	log.SetOutput(os.Stdout)
	StartUi(g)
}

// StartRepl sets up the application Window and initializes nimgobus/ebiten
func StartRepl() {
	// Set up resizeable window
	ebiten.SetWindowSize(1260, 1000)
	ebiten.SetWindowTitle("RM BASICx64")
	ebiten.SetWindowResizable(true)
	// Create a new game and pass it to RunGame method
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
