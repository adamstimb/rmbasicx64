package nimgobus

// ValidateColour validates if a colour/palette slot is valid for the current screen mode
func (n *Nimbus) ValidateColour(c int) bool {
	maxC := 15
	if n.mode == 80 {
		maxC = 3
	}
	if c < 0 || c > maxC {
		return false
	}
	return true
}

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
	if newSprite, ok := n.applyDrawingbox(rotatedSprite, 0); ok {
		n.drawSprite(newSprite)
	}
	//n.drawSprite(n.applyDrawingbox(rotatedSprite, 0))
}

// drawLine implements Bresenham's line algorithm to draw a line on a 2d array
// adapted from https://github.com/StephaneBunel/bresenham/blob/master/drawline.go
func (n *Nimbus) drawLine(img [][]int, x1, y1, x2, y2 int) [][]int {
	imgHeight := len(img) - 1
	var dx, dy, e, slope int

	//log.Printf("drawLine imgHeight=%d (%d, %d)-(%d, %d)", imgHeight, x1, y1, x2, y2)

	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		img[imgHeight-y1][x1] = 1

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			img[imgHeight-y1][x1] = 1
			x1++
		}
		img[imgHeight-y1][x1] = 1

	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			img[imgHeight-y1][x1] = 1
			y1++
		}
		img[imgHeight-y1][x1] = 1

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				img[imgHeight-y1][x1] = 1
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				img[imgHeight-y1][x1] = 1
				x1++
				y1--
			}
		}
		img[imgHeight-y1][x1] = 1

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			// BresenhamDxXRYD(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				img[imgHeight-y1][x1] = 1
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			// BresenhamDxXRYU(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				img[imgHeight-y1][x1] = 1
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		img[imgHeight-y1][x1] = 1

	// higher than wide.
	default:
		if y1 < y2 {
			// BresenhamDyXRYD(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				img[imgHeight-y1][x1] = 1
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			// BresenhamDyXRYU(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				img[imgHeight-y1][x1] = 1
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		img[imgHeight-y1][x1] = 1
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
	imgWidth := (maxX - minX) + 1
	imgHeight := (maxY - minY) + 1
	img := make2dArray(imgWidth, imgHeight)
	// draw lines
	for i := 0; i < len(coordList)-1; i++ {
		//log.Printf("i=%d minXY=(%d, %d) line=(%d, %d)-(%d-%d)", i, minX, minY, coordList[i].X, coordList[i].Y, coordList[i+1].X, coordList[i+1].Y)
		img = n.drawLine(img, coordList[i].X-minX, coordList[i].Y-minY, coordList[i+1].X-minX, coordList[i+1].Y-minY)
	}
	if newSprite, ok := n.applyDrawingbox(Sprite{img, minX, minY, opt.Brush, over}, 0); ok {
		n.drawSprite(newSprite)
	}
	//n.drawSprite(n.applyDrawingbox(Sprite{img, minX, minY, opt.Brush, over}, 0))
}

// Draw implements a "Bresenham-ish" circle drawing algorithm adapted from
// https://github.com/benhoyt/circle/blob/master/circle.go
func (n *Nimbus) drawCircle(r int, img [][]int) [][]int {
	x := 0
	y := r
	xsq := 0
	rsq := r * r
	ysq := rsq
	// Loop x from 0 to the line x==y. Start y at r and each time
	// around the loop either keep it the same or decrement it.
	for x <= y {
		img[r+y][r+x] = 1
		img[r+x][r+y] = 1
		img[r+y][r-x] = 1
		img[r+x][r-y] = 1
		img[r-y][r+x] = 1
		img[r-x][r+y] = 1
		img[r-y][r-x] = 1
		img[r-x][r-y] = 1

		// New x^2 = (x+1)^2 = x^2 + 2x + 1
		xsq = xsq + 2*x + 1
		x++
		// Potential new y^2 = (y-1)^2 = y^2 - 2y + 1
		y1sq := ysq - 2*y + 1
		// Choose y or y-1, whichever gives smallest error
		a := xsq + ysq
		b := xsq + y1sq
		if a-rsq >= rsq-b {
			y--
			ysq = y1sq
		}
	}
	return img
}

type CircleOptions struct {
	Brush int
	Over  int
}

