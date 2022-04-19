package nimgobus

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/adamstimb/rmbasicx64/pkg/nimgobus/resources/font"
	"github.com/adamstimb/rmbasicx64/pkg/nimgobus/resources/logo"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
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
	pixels    [][]int
	x, y      int
	colour    int
	over      bool
	fillStyle FillStyle
}

// repeatingChar is used to store and count repeating chars for dynamically limiting repeating key presses
type repeatingChar struct {
	char    int
	counter int
}

// colourFlashSetting is used to store colour flash settings for individual pallete slots
type colourFlashSetting struct {
	speed       int // 0 no flash, 1 slow, 2 fast
	flashColour int // the palette slot to of the flash colour
}

// imageBlock is how RM Basic used to store image data in the ANIMATE extension
type imageBlock struct {
	image   [][]int // 2d array containing the image
	mode    int     // The screen mode from when the image was captured/fetched in chars (40 or 80)
	deleted bool    // Set to true if Delblock is called on the block
}

// FillStyle describes the fill settings for AREA, FLOOD, CIRCLE, and SLICE
type FillStyle struct {
	Style    int // 1 for solid/dithered, 2 for hatched, 3 for hollow (edge)
	Hatching int // Hatching type if Style==2
	Colour2  int // 2nd hatching colour if Style==2
}

// FileObj describes a file object and whether its for writing or reading
type FileObj struct {
	File    *os.File
	Writing bool
}

// Nimbus acts as a container for all the components of the Nimbus monitor.  You
// only need to call the Init() method after declaring a new Nimbus.
type Nimbus struct {
	PaddingX               int
	PaddingY               int
	Scale                  float64
	Monitor                *ebiten.Image        // The Monitor image including background
	muVideoMemory          sync.Mutex           //
	videoMemory            [250][640]int        // Represents the video memory, basically a 640x250 array of integers represents the Nimbus colour index of each pixel
	muVideoMemoryOverlay   sync.Mutex           //
	videoMemoryOverlay     [250][640]int        // A copy of the video memory where temporal things like cursors can be drawn
	videoImage             *ebiten.Image        // An ebiten image derived from videoMemoryOverlay
	muDrawQueue            sync.Mutex           //
	drawQueue              []Sprite             // A queue of sprites to be written to videoMemory
	redrawComplete         bool                 // Flag to indicate if nimgobus is working on the drawQueue
	muBorderImage          sync.Mutex           //
	borderImage            *ebiten.Image        // An ebiten image representing the Nimbus monitor's background
	mode                   int                  // The current screen mode as set by SET MODE
	basicColours           []color.RGBA         // An array of the Nimbus's 16 built-in colours
	textBoxes              [10]textBox          // All defined textboxes
	drawingBoxes           [10]drawingBox       // All define drawing boxes
	imageBlocks            [100]imageBlock      // Images loaded into memory as "blocks"
	logoImage              [][]int              // The "RM Nimbus" branding image
	charImages0            [256][][]int         // An array of 2d arrays for each char in this charset
	charImages1            [256][][]int         // as above
	pointsStyles           [][][]int            // An array of 2d arrays representing the built-in points styles
	pointsStyle            int                  // The current points style
	lineStyles             [][]int              // An array of 2d arrays representing the built-in line styles
	lineStyle              int                  // The current line style
	patterns               [][4][4]int          // The brush patterns
	hatchings              [][16][16]int        // The fill hatchings
	fillStyle              FillStyle            // The current fill style
	useFillStyle           bool                 // Set to true to render with fill style
	borderSize             int                  // The width of the border in pixels
	borderColour           int                  // The current border colour
	paperColour            int                  // The current paper colour
	penColour              int                  // The current pen colour
	charset                int                  // The current char set (0 or 1)
	cursorCharset          int                  // The current cursor char set (0 or 1)
	cursorChar             int                  // The current cursor char
	brush                  int                  // The current drawing/plot brush
	plotDirection          int                  // The current plot direction
	plotFont               int                  // The current plot font
	plotSizeX              int                  // The current plot size x
	plotSizeY              int                  // The current plot size y
	over                   bool                 // The drawing mode (XOR)
	selectedTextBox        int                  // The current textbox
	selectedDrawingBox     int                  // The current drawing box
	defaultHighResPalette  []int                // The default palette for high-res (mode 80)
	defaultLowResPalette   []int                // The default palette for low-res (mode 40)
	palette                []int                // The current palette
	muColourFlashSettings  sync.Mutex           //
	colourFlashSettings    []colourFlashSetting // Colour flash settings for each palette slot
	cursorPosition         colRow               // The current cursor position
	cursorMode             int                  // The current cursor mode
	muCursorFlash          sync.Mutex           //
	cursorFlash            bool                 // Cursor flash flag: Everytime this changes the cursor with flash if enabled
	cursorFlashEnabled     bool                 // Flag to indicate if cursor flash is enabled
	colourFlash            int                  // The colour flash counter
	deleteMode             bool                 // true if delete mode selected
	deleteModeCursorImage  [][]int              // The special cursor for delete mode
	muKeyBuffer            sync.Mutex           //
	keyBuffer              []int                // Nimgobus needs it's own key buffer since ebiten's only deals with printable chars
	charRepeat             repeatingChar        // Used by the keyBuffer to dynamically limit key presses
	MouseX                 int                  // Mouse position and button press
	MouseY                 int                  //
	MouseButton            int                  //
	MouseOn                bool                 //
	sound                  bool                 // 3-channel Nimbus synth chip goodness
	voices                 []voice              //
	envelopes              []envelope           //
	selectedVoice          int                  //
	BreakInterruptDetected bool                 // Flag is set to true if user makes a <BREAK>
	FileChannels           map[int]*FileObj     // File channels and their objects are stored here when they're opened/created
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
	n.Clearblock()
	n.pointsStyles = append(n.pointsStyles, defaultPointsStyles...)
	n.pointsStyle = 1
	n.lineStyles = append(n.lineStyles, defaultLineStyles...)
	n.lineStyle = 0
	n.borderSize = 50
	n.patterns = append(n.patterns, defaultHighResPatterns...)
	n.hatchings = append(n.hatchings, defaultHatchings...)
	n.borderImage = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
	n.Monitor = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
	n.drawQueue = []Sprite{}
	n.redrawComplete = true
	n.videoImage = ebiten.NewImage(640, 250)
	n.mode = 80
	n.basicColours = append(n.basicColours, basicColours...)
	n.defaultHighResPalette = append(n.defaultHighResPalette, defaultHighResPalette...)
	n.defaultLowResPalette = append(n.defaultLowResPalette, defaultLowResPalette...)
	n.palette = append(n.palette, n.defaultHighResPalette...)
	n.useFillStyle = false
	n.initColourFlashSettings()
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
	n.MouseX = 0
	n.MouseY = 0
	n.MouseButton = 0
	n.MouseOn = false
	n.muBorderImage.Lock()
	n.borderImage.Fill(n.basicColours[n.palette[n.borderColour]])
	n.muBorderImage.Unlock()
	n.sound = false
	n.FileChannels = make(map[int]*FileObj)

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

	// Start tickers, etc.
	go n.cursorFlashTicker()
	n.colourFlash = 0
	go n.colourFlashTicker()
}

