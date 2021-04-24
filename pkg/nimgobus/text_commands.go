package nimgobus

import (
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

// SetPen sets penColour
func (n *Nimbus) SetPen(c int) {
	n.penColour = c
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
	n.cursorPosition = colRow{1, 1} // Relocate cursor
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

// SetCursor changes the cursor state.  Between 1 and 3 parameters can be passed. The
// first parameter sets the cursor mode (< 0 for invisible cursor, 0 for flashing
// cursor, > 0 for visible cursor without flashing), the second parameter sets the ASCII
// code of the cursor char, the third parameter sets the charset of the cursor char.
func (n *Nimbus) SetCursor(p ...int) {
	// Validate
	if len(p) < 1 || len(p) > 3 {
		panic("SetCursor requires 1 to 3 parameters")
	}
	// Set cursor mode
	if p[0] == 0 {
		n.cursorFlashEnabled = false
	}
	if p[0] == 1 {
		n.cursorFlashEnabled = true
	}
	// TODO: support > 1 for none-flashing cursor
	// If 2 parameters, set cursor char as well
	if len(p) > 1 {
		// Validate char
		if p[1] < 0 || p[1] > 255 {
			panic("Invalid cursor char")
		}
		// Set it
		n.cursorChar = p[1]
	}
	if len(p) > 2 {
		// Validate charset
		if p[2] != 0 && p[2] != 1 {
			panic("Invalid charset")
		}
		// Set it
		n.cursorCharset = p[2]
	}
}

// Put draws an ASCII char at the cursor position
func (n *Nimbus) Put(c int) {
	// Validate c
	if c < 0 || c > 255 {
		panic("Character code is out-of-range for extended ASCII (0-255)")
	}
	// Pick the textbox
	box := n.textBoxes[n.selectedTextBox]
	// Get x, y coordinate of cursor and draw the char
	relCurPos := n.cursorPosition
	var absCurPos colRow // we need the absolute cursor position
	absCurPos.col = relCurPos.col + box.col1 - 1
	absCurPos.row = relCurPos.row + box.row1 - 1
	curX, curY := n.convertColRow(absCurPos)
	// Draw the char (unless CR) and advance the cursor
	if c != 13 {
		var charPixels [][]int
		switch n.charset {
		case 0:
			charPixels = n.charImages0[c]
		case 1:
			charPixels = n.charImages1[c]
		}
		// add paper and pen colour to charPixels
		newCharPixels := make2dArray(8, 10)
		for x := 0; x < 8; x++ {
			for y := 0; y < 10; y++ {
				if charPixels[y][x] == 1 {
					newCharPixels[y][x] = n.penColour
				} else {
					newCharPixels[y][x] = n.paperColour
				}
			}
		}
		n.drawSprite(Sprite{newCharPixels, curX, curY, -1, true})
		n.AdvanceCursor(false)
	} else {
		// Force carriage return if 13 was passed
		n.AdvanceCursor(true)
	}
}

// Print
func (n *Nimbus) Print(s string) {
	for _, c := range s {
		n.Put(int(c))
	}
}

// Get returns a single character code input from the keyboard.
// If no key was pressed then -1 is returned.
func (n *Nimbus) Get() int {
	return n.popKeyBuffer()
}

// Input receives keyboard input into a string of up to 256 chars and returns
// the string when ENTER is pressed.
// The user can edit the string using the delete key and left and right arrow
// keys.  A prompt is printed on the screen at the current cursor position and
// the user's input is echoed to screen after the prompt.  The input buffer can
// also be pre-populated.
func (n *Nimbus) Input(prepopulateBuffer string) string {

	// Flush keyBuffer, initialize internal buffer and prepopulate it
	n.flushKeyBuffer()
	var buffer []int
	for _, c := range prepopulateBuffer {
		buffer = append(buffer, int(c))
	}
	bufferPosition := len(buffer)
	maxBufferSize := 255
	//startPos := n.cursorPosition
	//endPos := n.cursorPosition

	// popBuffer pops a char from the buffer at a given position
	popBuffer := func(buffer []int, indexToPop int) []int {
		var newBuffer []int
		for index, oldValue := range buffer {
			if index != indexToPop {
				newBuffer = append(newBuffer, oldValue)
			}
		}
		return newBuffer
	}

	// pushBuffer pushes a char on to the buffer at a given position
	pushBuffer := func(buffer []int, indexToPush int, newValue int) []int {
		var newBuffer []int
		if len(buffer) == indexToPush {
			// appending char to end of buffer
			newBuffer = append(buffer, newValue)
		} else {
			// pushing new char before the end of the buffer
			for index, oldValue := range buffer {
				if index == indexToPush {
					newBuffer = append(newBuffer, newValue)
					newBuffer = append(newBuffer, oldValue)
				} else {
					newBuffer = append(newBuffer, oldValue)
				}
			}
		}
		return newBuffer
	}

	// echoBuffer prints the chars in the buffer from a given buffer position
	echoBuffer := func(buffer []int, startIndex int) {
		for i := startIndex; i < len(buffer); i++ {
			n.Put(buffer[i])
		}
	}

	// moveCursorBack moves the cursor backwards one char along a line
	// of input that may span more than one line.  If andDelete is true
	// it will also delete the previous char.
	moveCursorBack := func(andDelete bool) {
		// handle deleting from the same line
		if n.cursorPosition.col > 1 {
			n.cursorPosition.col--
			if andDelete {
				n.Put(32)
				n.cursorPosition.col--
			}
			return
		}
		// handle deleting from the line above
		if n.cursorPosition.col == 1 && n.cursorPosition.row > 1 {
			// get width of current textbox
			box := n.textBoxes[n.selectedTextBox]
			width := box.col2 - box.col1
			// go up a line and delete end char
			n.cursorPosition.row--
			n.cursorPosition.col = width // + 1
			if andDelete {
				n.Put(32)
				n.cursorPosition.row--
				n.cursorPosition.col = width // + 1
			}
			return
		}
	}

	// moveCursorForward moves the cursor forward and if necessary onto the line below
	moveCursorForward := func() {
		// get width of current textbox
		box := n.textBoxes[n.selectedTextBox]
		width := box.col2 - box.col1
		// move cursor forward
		if n.cursorPosition.col < width {
			// just shift cursor right if we're not at the end of a line
			n.cursorPosition.col++
		} else {
			// otherwise go to line below
			n.cursorPosition.col = box.col1
			n.cursorPosition.row++
		}
	}

	// Print the buffer before looping to get user input
	echoBuffer(buffer, 0)

	// now loop to received and edit the input string until enter is pressed
	for !n.BreakInterruptDetected {
		// get most recent keyboard input
		char := n.Get()
		if char == -1 {
			// nothing pressed so update vars an skip
			continue
		}
		// handle control keys if any
		if char <= -10 {
			// is control key
			if char == -11 {
				// ENTER pressed so echo buffer beyond current position
				// one last time and break loop
				echoBuffer(buffer, bufferPosition)
				break
			}
			if char == -10 {
				// BACKSPACE pressed
				if bufferPosition > 0 {
					// only delete if not at beginning
					bufferPosition--
					buffer = popBuffer(buffer, bufferPosition)
					// delete char on screen
					moveCursorBack(true)
					// if deleting from before the end of the buffer, rewrite
					// the rest of the buffer and delete the trailing char
					if bufferPosition < len(buffer) {
						tempCursorPosition := n.cursorPosition
						echoBuffer(buffer, bufferPosition)
						n.Put(32)
						n.cursorPosition = tempCursorPosition
					}
				}
			}
			if char == -12 {
				// LEFT ARROW pressed
				if bufferPosition > 0 {
					// only move left if not at beginning
					bufferPosition--
					moveCursorBack(false)
				}
			}
			if char == -13 {
				// RIGHT ARROW pressed
				if bufferPosition < len(buffer) {
					// only move right if not at end of buffer
					bufferPosition++
					moveCursorForward()
				}
			}
			if char == -14 {
				// UP ARROW pressed
				lastCol := n.cursorPosition.col
				for bufferPosition > 0 {
					moveCursorBack(false)
					bufferPosition--
					if n.cursorPosition.col == lastCol {
						break
					}
				}
			}
			if char == -15 {
				// DOWN ARROW pressed
				lastCol := n.cursorPosition.col
				for bufferPosition < len(buffer) {
					moveCursorForward()
					bufferPosition++
					if n.cursorPosition.col == lastCol {
						break
					}
				}
			}
		} else {
			// is printable char
			// only accept it if we have space
			if bufferPosition <= maxBufferSize {
				// push new char into buffer
				buffer = pushBuffer(buffer, bufferPosition, char)
				// if new char added before end of buffer, rewrite rest of buffer
				// otherwise simply put the last char in the buffer
				if bufferPosition < len(buffer) {
					tempCursorPosition := n.cursorPosition
					echoBuffer(buffer, bufferPosition)
					n.cursorPosition = tempCursorPosition
					moveCursorForward()
				} else {
					n.Put(buffer[len(buffer)-1])
				}
				bufferPosition++
			}
		}
	}

	// Enter was pressed so carriage return
	n.Put(13)
	_ = n.popKeyBuffer()
	// Put inputBuffer into a string
	var returnString string
	for _, c := range buffer {
		returnString += string(rune(c))
	}
	// And return it
	return returnString
}
