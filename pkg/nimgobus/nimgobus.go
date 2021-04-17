package nimgobus

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/adamstimb/rmbasicx64/pkg/nimgobus/images"
	"github.com/hajimehoshi/ebiten/v2"
)

// colRow defines a column, row position
type colRow struct {
	col int
	row int
}

// textBox defines the bounding box of a scrollable text box
type textBox struct {
	col1 int
	row1 int
	col2 int
	row2 int
}

// StoreItem defines an item in the Nimgobus store, which is a map of StoreItem arrays
// that acts as a simple RAM.  Values are stored as strings but the original data
// type is remembered by assigning a code number to valueType.
type StoreItem struct {
	valueType int
	value     string
}

// Nimbus acts as a container for all the components of the Nimbus monitor.  You
// only need to call the Init() method after declaring a new Nimbus.
type Nimbus struct {
	Store                  map[string][]StoreItem
	Monitor                *ebiten.Image
	borderImage            *ebiten.Image
	paper                  *ebiten.Image
	basicColours           []color.RGBA
	borderSize             int
	borderColour           int
	paperColour            int
	penColour              int
	charset                int
	cursorChar             int
	defaultHighResPalette  []int
	defaultLowResPalette   []int
	palette                []int
	logoImage              *ebiten.Image
	textBoxes              [10]textBox
	imageBlocks            [16]*ebiten.Image
	selectedTextBox        int
	cursorPosition         colRow
	cursorMode             int
	cursorCharset          int
	cursorFlash            bool
	charImages0            [256]*ebiten.Image
	charImages1            [256]*ebiten.Image
	keyBuffer              []int
	keyBufferLock          bool
	ebitenInputChars       []rune
	ebitenInputCharsLock   bool
	BreakInterruptDetected bool
}

// Init initializes a new Nimbus.  You must call this method after declaring a
// new Nimbus variable.
func (n *Nimbus) Init() {
	// in case any randomonia is required we can run a seed on startup
	rand.Seed(time.Now().UnixNano())

	// should next exceed 60tps
	ebiten.SetMaxTPS(60)

	// Load Nimbus logo image and both charsets
	n.loadLogoImage()
	n.loadCharsetImages(0)
	n.loadCharsetImages(1)
	// Set init values of everything else
	n.borderSize = 50
	n.borderImage = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
	n.Monitor = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
	n.paper = ebiten.NewImage(640, 250)
	n.basicColours = basicColours
	n.defaultHighResPalette = defaultHighResPalette
	n.defaultLowResPalette = defaultLowResPalette
	n.palette = defaultHighResPalette
	n.borderColour = 0
	n.paperColour = 0
	n.penColour = 3
	n.charset = 0
	n.cursorMode = -1
	n.cursorChar = 95
	n.cursorCharset = 0
	n.cursorPosition = colRow{1, 1}
	n.cursorFlash = false
	n.selectedTextBox = 0
	n.keyBuffer = []int{}
	n.keyBufferLock = false
	n.ebitenInputChars = []rune{}
	n.ebitenInputCharsLock = true

	// Initialize with mode 80 textboxes
	for i := 0; i < 10; i++ {
		n.textBoxes[i] = textBox{1, 1, 25, 80}
	}

	// Set break int
	n.BreakInterruptDetected = false

	// Start flashCursor
	go n.flashCursor()
	// Start keyboardMonitor
	go n.keyboardMonitor()
}

// flashCursor flips the cursorFlash flag every half second
func (n *Nimbus) flashCursor() {
	for {
		time.Sleep(500 * time.Millisecond)
		if n.cursorMode == 0 {
			// Flashing cursor
			n.cursorFlash = !n.cursorFlash
		}
		if n.cursorMode < 0 {
			// Invisible cursor
			n.cursorFlash = false
		}
		if n.cursorMode > 1 {
			// Visible cursor but not flashing
			n.cursorFlash = true
		}
	}
}

// pushKeyBuffer obtains keyBufferLock then pushes a new char to the buffer
func (n *Nimbus) pushKeyBuffer(char int) {
	// Wait to obtain keyBufferLock
	for n.keyBufferLock {
		//
	}
	// keyBufferLock released so obtain the lock and push new char
	n.keyBufferLock = true
	n.keyBuffer = append(n.keyBuffer, char)
	// all done release lock
	n.keyBufferLock = false
}

