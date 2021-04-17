package nimgobus

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/StephaneBunel/bresenham"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PlonkLogo draws the RM Nimbus logo
func (n *Nimbus) PlonkLogo(x, y int) {
	// Convert position
	_, height := n.logoImage.Size()
	ex, ey := n.convertPos(x, y, height)

	// Draw the logo at the location
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(n.logoImage, op)
}

// PlotOptions describes optional parameters for the Plot command.  Plot will
// interpret zero values for SizeX and SizeY as 1.
type PlotOptions struct {
	Brush     int
	Font      int
	Direction int
	SizeX     int
	SizeY     int
}

// Plot draws a string of characters on the paper at a given location
// with the colour and size of your choice.
func (n *Nimbus) Plot(opt PlotOptions, text string, x, y int) {
	// Handle default size values
	if opt.SizeX == 0 {
		opt.SizeX = 1
	}
	if opt.SizeY == 0 {
		opt.SizeY = 1
	}
	// Validate brush
	n.validateColour(opt.Brush)
	// Create a new image big enough to contain the plotted chars
	// (without scaling)
	img := ebiten.NewImage(len(text)*10, 10)
	// draw chars on the image
	xpos := 0
	for _, c := range text {
		n.drawChar(img, int(c), xpos, 0, opt.Brush, opt.Font)
		xpos += 8
	}
	// Scale img and draw on paper
	ex, ey := n.convertPos(x, y, 10*opt.SizeY)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(opt.SizeX), float64(opt.SizeY))
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(img, op)
}

// LineOptions describes optional parameters for the Line command.
type LineOptions struct {
	Brush int
}

// Line draws connected lines on the screen.  x, y values are passed in the variadic
// p parameter.
func (n *Nimbus) Line(opt LineOptions, p ...int) {
	// Validate colour
	n.validateColour(opt.Brush)
	// Use drawLine to draw connected lines
	for i := 0; i < len(p)-2; i += 2 {
		n.drawLine(p[i], p[i+1], p[i+2], p[i+3], opt.Brush)
	}
}

// AreaOptions describes optional parameters for the Area command
type AreaOptions struct {
	Brush int
}

// Area draws a filled polygon on the screen.  x, y values are passed in the variadic
// p parameter.
func (n *Nimbus) Area(opt AreaOptions, p ...int) {
	// Validate colour
	n.validateColour(opt.Brush)
	// Use vector to draw the polygon
	var path vector.Path
	ex, ey := n.convertPos(p[0], p[1], 1)
	path.MoveTo(float32(ex), float32(ey)) // Go to start position
	for i := 2; i < len(p)-1; i += 2 {
		ex, ey = n.convertPos(p[i], p[i+1], 1)
		path.LineTo(float32(ex), float32(ey))
	}
	// Is the shape closed?  If not, draw a line back to start position
	if p[len(p)-2] != p[0] || p[len(p)-1] != p[1] {
		// Shape is open so close it
		ex, ey = n.convertPos(p[0], p[1], 1)
		path.MoveTo(float32(ex), float32(ey))
	}
	// Fill the shape on paper
	op := &vector.FillOptions{
		Color: n.convertColour(opt.Brush),
	}
	path.Fill(n.paper, op)
}

// CircleOptions describes optional parameters for the Circle command.
type CircleOptions struct {
	Brush int
}

// SliceOptions describes optional parameters for the Slice command.
type SliceOptions struct {
	Brush int
}

type xyCoord struct {
	x int
	y int
}