// Circle draws a a filled circle
func (n *Nimbus) Circle(opt CircleOptions, r, x, y int) {
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
	img := make2dArray((2*r)+10, (2*r)+1)
	sx := x - r
	sy := y - r
	// draw circle outline
	img = n.drawCircle(r, img)
	// Fill the shape using odd/even counter to determine if inside or outside shape.
	imgHeight := len(img)
	imgWidth := len(img[0])
	for y := 0; y < imgHeight; y++ {
		onLeftEdge := false
		leavingShape := false
		inside := false
		startX := 0
		for x := 0; x < imgWidth-1; x++ {
			// reset flags if beyond right-hand edge
			if leavingShape && img[y][x] == 0 {
				onLeftEdge = false
				leavingShape = false
				inside = false
				continue
			}
			// detect left-hand edge
			if !leavingShape && !inside && !onLeftEdge && img[y][x] == 1 {
				// found edge
				onLeftEdge = true
				inside = false
				continue
			}
			// detect entering shape
			if onLeftEdge && img[y][x] == 0 {
				onLeftEdge = false
				inside = true
				startX = x
			}
			// detect leaving shape
			if inside && img[y][x] == 1 {
				onLeftEdge = false
				inside = false
				leavingShape = true
				continue
			}
			// fill if inside shape
			if inside {
				img[y][x] = 1
			}
		}
		// Undo if we are still inside the shape at this position, because that's impossible
		// and probably due to apex of straight line.
		if inside {
			for x := startX; x <= imgWidth-1; x++ {
				img[y][x] = 0
			}
		}
	}
	if newSprite, ok := n.applyDrawingbox(Sprite{img, sx, sy, opt.Brush, over}, 0); ok {
		n.drawSprite(newSprite)
	}
	//n.drawSprite(n.applyDrawingbox(Sprite{img, sx, sy, opt.Brush, over}, 0))
}

type AreaOptions struct {
	Brush int
	Over  int
}

// Area draws a filled polygon of coordinates on the screen
func (n *Nimbus) Area(opt AreaOptions, coordList []XyCoord) {
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
	imgWidth := (maxX - minX) + 10 // deliberately padded for fill algorithm
	imgHeight := (maxY - minY) + 1
	img := make2dArray(imgWidth, imgHeight)
	// draw lines *and close the shape if required*
	if coordList[0].X != coordList[len(coordList)-1].X || coordList[0].Y != coordList[len(coordList)-1].Y {
		// shape if open so we need to close it
		coordList = append(coordList, XyCoord{coordList[0].X, coordList[0].Y})
	}
	for i := 0; i < len(coordList)-1; i++ {
		//log.Printf("i=%d minXY=(%d, %d) line=(%d, %d)-(%d-%d)", i, minX, minY, coordList[i].X, coordList[i].Y, coordList[i+1].X, coordList[i+1].Y)
		img = n.drawLine(img, coordList[i].X-minX, coordList[i].Y-minY, coordList[i+1].X-minX, coordList[i+1].Y-minY)
	}
	// Attempt to fill the shape using odd/even counter to determine if inside or outside shape. If the scanner
	// reaches the outer edge and the count is odd then fail with "too complicated" error.
	for y := 0; y < imgHeight; y++ {
		onLeftEdge := false
		leavingShape := false
		inside := false
		startX := 0
		for x := 0; x < imgWidth-1; x++ {
			// reset flags if beyond right-hand edge
			if leavingShape && img[y][x] == 0 {
				onLeftEdge = false
				leavingShape = false
				inside = false
				continue
			}
			// detect left-hand edge
			if !leavingShape && !inside && !onLeftEdge && img[y][x] == 1 {
				// found edge
				onLeftEdge = true
				inside = false
				continue
			}
			// detect entering shape
			if onLeftEdge && img[y][x] == 0 {
				onLeftEdge = false
				inside = true
				startX = x
			}
			// detect leaving shape
			if inside && img[y][x] == 1 {
				onLeftEdge = false
				inside = false
				leavingShape = true
				continue
			}
			// fill if inside shape
			if inside {
				img[y][x] = 1
			}
		}
		// Undo if we are still inside the shape at this position, because that's impossible
		// and probably due to apex of straight line.
		if inside {
			for x := startX; x <= imgWidth-1; x++ {
				img[y][x] = 0
			}
		}
	}
	if newSprite, ok := n.applyDrawingbox(Sprite{img, minX, minY, opt.Brush, over}, 0); ok {
		n.drawSprite(newSprite)
	}
	//n.drawSprite(n.applyDrawingbox(Sprite{img, minX, minY, opt.Brush, over}, 0))
}
