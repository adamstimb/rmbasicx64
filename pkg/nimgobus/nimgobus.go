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

// textBox defines the bounding box of a scrollable text box
type textBox struct {
	col1 int
	row1 int
	col2 int
	row2 int
}

type Sprite struct {
	pixels [][]int
	x, y   int
	colour int
	over   bool
}

// Nimbus acts as a container for all the components of the Nimbus monitor.  You
// only need to call the Init() method after declaring a new Nimbus.
type Nimbus struct {
	Monitor               *ebiten.Image
	videoMemory           [250][640]int //[height][width]
	videoImage            *ebiten.Image
	muDrawQueue           sync.Mutex
	drawQueue             []Sprite
	muBorderImage         sync.Mutex
	borderImage           *ebiten.Image
	paper                 *ebiten.Image
	basicColours          []color.RGBA
	borderSize            int
	borderColour          int
	paperColour           int
	penColour             int
	charset               int
	cursorChar            int
	defaultHighResPalette []int
	defaultLowResPalette  []int
	palette               []int
	logoImage             [][]int
	textBoxes             [10]textBox
	imageBlocks           [16]*ebiten.Image
	selectedTextBox       int
	cursorPosition        colRow
	cursorMode            int
	cursorCharset         int
	muCursorFlash         sync.Mutex
	cursorFlash           bool
	cursorFlashEnabled    bool
	cursorIsVisible       bool
	charImages0           [256][][]int
	charImages1           [256][][]int
	//charImages0            [256]*ebiten.Image
	//charImages1            [256]*ebiten.Image
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
	n.drawQueue = []Sprite{}
	n.videoImage = ebiten.NewImage(640, 250)
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
	n.muCursorFlash.Lock()
	n.cursorFlash = false
	n.muCursorFlash.Unlock()
	n.cursorFlashEnabled = true
	n.cursorIsVisible = true
	n.selectedTextBox = 0
	n.keyBuffer = []int{}
	n.keyBufferLock = false
	n.ebitenInputChars = []rune{}
	n.ebitenInputCharsLock = true
	n.muBorderImage.Lock()
	n.borderImage.Fill(n.basicColours[n.palette[n.borderColour]])
	n.muBorderImage.Unlock()

	// Initialize with mode 80 textboxes
	for i := 0; i < 10; i++ {
		n.textBoxes[i] = textBox{1, 1, 80, 25}
	}

	// Set break int
	n.BreakInterruptDetected = false

	// Start tickers
	go n.cursorFlashTicker()
}

// Update redraws the Nimbus monitor image and checks input buffers
func (n *Nimbus) Update() {

	if n.cursorFlashEnabled && n.cursorFlash {
		n.drawCursor()
	}
	// If there's anything in the drawQueue then update video
	if len(n.drawQueue) > 0 {
		n.muDrawQueue.Lock()
		n.updateVideoMemory()
		n.redraw()
		n.muDrawQueue.Unlock()
	}
}

// cursorFlashTicker switches the cursorFlash flag to true every n seconds
func (n *Nimbus) cursorFlashTicker() {
	for {
		if n.cursorFlashEnabled {
			n.muCursorFlash.Lock()
			n.cursorFlash = true
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
