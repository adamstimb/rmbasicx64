package nimgobus

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

// ValidateStyle validates if a style slot is valid
func (n *Nimbus) ValidateStyle(s int) bool {
	if s < 1 || s > 5 {
		return false
	}
	return true
}

// PlonkLogo draws the RM Nimbus logo
func (n *Nimbus) PlonkLogo(x, y int) {
	n.drawSprite(Sprite{pixels: n.logoImage, x: x + 1, y: y, colour: -1, over: true})
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
	resizedSprite := n.resizeSprite(Sprite{pixels: img, x: x, y: y, colour: opt.Brush, over: over}, imgWidth*opt.SizeX, imgHeight*opt.SizeY)
	rotatedSprite := n.rotateSprite(Sprite{pixels: resizedSprite.pixels, x: x, y: y, colour: opt.Brush, over: over}, opt.Direction)
	if newSprite, ok := n.applyDrawingbox(rotatedSprite, 0); ok {
		n.drawSprite(newSprite)
	}
}

// drawLine implements Bresenham's line algorithm to draw a line on a 2d array
// adapted from https://github.com/StephaneBunel/bresenham/blob/master/drawline.go
func (n *Nimbus) drawLine(img [][]int, x1, y1, x2, y2 int) [][]int {
	imgHeight := len(img) - 1
	var dx, dy, e, slope int

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
	if newSprite, ok := n.applyDrawingbox(Sprite{pixels: img, x: minX, y: minY, colour: opt.Brush, over: over}, 0); ok {
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
	Brush     int
	Over      int
	FillStyle FillStyle
}

// Circle draws a a filled circle
func (n *Nimbus) Circle(opt CircleOptions, r, x, y int) {
	// Handle default values
	if opt.Brush == -255 {
		opt.Brush = n.brush
	}
	var fillStyle FillStyle
	if opt.FillStyle.Style < 0 {
		fillStyle = n.AskFillStyle()
	} else {
		fillStyle = opt.FillStyle
	}
	// Handle Brush
	highestColour := 3
	if n.AskMode() == 40 {
		highestColour = 15
	}
	if opt.Brush < 128 {
		opt.Brush = overflow(opt.Brush, highestColour)
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
	if newSprite, ok := n.applyDrawingbox(Sprite{pixels: img, x: sx, y: sy, colour: opt.Brush, over: over, fillStyle: fillStyle}, 0); ok {
		n.drawSprite(newSprite)
	}
}

type AreaOptions struct {
	Brush     int
	Over      int
	FillStyle FillStyle
}

// Area draws a filled polygon of coordinates on the screen
func (n *Nimbus) Area(opt AreaOptions, coordList []XyCoord) {
	// Handle default values
	if opt.Brush == -255 {
		opt.Brush = n.brush
	}
	var fillStyle FillStyle
	if opt.FillStyle.Style < 0 {
		fillStyle = n.AskFillStyle()
	} else {
		fillStyle = opt.FillStyle
	}
	// Handle Brush
	highestColour := 3
	if n.AskMode() == 40 {
		highestColour = 15
	}
	if opt.Brush < 128 {
		opt.Brush = overflow(opt.Brush, highestColour)
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
	if newSprite, ok := n.applyDrawingbox(Sprite{pixels: img, x: minX, y: minY, colour: opt.Brush, over: over, fillStyle: fillStyle}, 0); ok {
		n.drawSprite(newSprite)
	}
}

type PointsOptions struct {
	Style int
	Brush int
	Over  int
}

// Points draws points at some given coordinates on the screen
func (n *Nimbus) Points(opt PointsOptions, coordList []XyCoord) {
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
	if opt.Style == -255 {
		opt.Style = n.pointsStyle
	}
	// Draw sprites
	for _, coord := range coordList {
		if newSprite, ok := n.applyDrawingbox(Sprite{pixels: n.pointsStyles[opt.Style-1], x: coord.X - 4, y: coord.Y - 4, colour: opt.Brush, over: over}, 0); ok {
			n.drawSprite(newSprite)
		}
	}
}

// Flood fill algorithm adapted from https://stackoverflow.com/questions/2783204/flood-fill-using-a-stack
func (n *Nimbus) floodFillDo(maxX int, hits [250][640]bool, x, y, srcColor, tgtColor int, useEdgeColour bool, edgeColour int, fillStyle FillStyle) bool {
	if (y < 0) || (x < 0) || (y > 249) || (x > maxX) {
		return false
	}
	if hits[y][x] {
		return false
	}
	if useEdgeColour {
		if n.videoMemory[249-y][x] == edgeColour {
			return false
		}
	} else {
		if n.videoMemory[249-y][x] != srcColor {
			return false
		}
	}
	// valid, paint it
	n.videoMemory[249-y][x] = n.handlePattern(x, y, tgtColor, fillStyle)
	return true
}

func (n *Nimbus) floodFill(x, y, color int, useEdgeColour bool, edgeColour int, fillStyle FillStyle) {
	maxX := 639
	if n.AskMode() == 40 {
		maxX = 319
	}
	srcColor := n.GetPixel(x, y)
	hits := [250][640]bool{}
	queue := []XyCoord{}
	queue = append(queue, XyCoord{x, y})
	n.muVideoMemory.Lock()
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		result := n.floodFillDo(maxX, hits, p.X, p.Y, srcColor, color, useEdgeColour, edgeColour, fillStyle)
		if result {
			hits[p.Y][p.X] = true
			queue = append(queue, XyCoord{p.X + 1, p.Y + 1})
			queue = append(queue, XyCoord{p.X - 1, p.Y - 1})
			queue = append(queue, XyCoord{p.X + 1, p.Y - 1})
			queue = append(queue, XyCoord{p.X - 1, p.Y + 1})
			queue = append(queue, XyCoord{p.X - 1, p.Y})
			queue = append(queue, XyCoord{p.X + 1, p.Y})
		}
	}
	n.muVideoMemory.Unlock()
}

type FloodOptions struct {
	Brush         int
	UseEdgeColour bool
	EdgeColour    int
	FillStyle     FillStyle
}

// Flood seeds a boundary fill at x, y
func (n *Nimbus) Flood(opt FloodOptions, coord XyCoord) {
	// Handle default values
	if opt.Brush == -255 {
		opt.Brush = n.brush
	}
	var fillStyle FillStyle
	if opt.FillStyle.Style < 0 {
		fillStyle = n.AskFillStyle()
	} else {
		fillStyle = opt.FillStyle
	}
	// Start fill
	n.floodFill(coord.X, coord.Y, opt.Brush, opt.UseEdgeColour, opt.EdgeColour, fillStyle)
}

// Clearblock resets all the image blocks
func (n *Nimbus) Clearblock() {
	n.imageBlocks = [100]imageBlock{}
	for i := range n.imageBlocks {
		n.imageBlocks[i].deleted = true
	}
}

// Fetch receives an image, downsamples the number of colours to 4 or 16 depending
// on current screen mode, and assigns it to a Nimbus image block
func (n *Nimbus) Fetch(b int, path string) bool {
	// Load image from disk
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Printf("Error loading %s: %v", path, err)
		return false
	}
	width, height := img.Size()

	// Make a temp pallete of the current screen colours
	rgbaColours := make([]color.Color, len(n.palette))
	for i := 0; i < len(n.palette); i++ {
		rgbaColours[i] = n.basicColours[n.palette[i]]
	}
	tempPalette := color.Palette(rgbaColours)

	// Replace the colour of each pixel in the image with the nearest Nimbus colour
	spriteImg := make2dArray(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nearestRgba := tempPalette.Convert(img.At(x, y))
			for c := 0; c < len(n.palette); c++ {
				if nearestRgba == tempPalette[c] {
					spriteImg[y][x] = c
					break
				}
			}
		}
	}

	// Store the image block
	n.imageBlocks[b] = imageBlock{image: spriteImg, mode: n.AskMode()}
	return true
}

