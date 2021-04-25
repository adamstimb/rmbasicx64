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
		n.writeSpriteToOverlay(Sprite{charPixels, curX, curY, n.penColour, false})
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
		textBoxImg := make2dArray(x2-x1, y1-y2)
		for y := y2; y < y1; y++ {
			textBoxImg[y] = n.videoMemory[y][:]
		}
		// Empty paper on bottom row of textbox
		paperImg := make2dArray(x2-x1, 10)
		for x := x1; x < x2; x++ {
			for y := 0; y < 10; y++ {
				paperImg[y][x] = n.paperColour
			}
		}
		n.muVideoMemory.Unlock()
		n.muDrawQueue.Unlock()
		n.drawSprite(Sprite{textBoxImg[10:], x1, y2 + 10, -1, true})
		n.drawSprite(Sprite{paperImg, x1, y2, -1, true})
	}
	// Set new cursor position
	n.cursorPosition = relCurPos
}