// popKeyBuffer obtains keyBufferLock then pops the oldest char in the buffer
// and returns the char.  If the buffer is empty -1 is returned.
func (n *Nimbus) popKeyBuffer() int {
	// check if buffer is empty and return -1 if so
	if len(n.keyBuffer) == 0 {
		// is empty
		return -1
	}
	// Otherwise wait to obtain keyBufferLock
	for n.keyBufferLock {
		//
	}
	// keyBufferLock released so obtain the lock and pop the buffer
	n.keyBufferLock = true
	char := n.keyBuffer[0]
	// if buffer only has one char re-initialize it otherwise shorten it by 1 element
	if len(n.keyBuffer) <= 1 {
		n.keyBuffer = []int{}
	} else {
		n.keyBuffer = n.keyBuffer[1:]
	}
	// all done release lock and return char
	n.keyBufferLock = false
	return char
}

// keyboardMonitor checks for key presses, handles repeating keys and adds to keyBuffer
func (n *Nimbus) keyboardMonitor() {

	handleRepeatingChars := func(keyChars []rune, lastChar int, repeatCount int) (int, int) {
		sleepTime := 10 * time.Microsecond
		repeatThreshold := 5
		for _, thisRune := range keyChars {
			char := int(thisRune)
			if char == lastChar && repeatCount < repeatThreshold {
				// Is repeating char
				// let use hold key down for several frames to prevent
				// spewing the same letter
				repeatCount++
				lastChar = char
				time.Sleep(sleepTime)
				continue
			} else {
				// Key has been held down long enough to be repeated or
				// it's not a repeating key so push it to buffer
				repeatCount = 0
				n.pushKeyBuffer(char)
				lastChar = char
			}
		}
		return lastChar, repeatCount
	}

	lastChar := -1
	repeatCount := 0
	for {
		// evaluate printable chars
		for !n.ebitenInputCharsLock {
			// wait for unlock
		}
		keyChars := n.ebitenInputChars
		n.ebitenInputCharsLock = false
		if len(keyChars) > 0 {
			lastChar, repeatCount = handleRepeatingChars(keyChars, lastChar, repeatCount)
		} else {
			// no printable char so evaluate control keys
			if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyKPEnter) {
				lastChar, repeatCount = handleRepeatingChars([]rune{-11}, lastChar, repeatCount)
				continue
			}
			if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
				lastChar, repeatCount = handleRepeatingChars([]rune{-10}, lastChar, repeatCount)
				continue
			}
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				lastChar, repeatCount = handleRepeatingChars([]rune{-12}, lastChar, repeatCount)
				continue
			}
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				lastChar, repeatCount = handleRepeatingChars([]rune{-13}, lastChar, repeatCount)
				continue
			}
			// handle BREAK interrupt (Ctrl+ScrollLock)
			if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyScrollLock) {
				n.BreakInterruptDetected = true
				continue
			}
			// if not control chars then next iteration
			repeatCount = 0
			lastChar = -1
			continue
		}
	}
}

func updateInputChars(n *Nimbus) {
	n.ebitenInputCharsLock = true
	n.ebitenInputChars = ebiten.InputChars()
}

// Update redraws the Nimbus monitor image
func (n *Nimbus) Update() {

	// Update input chars
	go updateInputChars(n)

	// Copy paper so we can apply overlays (e.g. cursor)
	//paperCopy := ebiten.NewImageFromImage(n.paper) // <---- Seems to break on Windows and is slow on Linux
	//...therefore we create a new image and paste the current paper on that
	paperCopy := ebiten.NewImage(n.paper.Size())
	op := &ebiten.DrawImageOptions{}
	paperCopy.DrawImage(n.paper, op)

	// Apply overlays
	// Cursor
	if n.cursorFlash {
		curX, curY := n.convertColRow(n.cursorPosition)
		n.drawChar(paperCopy, n.cursorChar, int(curX), int(curY), n.penColour, n.cursorCharset)
	}

	// calculate y scale for paper and draw it on border image
	paperX, paperY := paperCopy.Size()
	scaleX := 640.0 / float64(paperX)
	scaleY := 500.0 / float64(paperY)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(n.borderSize), float64(n.borderSize))
	n.borderImage.DrawImage(paperCopy, op)

	// Draw border image on monitor
	op = &ebiten.DrawImageOptions{}
	n.Monitor.DrawImage(n.borderImage, op)
}