// Readblock reads an area x1, y1, x2, y2 of the screen into block b
func (n *Nimbus) Readblock(b, x1, y1, x2, y2 int) {
	// Clamp x, y values to within screen
	maxX := 319
	if n.AskMode() == 80 {
		maxX = 639
	}
	if y1 < 0 {
		y1 = 0
	}
	if y1 > 249 {
		y1 = 249
	}
	if x1 < 0 {
		x1 = 0
	}
	if x2 > maxX {
		x2 = maxX
	}
	if y2 < 0 {
		y2 = 0
	}
	if y2 > 249 {
		y2 = 249
	}
	if x2 < 0 {
		x2 = 0
	}
	if x2 > maxX {
		x2 = maxX
	}
	// Define a 2d array to store the image data
	width := x2 - x1
	height := y2 - y1
	img := make2dArray(width, height)
	// Use GetPixel to allow the drawqueue to flush then grab the video memory lock
	_ = n.GetPixel(0, 0)
	n.muVideoMemory.Lock()
	// Copy the section, unlock, and bung it in the blocks
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			img[len(img)-(y-y1+1)][x-x1] = n.videoMemory[249-y][x]
		}
	}
	n.muVideoMemory.Unlock()
	n.imageBlocks[b] = imageBlock{image: img, mode: n.AskMode()}
}

// Writeblock draws an image block on the screen at position x, y
func (n *Nimbus) Writeblock(b, x, y int, over bool) {
	// Retrieve image block and draw it
	block := n.imageBlocks[b]
	if !block.deleted {
		n.drawSprite(Sprite{pixels: block.image, x: x, y: y, colour: -1, over: over})
	}
}