// Update is called by every game.Update()
func (n *Nimbus) Update(PaddingX, PaddingY int, Scale float64) {
	// sound
	if n.sound {
		// NewContext only needs to come from one voice
		if n.voices[0].audioContext == nil {
			n.voices[0].audioContext = audio.NewContext(sampleRate)
		}
		for v := 0; v < len(n.voices); v++ {
			if n.voices[v].player == nil {
				var err error
				n.voices[v].player, err = audio.NewPlayer(n.voices[0].audioContext, &stream{voice: &n.voices[v]})
				if err != nil {
					log.Printf("Could not start voice %d audio: %e", v, err)
				} else {
					n.voices[v].player.Play()
				}
			}
		}
	}
	// video
	n.PaddingX = PaddingX
	n.PaddingY = PaddingY
	n.Scale = Scale
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
	// handle BREAK interrupt first (Ctrl+ScrollLock or Ctrl+B for computers without ScrollLock)
	if ebiten.IsKeyPressed(ebiten.KeyControl) && (ebiten.IsKeyPressed(ebiten.KeyScrollLock) || ebiten.IsKeyPressed(ebiten.KeyB)) {
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

// colourFlashTicker increments the colourFlash counter every 500 ms
func (n *Nimbus) colourFlashTicker() {
	for {
		n.colourFlash++
		if n.colourFlash > 3 {
			n.colourFlash = 0
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
					// green
					newArray[y][x] = 1
					continue
				}
			}
		}
		return newArray
	}

	img, _, err := image.Decode(bytes.NewReader(logo.NimbusLogoFinal_png))
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
		img, _, err = image.Decode(bytes.NewReader(font.Charsets_png)) //images.CharSetZeroImage
	} else {
		img, _, err = image.Decode(bytes.NewReader(font.Charsets_png)) //images.CharSetZeroImage
	}
	if err != nil {
		log.Fatal(err)
	}
	imgArray = convertToArray(img)
	for i := 0; i <= 255; i++ {
		if charset == 0 {
			n.charImages0[i] = n.charImageSelecta(imgArray, i, charset)
		} else {
			n.charImages1[i] = n.charImageSelecta(imgArray, i, charset)
		}
	}
}

// charImageSelecta returns the subimage pointer of a char from the charset
// image for any given ASCII code.  If control char is received, a blank char
// is returned instead.
func (n *Nimbus) charImageSelecta(img [][]int, c int, charSet int) [][]int {
	// Hotfix: RM Basic on the emulator somehow skips chars 28 and 29 -
	// we'll redo the dump again later.  For now, we'll just return a blank for
	// those chars and increment c by 2 if c > 29
	if c >= 28 && c <= 29 {
		c = 127
	}
	if c > 29 {
		c -= 2
	}
	charsPerRow := 40
	c++
	row := int(math.Ceil(float64(c) / float64(charsPerRow)))
	column := c - (charsPerRow * (row - 1))
	// Calculate corners of rectangle around the char
	var yOffset int
	if charSet == 0 {
		yOffset = 0
	} else {
		yOffset = 70
	}
	x1 := (column - 1) * 8
	x2 := x1 + 8
	y1 := yOffset + ((row - 1) * 10)
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