// drawCircle draws a Bresenham circle or slice and should only be called by either
// Circle or Slice functions.  To draw a full circle set startAngle and stopAngle
// both to 0.
func (n *Nimbus) drawCircle(opt CircleOptions, r, startAngle, stopAngle, xc, yc int) {
	// Handle slice or whole circle
	var slice bool
	if startAngle == 0 && stopAngle == 0 {
		slice = false
	} else {
		slice = true
	}
	// Convert co-ordinates
	ex, ey := n.convertPos(xc, yc, 1)
	xc = int(ex)
	yc = int(ey)
	// Calculate points and corresponding angle using Bresenham's algorithm
	x := 0
	y := r
	d := 3 - 2*r
	points := make(map[float64]xyCoord)
	points = addCirclePoints(points, xc, yc, x, y)
	for y >= x {
		x++
		if d > 0 {
			y--
			d = d + 4*(x-y) + 10
		} else {
			d = d + 4*x + 6
		}
		points = addCirclePoints(points, xc, yc, x, y)
	}
	// Draw the shape as a filled polygon
	op := &vector.FillOptions{
		Color: n.convertColour(opt.Brush),
	}
	var path vector.Path
	// Render a whole circle or slice
	if slice {
		// slice
		path = makeSliceVectors(points, xc, yc, startAngle, stopAngle)
	} else {
		// whole circle
		path = makeCircleVectors(points)
	}
	path.Fill(n.paper, op)
}

// Circle draws a circle of radius r with the center located at x, y.
func (n *Nimbus) Circle(opt CircleOptions, r, x, y int) {
	// Validate
	n.validateColour(opt.Brush)
	if r < 1 {
		panic("Radius is out of range")
	}
	// Delegate to drawCircle
	n.drawCircle(opt, r, 0, 0, x, y)
}

// Slice draws a slice of a circle of radius r with the center located at x, y.  The slice
// begins in degrees startAngle from the vertical (where vertical is 0 or 360 degrees) and
// ends at stopAngle.
func (n *Nimbus) Slice(opt SliceOptions, r, startAngle, stopAngle, x, y int) {
	// Validate
	n.validateColour(opt.Brush)
	if r < 1 {
		panic("Radius is out of range")
	}
	if startAngle < 0 || startAngle > 360 {
		panic("startAngle is out of range")
	}
	if stopAngle < 0 || stopAngle > 360 {
		panic("stopAngle is out of range")
	}
	if stopAngle == startAngle {
		panic("startAngle and stopAngle cannot be equal")
	}
	// Convert sliceoptions to circleoptions and delegate to drawCircle
	var circleOpts CircleOptions
	circleOpts.Brush = opt.Brush
	n.drawCircle(circleOpts, r, startAngle, stopAngle, x, y)
}

