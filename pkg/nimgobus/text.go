package nimgobus

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	absCurPos.col = relCurPos.col + box.col1
	absCurPos.row = relCurPos.row + box.row1
	ex, ey := n.convertColRow(absCurPos)
	ex -= 8
	ey -= 10

	// Draw paper under the char
	img := ebiten.NewImage(8, 10)
	img.Fill(n.convertColour(n.paperColour))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(img, op)
	// Draw the char (unless CR)
	if c != 13 {
		n.drawChar(n.paper, c, int(ex), int(ey), n.penColour, n.charset)
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
		x1, y1 := n.convertColRow(colRow{box.col1, box.row1})
		x2, y2 := n.convertColRow(colRow{box.col2, box.row2})
		x2 += 8
		y2 += 10
		// Copy actual textbox image
		oldTextBoxImg := n.paper.SubImage(image.Rect(int(x1), int(y1), int(x2), int(y2))).(*ebiten.Image)
		// Create a new textbox image and fill it with paper colour
		newTextBoxImg := ebiten.NewImage(int(x2-x1), int(y2-y1))
		newTextBoxImg.Fill(n.convertColour(n.paperColour))
		// Place old textbox image on new image 10 pixels higher
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, -10)
		newTextBoxImg.DrawImage(oldTextBoxImg, op)
		// Redraw the textbox on the paper
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x1, y1)
		n.paper.DrawImage(newTextBoxImg, op)
	}
	// Set new cursor position
	n.cursorPosition = relCurPos
}

// Print prints a string in the selected textbox
func (n *Nimbus) Print(text string) {
	// turn off cursor until finished printing
	originalCursorMode := n.cursorMode
	n.cursorMode = -1
	for _, c := range text {
		n.Put(int(c))
	}
	// Send carriage return
	n.Put(13)
	// reset cursor mode to whatever it was before
	n.cursorMode = originalCursorMode
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
func (n *Nimbus) Input(prompt string, prepopulateBuffer string) string {
	// Print the prompt message without carriage return
	for _, c := range prompt {
		n.Put(int(c))
	}

	// Initialize buffer and print it
	var buffer []int
	for _, c := range prepopulateBuffer {
		n.Put(int(c))
		buffer = append(buffer, int(c))
	}
	bufferPosition := len(buffer)
	maxBufferSize := 255

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
			n.cursorPosition.col = width + 1
			if andDelete {
				n.Put(32)
				n.cursorPosition.row--
				n.cursorPosition.col = width + 1
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

	// now loop to received and edit the input string until enter is pressed
	for {
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
				// ENTER pressed so echo buffer beyong current position
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
				if bufferPosition > 1 {
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
