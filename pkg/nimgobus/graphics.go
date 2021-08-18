package nimgobus

// PlonkLogo draws the RM Nimbus logo
func (n *Nimbus) PlonkLogo(x, y int) {
	n.drawSprite(Sprite{n.logoImage, x, y, -1, true})
}

type PlotOptions struct {
	Brush     int
	Font      int
	Direction int
	SizeX     int
	SizeY     int
}

// Plot draws a string of characters on the paper at a given location
// with the colour, size and orientation of your choice.
func (n *Nimbus) Plot(opt PlotOptions, text string, x, y int) {
	// Handle default size values
	if opt.SizeX == 0 {
		opt.SizeX = 1
	}
	if opt.SizeY == 0 {
		opt.SizeY = 1
	}
	// TODO: Handle fonts
	// Validate brush
	// n.validateColour(opt.Brush)  // TODO: Decide once and for all how to handle this.
	// Create a new image big enough to contain the plotted chars
	// (without scaling)
	imgWidth := len(text) * 8
	imgHeight := 10
	img := make2dArray(imgWidth, imgHeight)
	// Select charset and draw chars on image
	xOffset := 0
	for _, c := range text {
		// draw char on image
		var charPixels [][]int
		switch opt.Font {
		case 0:
			charPixels = n.charImages0[c]
		case 1:
			charPixels = n.charImages1[c]
		}
		for x := 0; x < 8; x++ {
			for y := 0; y < 10; y++ {
				img[y][x+xOffset] = charPixels[y][x]
			}
		}
		xOffset += 8
	}
	// TODO: Stretch
	resizedSprite := n.resizeSprite(Sprite{img, x, y, opt.Brush, true}, imgWidth*opt.SizeX, imgHeight*opt.SizeY)
	// TODO: Direction
	// TODO: Over
	n.drawSprite(resizedSprite)
}