// loadLogoImage loads the Nimbus logo image
func (n *Nimbus) loadLogoImage() {
	img, _, err := image.Decode(bytes.NewReader(images.NimbusLogoImage))
	if err != nil {
		log.Fatal(err)
	}
	n.logoImage = ebiten.NewImageFromImage(img)
}

// loadCharsetImages loads the charset images
func (n *Nimbus) loadCharsetImages(charset int) {
	var img image.Image
	var err error
	if charset == 0 {
		img, _, err = image.Decode(bytes.NewReader(images.CharSetZeroImage))
	} else {
		img, _, err = image.Decode(bytes.NewReader(images.CharSetOneImage))
	}
	if err != nil {
		log.Fatal(err)
	}
	img2 := ebiten.NewImageFromImage(img)
	for i := 0; i <= 255; i++ {
		if charset == 0 {
			n.charImages0[i] = n.charImageSelecta(img2, i)
		} else {
			n.charImages1[i] = n.charImageSelecta(img2, i)
		}
	}
}

// convertPos receives Nimbus-style screen coords and returns then as ebiten-style
func (n *Nimbus) convertPos(x, y, imageHeight int) (ex, ey float64) {
	_, paperHeight := n.paper.Size()
	return float64(x), float64(paperHeight) - float64(y) - float64(imageHeight)
}

// convertColRow receives a Nimbus-style column, row position and returns an
// ebiten-style graphical coordinate
func (n *Nimbus) convertColRow(cr colRow) (ex, ey float64) {
	ex = (float64(cr.col) * 8) - 8
	ey = (float64(cr.row) * 10) - 10
	return ex, ey
}

// validateColour checks that a Nimbus colour index is valid for the current
// screen mode.  If validation fails then a panic is issued.
func (n *Nimbus) validateColour(c int) {
	// Negative values and anything beyond the pallete range is not allowed
	if c < 0 {
		panic("Negative values are not allowed for colours")
	}
	if c > len(n.palette)-1 {
		panic("Colour is out of range for this screen mode")
	}
}

// convertColour receives a Nimbus colour index and returns the RGBA
func (n *Nimbus) convertColour(c int) color.RGBA {
	return n.basicColours[n.palette[c]]
}

// charImageSelecta returns the subimage pointer of a char from the charset
// image for any given ASCII code.  If control char is received, a blank char
// is returned instead.
func (n *Nimbus) charImageSelecta(img *ebiten.Image, c int) *ebiten.Image {

	// select blank char 127 if control char
	if c < 33 {
		c = 127
	}

	// Calculate row and column position of the char on the charset image
	mapNumber := c - 32 // codes < 33 are not on the map
	row := int(math.Ceil(float64(mapNumber) / float64(30)))
	column := mapNumber - (30 * (row - 1))

	// Calculate corners of rectangle around the char
	x1 := (column - 1) * 10
	x2 := x1 + 10
	y1 := (row - 1) * 10
	y2 := y1 + 10

	// Return pointer to sub image
	return img.SubImage(image.Rect(x1, y1, x2, y2)).(*ebiten.Image)
}

// drawChar draws a character at a specific location on an image
func (n *Nimbus) drawChar(image *ebiten.Image, c, x, y, colour, charset int) {
	// Draw char on image and apply colour
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	rgba := n.convertColour(colour)
	r := float64(rgba.R) / 0xff
	g := float64(rgba.G) / 0xff
	b := float64(rgba.B) / 0xff
	op.ColorM.Translate(r, g, b, 0)
	if charset == 0 {
		image.DrawImage(n.charImages0[c], op)
	} else {
		image.DrawImage(n.charImages1[c], op)
	}
}
