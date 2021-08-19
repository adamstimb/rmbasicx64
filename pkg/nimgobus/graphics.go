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
	Over      int
}

// Plot draws a string of characters on the paper at a given location
// with the colour, size and orientation of your choice.
func (n *Nimbus) Plot(opt PlotOptions, text string, x, y int) {
	// Handle default values
	if opt.SizeX == -255 {
		opt.SizeX = n.plotSizeX
	}
	if opt.SizeY == -255 {
		opt.SizeY = n.plotSizeY
	}
	if opt.Brush == -255 {
		opt.Brush = n.brush
	}
	if opt.Direction == -255 {
		opt.Direction = n.plotDirection
	}
	if opt.Font == -255 {
		opt.Font = n.plotFont
	}
	var over bool
	switch opt.Over {
	case -255:
		over = n.over
	case 0:
		over = false
	case -1:
		over = true
	}
	// Validate brush
	// n.validateColour(opt.Brush)  // TODO: Decide once and for all how to handle this.
	// Plot chars and applying scaling/direction
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
	resizedSprite := n.resizeSprite(Sprite{img, x, y, opt.Brush, over}, imgWidth*opt.SizeX, imgHeight*opt.SizeY)
	rotatedSprite := n.rotateSprite(Sprite{resizedSprite.pixels, x, y, opt.Brush, over}, opt.Direction)
	n.drawSprite(n.applyDrawingbox(rotatedSprite, 0))
}

// drawLine implements Bresenham's line algorithm to draw a line on a 2d array
func (n *Nimbus) drawLine(img [][]int, x0, y0, x1, y1 int) [][]int {
	dx := x1 - x0
	dy := y1 - y0
	x := x0
	y := y0
	p := 2*dy - dx
	imgHeight := len(img)
	for x < x1 {
		if p >= 0 {
			img[(imgHeight-1)-y][x] = 1
			y++
			p = p + 2*dy - 2*dx
		} else {
			img[(imgHeight-1)-y][x] = 1
			p = p + 2*dy
		}
		x++
	}
	return img
}

type XyCoord struct {
	X int
	Y int
}

type LineOptions struct {
	Brush     int
	Font      int
	Direction int
	SizeX     int
	SizeY     int
	Over      int
}

// Line draws a list of coordinates on the screen connected by lines
func (n *Nimbus) Line(opt LineOptions, coordList []XyCoord) {
	// Handle default values
	if opt.Brush == -255 {
		opt.Brush = n.brush
	}
	var over bool
	switch opt.Over {
	case -255:
		over = n.over
	case 0:
		over = false
	case -1:
		over = true
	}
	// Find optimal image size and minimum x,y
	minX := 1000
	maxX := 0
	minY := 1000
	maxY := 0
	for _, coord := range coordList {
		if coord.X < minX {
			minX = coord.X
		}
		if coord.X > maxX {
			maxX = coord.X
		}
		if coord.Y < minY {
			minY = coord.Y
		}
		if coord.Y > maxY {
			maxY = coord.Y
		}
	}
	imgWidth := maxX - minX
	imgHeight := maxY - minY
	img := make2dArray(imgWidth, imgHeight)
	// draw lines
	for i := 0; i < len(coordList)-1; i++ {
		img = n.drawLine(img, coordList[i].X-minX, coordList[i].Y-minY, coordList[i+1].X-minX, coordList[i+1].Y-minY)
	}
	n.drawSprite(n.applyDrawingbox(Sprite{img, minX, minY, opt.Brush, over}, 0))
}