// AskBlocksize returns the width, height and original screen mode of an image block
func (n *Nimbus) AskBlocksize(b int) (width, height, mode int) {
	block := n.imageBlocks[b]
	if !block.deleted {
		mode = block.mode
		width = len(block.image[0])
		height = len(block.image)
	}
	// return all zeroes if block deleted
	return width, height, mode
}

// Squash is the same as Writeblock but scales the image by 1/4
func (n *Nimbus) Squash(b, x, y int, over bool) {
	// Retrieve image block, rescale and draw it
	block := n.imageBlocks[b]
	if !block.deleted {
		width, height, _ := n.AskBlocksize(b)
		newWidth := int(width / 4)
		newHeight := int(height / 4)
		rescaledImg := n.resizeSprite(Sprite{pixels: block.image}, newWidth, newHeight)
		n.drawSprite(Sprite{pixels: rescaledImg.pixels, x: x, y: y, colour: -1, over: over})
	}
}

// Delblock "deletes" an image block by setting it's deleted flag to true
func (n *Nimbus) Delblock(b int) {
	n.imageBlocks[b].deleted = true
}

// Keep saves an image block b with a specific format
func (n *Nimbus) Keep(b int, format, path string) error {
	block := n.imageBlocks[b]
	if block.deleted {
		return nil
	}
	// Turn the block into an image
	width, height, _ := n.AskBlocksize(b)
	img := image.NewRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{width, height}})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.SetRGBA(x, y, n.basicColours[n.palette[block.image[y][x]]])
		}
	}
	// Create the file then attempt to encode and write with the appropriate format
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	switch format {
	case "jpeg":
		if err = jpeg.Encode(f, img, nil); err != nil {
			log.Printf("failed to encode jpeg: %v", err)
			return err
		}
	case "png":
		if err = png.Encode(f, img); err != nil {
			log.Printf("failed to encode png: %v", err)
			return err
		}
	}
	return nil
}

// SetPattern defines the user-definable patterns
func (n *Nimbus) SetPattern(slot, row, c1, c2, c3, c4 int) {
	// clamp colours to maximum for current screen mode
	maxC := 4
	if n.AskMode() == 40 {
		maxC = 15
	}
	if c1 > maxC {
		c1 = maxC
	}
	if c2 > maxC {
		c2 = maxC
	}
	if c3 > maxC {
		c3 = maxC
	}
	if c4 > maxC {
		c4 = maxC
	}
	if c1 < 0 {
		c1 = 0
	}
	if c2 < 0 {
		c2 = 0
	}
	if c3 < 0 {
		c3 = 0
	}
	if c4 < 0 {
		c4 = 0
	}
	n.patterns[slot-128][row] = [4]int{c1, c2, c3, c4}
}

// SetFillStyle sets the fill style for AREA, CIRCLE, SLICE and FLOOD
func (n *Nimbus) SetFillStyle(style, hatching, colour2 int) {
	_ = n.GetPixel(0, 0)
	n.fillStyle = FillStyle{Style: style, Hatching: hatching, Colour2: colour2}
}

// AskFillStyle gets the current fill style
func (n *Nimbus) AskFillStyle() FillStyle {
	return n.fillStyle
}

// SetDrawing selects a drawingbox if only 1 parameter is passed (index), or
// defines a drawingbox if 5 parameters are passed (index, col1, row1, col2,
// row2)
func (n *Nimbus) SetDrawing(p ...int) {
	// Validate number of parameters
	if len(p) != 1 && len(p) != 5 {
		// invalid
		panic("SetDrawing accepts either 1 or 5 parameters")
	}
	if len(p) == 1 {
		// Select drawingbox - validate choice first then set it
		if p[0] < 0 || p[0] > 10 {
			panic("SetDrawing index out of range")
		}
		n.selectedDrawingBox = p[0]
		return
	}
	// Otherwise define textbox if index is not 0
	if p[0] == 0 {
		panic("SetDrawing cannot define index zero")
	}
	// Clamp x and y values (x > 0, x < screenwidth, y > 0, y < 250)
	screenWidth := 320
	if n.AskMode() == 80 {
		screenWidth = 640
	}
	// x values
	for i := 1; i < 5; i += 2 {
		if p[i] < 0 {
			p[i] = 0
			continue
		}
		if p[i] >= screenWidth {
			p[i] = screenWidth - 1
			continue
		}
	}
	// y values
	for i := 2; i < 5; i += 2 {
		if p[i] < 0 {
			p[i] = 0
			continue
		}
		if p[i] >= 250 {
			p[i] = 250 - 1
			continue
		}
	}
	// Set bottomLeft and topRight colrows
	var upper, lower, left, right int
	if p[1] < p[3] {
		left = p[1]
		right = p[3]
	} else {
		left = p[3]
		right = p[1]
	}
	if p[2] < p[4] {
		lower = p[2]
		upper = p[4]
	} else {
		lower = p[4]
		upper = p[2]
	}
	// Set drawingbox
	n.drawingBoxes[p[0]] = drawingBox{left, upper, right, lower}

	return
}
