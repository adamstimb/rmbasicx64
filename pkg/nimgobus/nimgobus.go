package nimgobus

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/adamstimb/rmbasicx64/pkg/nimgobus/images"
	"github.com/hajimehoshi/ebiten/v2"
)

// colRow defines a column, row position
type colRow struct {
	col int
	row int
}

func (c *colRow) AskCurpos() (int, int) {
	return c.col, c.row
}

// textBox defines the bounding box of a scrollable text box
type textBox struct {
	col1 int
	row1 int
	col2 int
	row2 int
}

// drawingBox defines the bounding box of a drawing box
type drawingBox struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

// Sprite defines a sprite that contains a 2d image array, a screen co-ordinate, colour and overwrite (XOR) mode
type Sprite struct {
	pixels [][]int
	x, y   int
	colour int
	over   bool
}

// repeatingChar is used to store and count repeating chars for dynamically limiting repeating key presses
type repeatingChar struct {
	char    int
	counter int
}

// Nimbus acts as a container for all the components of the Nimbus monitor.  You
// only need to call the Init() method after declaring a new Nimbus.
type Nimbus struct {
	Monitor                *ebiten.Image     // The Monitor image including background
	muVideoMemory          sync.Mutex        //
	videoMemory            [250][640]int     // Represents the video memory, basically a 640x250 array of integers represents the Nimbus colour index of each pixel
	muVideoMemoryOverlay   sync.Mutex        //
	videoMemoryOverlay     [250][640]int     // A copy of the video memory where temporal things like cursors can be drawn
	videoImage             *ebiten.Image     // An ebiten image derived from videoMemoryOverlay
	muDrawQueue            sync.Mutex        //
	drawQueue              []Sprite          // A queue of sprites to be written to videoMemory
	redrawComplete         bool              // Flag to indicate if nimgobus is working on the drawQueue
	muBorderImage          sync.Mutex        //
	borderImage            *ebiten.Image     // An ebiten image representing the Nimbus monitor's background
	basicColours           []color.RGBA      // An array of the Nimbus's 16 built-in colours
	textBoxes              [10]textBox       // All defined textboxes
	drawingBoxes           [10]drawingBox    // All define drawing boxes
	imageBlocks            [16]*ebiten.Image // Images loaded into memory as "blocks"
	logoImage              [][]int           // The "RM Nimbus" branding image
	charImages0            [256][][]int      // An array of 2d arrays for each char in this charset
	charImages1            [256][][]int      // as above
	borderSize             int               // The width of the border in pixels
	borderColour           int               // The current border colour
	paperColour            int               // The current paper colour
	penColour              int               // The current pen colour
	charset                int               // The current char set (0 or 1)
	cursorCharset          int               // The current cursor char set (0 or 1)
	cursorChar             int               // The current cursor char
	brush                  int               // The current drawing/plot brush
	plotDirection          int               // The current plot direction
	plotFont               int               // The current plot font
	plotSizeX              int               // The current plot size x
	plotSizeY              int               // The current plot size y
	over                   bool              // The drawing mode (XOR)
	selectedTextBox        int               // The current textbox
	selectedDrawingBox     int               // The current drawing box
	defaultHighResPalette  []int             // The default palette for high-res (mode 80)
	defaultLowResPalette   []int             // The default palette for low-res (mode 40)
	palette                []int             // The current palette
	cursorPosition         colRow            // The current cursor position
	cursorMode             int               // The current cursor mode
	muCursorFlash          sync.Mutex        //
	cursorFlash            bool              // Cursor flash flag: Everytime this changes the cursor with flash if enabled
	cursorFlashEnabled     bool              // Flag to indicate if cursor flash is enabled
	deleteMode             bool              // true if delete mode selected
	deleteModeCursorImage  [][]int           // The special cursor for delete mode
	muKeyBuffer            sync.Mutex        //
	keyBuffer              []int             // Nimgobus needs it's own key buffer since ebiten's only deals with printable chars
	charRepeat             repeatingChar     // Used by the keyBuffer to dynamically limit key presses
	BreakInterruptDetected bool              // Flag is set to true if user makes a <BREAK>
}

