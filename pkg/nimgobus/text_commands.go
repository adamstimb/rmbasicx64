package nimgobus

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// AskMode returns the current screen mode (40 column or 80 column)
func (n *Nimbus) AskMode() int {
	width, _ := n.videoImage.Size()
	if width == 320 {
		return 40 // low-res mode 40
	}
	if width == 640 {
		return 80 // high-res mode 80
	}
	return 0 // this never happens
}

// make2dArray initializes an empty 2d array and returns it
func make2dArray(width, height int) [][]int {
	a := make([][]int, height)
	for i := range a {
		a[i] = make([]int, width)
	}
	return a
}

// Cls clears the selected textbox if no parameters are passed, or clears another
// textbox if one parameter is passed
func (n *Nimbus) Cls(p ...int) {
	// Validate number of parameters
	if len(p) != 0 && len(p) != 1 {
		// invalid
		panic("Cls accepts either 0 or 1 parameters")
	}
	// Create one big sprite with every pixel set to paperColour and draw it
	width, height := n.videoImage.Size()
	blankPaper := make2dArray(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			blankPaper[y][x] = 1
		}
	}
	n.drawSprite(Sprite{blankPaper, 0, 0, n.paperColour, true})
}

// SetPaper sets paperColour
func (n *Nimbus) SetPaper(c int) {
	n.paperColour = c
}

// SetBorder sets the borderColour
func (n *Nimbus) SetBorder(c int) {
	// Repaint border image with new colour and force redraw
	n.borderColour = c
	n.muBorderImage.Lock()
	n.borderImage.Fill(n.basicColours[n.palette[n.borderColour]])
	n.muBorderImage.Unlock()
	n.ForceRedraw()
}

// SetMode sets the screen mode. 40 is low-resolution, high-colour mode (320x250) and
// 80 is high-resolutions, low-colour mode (640x250)
func (n *Nimbus) SetMode(columns int) {
	if columns != 40 && columns != 80 {
		// RM Basic would just pass if an invalid column number was given so
		// we'll do the same
		return
	}
	// Need to manipulate videoImage so force redraw and get the lock
	n.ForceRedraw()
	n.muDrawQueue.Lock()
	n.muBorderImage.Lock()
	if columns == 40 {
		// low-resolution, high-colour mode (320x250)
		n.videoImage = ebiten.NewImage(320, 250)
		n.paperColour = 0
		n.borderColour = 0
		n.penColour = 15
		n.palette = n.defaultLowResPalette
		// reinit border image and fill with new colour
		n.borderImage = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
		n.borderImage.Fill(n.basicColours[n.palette[n.borderColour]])
	}
	if columns == 80 {
		// high-resolutions, low-colour mode (640x250)
		n.videoImage = ebiten.NewImage(640, 250)
		n.palette = n.defaultHighResPalette
		n.paperColour = 0
		n.borderColour = 0
		n.penColour = 3
		// reinit border image and fill with new colour
		n.borderImage = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
		n.borderImage.Fill(n.basicColours[n.palette[n.borderColour]])
	}
	n.cursorPosition = colRow{1, 1}                        // Relocate cursor
	n.paper.Fill(n.basicColours[n.palette[n.paperColour]]) // Apply paper colour
	// Redefine textboxes, imageBlocks and clear screen
	for i := 0; i < 10; i++ {
		n.textBoxes[i] = textBox{1, 1, columns, 25}
	}
	n.imageBlocks = [16]*ebiten.Image{}
	n.muBorderImage.Unlock()
	n.muDrawQueue.Unlock()
	n.Cls()
}

// convertColRow receives a Nimbus-style column, row position and returns a
// Nimbus-style x, y coordinate
func (n *Nimbus) convertColRow(cr colRow) (x, y int) {
	x = (cr.col - 1) * 8
	y = len(n.videoMemory) - (cr.row * 10)
	return x, y
}

