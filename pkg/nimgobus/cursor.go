package nimgobus

// drawCursor draws the cursor
func (n *Nimbus) drawCursor() {
	// Set up cursor
	var charPixels [][]int
	if n.deleteMode {
		charPixels = n.deleteModeCursorImage
	} else {
		switch n.cursorCharset {
		case 0:
			charPixels = n.charImages0[n.cursorChar]
		case 1:
			charPixels = n.charImages1[n.cursorChar]
		}
	}
	// Pick the textbox, get x, y coordinate of cursor and draw the char
	box := n.textBoxes[n.selectedTextBox]
	relCurPos := n.cursorPosition
	var absCurPos colRow // we need the absolute cursor position
	absCurPos.col = relCurPos.col + box.col1 - 1
	absCurPos.row = relCurPos.row + box.row1 - 1
	curX, curY := n.convertColRow(absCurPos)
	// cleanup cursor if disables and/or skip to next iteration
	n.muCursorFlash.Lock()
	if n.cursorFlash {
		n.writeSpriteToOverlay(Sprite{pixels: charPixels, x: curX, y: curY, colour: n.penColour, over: false})
	}
	n.muCursorFlash.Unlock()
}

// AdvanceCursor moves the cursor forward and handles line feeds and carriage returns
func (n *Nimbus) AdvanceCursor(forceCarriageReturn bool) {

	// Pick the textbox
	box := n.textBoxes[n.selectedTextBox]
	width := box.col2 - box.col1 // width and height in chars
	height := box.row2 - box.row1

	// Get relative cursor position and move cursor forward
	relCurPos := n.cursorPosition
	relCurPos.col++

	// Carriage return?
	if relCurPos.col > width+1 || forceCarriageReturn {
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
		y1 += 10
		x2 += 8
		// We have to manipulate videoMemory itself next, so force redraw and get the drawQueue lock
		n.ForceRedraw()
		n.muDrawQueue.Lock()
		n.muVideoMemory.Lock()
		// Copy the textbox segment of videoMemory
		textBoxImg := make2dArray((x2-x1)+1, y1-y2)
		for y := y2; y < y1; y++ {
			textBoxImg[(y - y2)] = n.videoMemory[y][x1:x2]
		}
		// Empty paper on bottom row of textbox
		paperImg := make2dArray((x2-x1)+9, 10)
		for x := x1; x <= x2; x++ {
			for y := 0; y < 10; y++ {
				//log.Printf("x: %d, y: %d, x1: %d, y1: %d, x2: %d, y2: %d, x2-x1: %d", x2, y, x1, y1, x2, y2, x2-x1)
				paperImg[y][x] = n.paperColour
			}
		}
		n.muVideoMemory.Unlock()
		n.muDrawQueue.Unlock()
		n.drawSprite(Sprite{pixels: textBoxImg[10:], x: x1, y: y2 + 10, colour: -1, over: true})
		n.drawSprite(Sprite{pixels: paperImg, x: x1, y: y2, colour: -1, over: true})
	}
	// Set new cursor position
	n.cursorPosition = relCurPos
}

// SetWriting selects a textbox if only 1 parameter is passed (index), or
// defines a textbox if 5 parameters are passed (index, col1, row1, col2,
// row2)
func (n *Nimbus) SetWriting(p ...int) {
	// Validate number of parameters
	if len(p) != 1 && len(p) != 5 {
		// invalid
		panic("SetWriting accepts either 1 or 5 parameters")
	}
	if len(p) == 1 {
		// Select textbox - validate choice first then set it
		if p[0] < 0 || p[0] > 10 {
			panic("SetWriting index out of range")
		}
		oldTextBox := n.selectedTextBox
		n.selectedTextBox = p[0]
		// Set cursor position to 1,1 if different textbox selected
		if oldTextBox != n.selectedTextBox {
			n.cursorPosition = colRow{1, 1}
		}
		return
	}
	// Otherwise define textbox if index is not 0
	if p[0] == 0 {
		panic("SetWriting cannot define index zero")
	}
	// Validate column and row values
	for i := 1; i < 5; i++ {
		if p[i] < 0 {
			panic("Negative row or column values are not allowed")
		}
	}
	if p[2] > 25 || p[4] > 25 {
		panic("Row values above 25 are not allowed")
	}
	maxColumns := n.mode
	if p[1] > maxColumns || p[3] > maxColumns {
		panic("Column value out of range for this screen mode")
	}
	// Validate passed - set bottomLeft and topRight colrows
	var upper, lower, left, right int
	if p[1] < p[3] {
		left = p[1]
		right = p[3]
	} else {
		left = p[3]
		right = p[1]
	}
	if p[2] < p[4] {
		upper = p[2]
		lower = p[4]
	} else {
		upper = p[4]
		lower = p[2]
	}
	// Set textbox
	n.textBoxes[p[0]] = textBox{left, upper, right, lower}

	return
}