// makeSliceVectors takes point coorindates for a circle and creates vectors
// that describe a slice
func makeSliceVectors(points map[float64]xyCoord, xc, yc, startAngle, stopAngle int) vector.Path {
	var keys []float64
	var path vector.Path
	for k := range points {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	start := false
	drawing := true
	i := 0
	for drawing {
		k := keys[i]

		if start {
			fmt.Println(k, float32(points[k].x), float32(points[k].y))
			// slice has started so draw vectors around circle
			path.LineTo(float32(points[k].x), float32(points[k].y))
		}
		if k >= float64(startAngle) && start == false {
			// slice needs to start
			path.MoveTo(float32(xc), float32(yc))
			path.LineTo(float32(points[k].x), float32(points[k].y))
			start = true
		}
		if k >= float64(stopAngle) && start {
			// slice needs to stop
			path.LineTo(float32(points[k].x), float32(points[k].y))
			path.LineTo(float32(xc), float32(yc))
			drawing = false
		}
		// increment i and overflow if required
		i++
		if i >= len(keys) {
			// overflow
			i = 0
		}
	}
	return path
}

// makeCircleVectors takes point coordinates generated by the Bresenham
// algorithm and converts them into vectors
func makeCircleVectors(points map[float64]xyCoord) vector.Path {
	var keys []float64
	var path vector.Path
	for k := range points {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	start := true
	for _, k := range keys {
		if start {
			path.MoveTo(float32(points[k].x), float32(points[k].y))
			start = false
		} else {
			path.LineTo(float32(points[k].x), float32(points[k].y))
		}
	}
	return path
}

// addCirclePoints calculates 8 symmetrical points around a circle according
// to Bresenham.  The results are returned in the points map where the key
// value is the angle from vertical of the coorindate in degrees.
func addCirclePoints(points map[float64]xyCoord, xc, yc, x, y int) map[float64]xyCoord {
	var coords [8]xyCoord
	coords[0] = xyCoord{xc + x, yc + y}
	coords[1] = xyCoord{xc - x, yc + y}
	coords[2] = xyCoord{xc + x, yc - y}
	coords[3] = xyCoord{xc - x, yc - y}
	coords[4] = xyCoord{xc + y, yc + x}
	coords[5] = xyCoord{xc - y, yc + x}
	coords[6] = xyCoord{xc + y, yc - x}
	coords[7] = xyCoord{xc - y, yc - x}
	// To draw vectors we need to map each coordinate by it's angle from
	// vertical.  This involves calculating the angle of triangle (partial
	// angle) and then calculating the angle from vertical depending on
	// which quadrant the coordinate is in, or something...
	for _, coord := range coords {
		opp := math.Abs(float64(coord.y - yc))
		adj := math.Abs(float64(coord.x - xc))
		partialAngle := math.Atan(opp/adj) * 180 / math.Pi
		var angle float64
		if coord.y >= yc && coord.x >= xc {
			// top-right quadrant
			angle = 90 - partialAngle
		}
		if coord.y <= yc && coord.x >= xc {
			// bottom-right quadrant
			angle = 90 + partialAngle
		}
		if coord.y <= yc && coord.x <= xc {
			// bottom-left quadrant
			angle = 180 + (90 - partialAngle)
		}
		if coord.y >= yc && coord.x <= xc {
			// top-left quadrant
			angle = 270 + partialAngle
		}
		points[angle] = xyCoord{coord.x, coord.y}
	}
	return points
}

// drawLine uses the Bresenham algorithm to draw a straight line on the Nimbus paper
func (n *Nimbus) drawLine(x1, y1, x2, y2, colour int) {
	// convert coordinates
	ex1, ey1 := n.convertPos(x1, y1, 1)
	ex2, ey2 := n.convertPos(x2, y2, 1)
	// create a temp image on which to draw the line
	paperWidth, paperHeight := n.paper.Size()
	dest := image.NewRGBA(image.Rect(0, 0, paperWidth, paperHeight))
	bresenham.Bresenham(dest, int(ex1), int(ey1), int(ex2), int(ey2), n.convertColour(colour))
	// create a copy of the image as an ebiten.image and paste it on to the Nimbus paper
	img := ebiten.NewImageFromImage(dest)
	op := &ebiten.DrawImageOptions{}
	n.paper.DrawImage(img, op)
}

// Fetch receives an image, downsamples the number of colours to 4 or 16 depending
// on current screen mode, and assigns it to a Nimbus image block
func (n *Nimbus) Fetch(img *ebiten.Image, b int) {
	width, height := img.Size()
	newImg := ebiten.NewImage(width, height)

	// Make a temp pallete of the current screen colours
	rgbaColours := make([]color.Color, len(n.palette))
	for i := 0; i < len(n.palette); i++ {
		rgbaColours[i] = n.convertColour(i)
	}
	tempPalette := color.Palette(rgbaColours)

	// Replace the colour of each pixel in the image with the nearest Nimbus colour
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImg.Set(x, y, tempPalette.Convert(img.At(x, y)))
		}
	}
	// Store the image block
	n.imageBlocks[b] = newImg
}

// Writeblock draws an image block on the screen at position x, y
func (n *Nimbus) Writeblock(b, x, y int) {
	// Retrieve image block and draw it
	img := n.imageBlocks[b]
	_, height := img.Size()
	ex, ey := n.convertPos(x, y, height)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(img, op)
}