// convertColRowEbiten receives a Nimbus-style column, row position and returns
// an ebiten-style x, y coordinate
func (n *Nimbus) convertColRowEbiten(cr colRow) (x, y int) {
	x = (cr.col * 8) - 8
	y = (cr.row * 10) - 10
	return x, y
}

// SetCurpos sets the cursor position within the selected text box
func (n *Nimbus) SetCurpos(col, row int) {
	// Pick the textbox
	box := n.textBoxes[n.selectedTextBox]
	// Validate col and row position
	if col < 0 || row < 0 {
		panic("Negative column or row values are not allowed")
	}
	width := box.col2 - box.col1
	height := box.row2 - box.row1
	if col > width {
		panic("Column value is outside selected textbox")
	}
	if row > height {
		panic("Row value is outside selected textbox")
	}
	// Validation passed, set cursor position
	n.cursorPosition = colRow{col, row}
}

// Put draws an ASCII char at the cursor position
func (n *Nimbus) Put(c int) {
	// Validate c
	if c < 0 || c > 255 {
		panic("Character code is out-of-range for extended ASCII (0-255)")
	}
	// Pick the textbox
	box := n.textBoxes[n.selectedTextBox]
	width := box.col2 - box.col1
	height := box.row2 - box.row1

	// Get x, y coordinate of cursor and draw the char
	relCurPos := n.cursorPosition
	var absCurPos colRow // we need the absolute cursor position
	absCurPos.col = relCurPos.col + box.col1 - 1
	absCurPos.row = relCurPos.row + box.row1 - 1
	curX, curY := n.convertColRow(absCurPos)

	// Draw paper under the char
	paper := make2dArray(8, 10)
	for x := 0; x < 8; x++ {
		for y := 0; y < 10; y++ {
			paper[y][x] = 1
		}
	}
	n.drawSprite(Sprite{paper, curX, curY, n.paperColour, true})

	// Draw the char (unless CR)
	if c != 13 {
		var charPixels [][]int
		switch n.charset {
		case 0:
			charPixels = n.charImages0[c]
		case 1:
			charPixels = n.charImages1[c]
		}
		n.drawSprite(Sprite{charPixels, curX, curY, n.penColour, true})
		// Update relative cursor position
		relCurPos.col++
	}

	// Carriage return?
	if relCurPos.col > width+1 || c == 13 {
		// over the edge so carriage return
		relCurPos.col = 1
		relCurPos.row++
	}
	// Line feed?
	if relCurPos.row > height+1 {
		// move cursor up and scroll textbox
		relCurPos.row--
		// Scroll up:
		// Define bounding rectangle for the textbox
		x1, y1 := n.convertColRowEbiten(colRow{box.col1, box.row1})
		x2, y2 := n.convertColRowEbiten(colRow{box.col2, box.row2})
		x2 += 8
		y2 += 10
		// We have to manipulate videoImage itself next, so force redraw and get the drawQueue lock
		n.ForceRedraw()
		n.muDrawQueue.Lock()
		// Copy actual textbox image
		oldTextBoxImg := n.videoImage.SubImage(image.Rect(int(x1), int(y1), int(x2), int(y2))).(*ebiten.Image)
		// Create a new textbox image and fill it with paper colour
		newTextBoxImg := ebiten.NewImage(int(x2-x1), int(y2-y1))
		newTextBoxImg.Fill(n.basicColours[n.palette[n.paperColour]])
		// Place old textbox image on new image 10 pixels higher
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, -10)
		newTextBoxImg.DrawImage(oldTextBoxImg, op)
		// Redraw the textbox on the paper
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x1), float64(y1))
		n.paper.DrawImage(newTextBoxImg, op)
		n.muDrawQueue.Unlock()
	}
	// Set new cursor position
	n.cursorPosition = relCurPos
}

// Print
func (n *Nimbus) Print(s string) {
	for _, c := range s {
		n.Put(int(c))
	}
}

// Input
func (n *Nimbus) Input(a, b string) string {
	return ""
}