// Init initializes a new Nimbus.  You must call this method after declaring a
// new Nimbus variable.
func (n *Nimbus) Init() {
	// in case any randomonia is required we can run a seed on startup
	rand.Seed(time.Now().UnixNano())

	// should next exceed 50tps
	ebiten.SetMaxTPS(50)

	// Load Nimbus logo image and both charsets
	n.loadLogoImage()
	n.loadCharsetImages(0)
	n.loadCharsetImages(1)
	// Set init values of everything else
	n.borderSize = 50
	n.borderImage = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
	n.Monitor = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
	n.drawQueue = []Sprite{}
	n.redrawComplete = true
	n.videoImage = ebiten.NewImage(640, 250)
	n.basicColours = append(n.basicColours, basicColours...)
	n.defaultHighResPalette = append(n.defaultHighResPalette, defaultHighResPalette...)
	n.defaultLowResPalette = append(n.defaultLowResPalette, defaultLowResPalette...)
	n.palette = append(n.palette, n.defaultHighResPalette...)
	n.borderColour = 0
	n.paperColour = 0
	n.penColour = 3
	n.charset = 0
	n.cursorMode = -1
	n.cursorChar = 95
	n.cursorCharset = 0
	n.brush = 3
	n.plotDirection = 0
	n.plotFont = 0
	n.plotSizeX = 1
	n.plotSizeY = 1
	n.over = true
	n.cursorPosition = colRow{1, 1}
	n.muCursorFlash.Lock()
	n.cursorFlash = false
	n.muCursorFlash.Unlock()
	n.cursorFlashEnabled = true
	n.deleteMode = false
	n.selectedTextBox = 0
	n.selectedDrawingBox = 0
	n.keyBuffer = []int{}
	n.charRepeat = repeatingChar{0, 0}
	n.muBorderImage.Lock()
	n.borderImage.Fill(n.basicColours[n.palette[n.borderColour]])
	n.muBorderImage.Unlock()

	// Draw delete mode cursor
	n.deleteModeCursorImage = make2dArray(8, 10)
	for x := 0; x < 8; x++ {
		for y := 0; y < 10; y++ {
			n.deleteModeCursorImage[y][x] = 1
		}
	}

	// Initialize with mode 80 textboxes and drawingboxes
	for i := 0; i < 10; i++ {
		n.textBoxes[i] = textBox{1, 1, 80, 25}
		n.drawingBoxes[i] = drawingBox{0, 0, 639, 249}
	}

	// Set break int
	n.BreakInterruptDetected = false

	// Start tickers
	go n.cursorFlashTicker()
}

// Update is called by every Draw() call from ebiten
func (n *Nimbus) Update() {
	n.updateKeyBuffer()
	n.updateVideoMemory()
	n.redraw()
}

// flushKeyBuffer empties the keyBuffer
func (n *Nimbus) flushKeyBuffer() {
	n.muKeyBuffer.Lock()
	n.keyBuffer = []int{}
	n.muKeyBuffer.Unlock()
}

// updateKeyBuffer updates nimgobus's key buffer
func (n *Nimbus) updateKeyBuffer() {

	// acceptRepeatingChar prevents unprintable chars that aren't controlled by ebiten's
	// keyboard buffer bombing nimgobus' buffer.
	acceptRepeatingChar := func(char int) {
		if char != n.charRepeat.char {
			// not the same as last char so add it and reset counter
			n.charRepeat.char = char
			n.charRepeat.counter = 0
			n.keyBuffer = append(n.keyBuffer, char)
		} else {
			// is the same char so only add if repeated 20 times or more than 20 times and a multiple of 5
			if n.charRepeat.counter == 40 || (n.charRepeat.counter > 40 && n.charRepeat.counter%5 == 0) {
				n.keyBuffer = append(n.keyBuffer, char)
				n.charRepeat.counter++
			} else {
				n.charRepeat.counter++
			}
		}
	}

	n.muKeyBuffer.Lock()
	inputChars := ebiten.InputChars()
	for _, r := range inputChars {
		// Copy printable keys from ebiten's InputChars
		n.keyBuffer = append(n.keyBuffer, int(r))
	}
	// Evaluate control keys
	// handle BREAK interrupt first (Ctrl+ScrollLock)
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyScrollLock) {
		n.BreakInterruptDetected = true
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyKPEnter) {
		acceptRepeatingChar(-11)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		acceptRepeatingChar(-10)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		acceptRepeatingChar(-12)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		acceptRepeatingChar(-13)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		acceptRepeatingChar(-14)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		acceptRepeatingChar(-15)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyHome) {
		acceptRepeatingChar(-16)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnd) {
		acceptRepeatingChar(-17)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		acceptRepeatingChar(-18)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		acceptRepeatingChar(-19)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyDelete) {
		acceptRepeatingChar(-20)
		n.muKeyBuffer.Unlock()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyInsert) {
		acceptRepeatingChar(-21)
		n.muKeyBuffer.Unlock()
		return
	}
	n.charRepeat.char = 0
	n.charRepeat.counter = 0
	n.muKeyBuffer.Unlock()
}

// cursorFlashTicker flips the cursorFlash flag every 500 ms
func (n *Nimbus) cursorFlashTicker() {
	for {
		if n.cursorFlashEnabled {
			n.muCursorFlash.Lock()
			n.cursorFlash = !n.cursorFlash
			n.muCursorFlash.Unlock()
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// loadLogoImage loads the Nimbus logo image
func (n *Nimbus) loadLogoImage() {

	// convertToArray receives the logo image and returns it as 3-colour 2d array
	convertToArray := func(img image.Image) [][]int {
		b := img.Bounds()
		width := b.Max.X
		height := b.Max.Y
		newArray := make2dArray(width, height)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				// Get colour at x, y and downsample to 3 colours
				c := img.At(x, y)
				r, g, b, a := c.RGBA()
				if r > 60000 && g > 60000 && b > 60000 && a > 60000 {
					// white
					newArray[y][x] = 3
					continue
				}
				if r > 60000 && (g > 20000 && g < 30000) && (b > 2000 && b < 30000) && a > 60000 {
					// red
					newArray[y][x] = 2
					continue
				}
				if r < 10000 && g < 10000 && b < 10000 && a > 60000 {
					// black
					newArray[y][x] = 0
					continue
				}
			}
		}
		return newArray
	}

	img, _, err := image.Decode(bytes.NewReader(images.NimbusLogoImage))
	if err != nil {
		log.Fatal(err)
	}
	n.logoImage = convertToArray(img)
}

// loadCharsetImages loads the charset images
func (n *Nimbus) loadCharsetImages(charset int) {

	// convertToArray receives an char image and returns it as black-and-white 2d array
	convertToArray := func(img image.Image) [][]int {
		b := img.Bounds()
		width := b.Max.X
		height := b.Max.Y
		newArray := make2dArray(width, height)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				// Get colour at x, y and if black set it to 1 in the 2d array
				c := img.At(x, y)
				_, _, _, a := c.RGBA()
				if a == 65535 {
					newArray[y][x] = 1
				}
			}
		}
		return newArray
	}

	var img image.Image
	var imgArray [][]int
	var err error
	if charset == 0 {
		img, _, err = image.Decode(bytes.NewReader(images.CharSetZeroImage))
	} else {
		img, _, err = image.Decode(bytes.NewReader(images.CharSetOneImage))
	}
	if err != nil {
		log.Fatal(err)
	}
	imgArray = convertToArray(img)
	for i := 0; i <= 255; i++ {
		if charset == 0 {
			n.charImages0[i] = n.charImageSelecta(imgArray, i)
		} else {
			n.charImages1[i] = n.charImageSelecta(imgArray, i)
		}
	}
}

// charImageSelecta returns the subimage pointer of a char from the charset
// image for any given ASCII code.  If control char is received, a blank char
// is returned instead.
func (n *Nimbus) charImageSelecta(img [][]int, c int) [][]int {
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
	// Get char and return
	charImgArray := make2dArray(10, 10)
	charImgArrayX := 0
	for x := x1; x < x2; x++ {
		charImgArrayY := 0
		for y := y1; y < y2; y++ {
			charImgArray[charImgArrayY][charImgArrayX] = img[y][x]
			charImgArrayY++
		}
		charImgArrayX++
	}
	return charImgArray
}

// popKeyBuffer pops the oldest char in the buffer
// If the buffer is empty -1 is returned.
func (n *Nimbus) popKeyBuffer() int {
	// check if buffer is empty and return -1 if so
	if len(n.keyBuffer) == 0 {
		// is empty
		return -1
	}
	// Otherwise pop the buffer
	n.muKeyBuffer.Lock()
	char := n.keyBuffer[0]
	// if buffer only has one char re-initialize it otherwise shorten it by 1 element
	if len(n.keyBuffer) <= 1 {
		n.keyBuffer = []int{}
	} else {
		n.keyBuffer = n.keyBuffer[1:]
	}
	// all done
	n.muKeyBuffer.Unlock()
	return char
}
