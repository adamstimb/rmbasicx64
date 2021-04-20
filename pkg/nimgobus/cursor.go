package nimgobus

// drawCursor draws the cursor
func (n *Nimbus) drawCursor() {
	n.muCursorFlash.Lock()
	var charPixels [][]int
	switch n.cursorCharset {
	case 0:
		charPixels = n.charImages0[n.cursorChar]
	case 1:
		charPixels = n.charImages1[n.cursorChar]
	}
	n.drawSprite(Sprite{charPixels, 0, 240, n.penColour, false})
	n.cursorIsVisible = !n.cursorIsVisible
	n.cursorFlash = false
	n.muCursorFlash.Unlock()
}

// HideCursor makes the cursor invisible and stops it flashing
func (n *Nimbus) HideCursor() {
	n.cursorFlashEnabled = false
	if n.cursorIsVisible {
		// cursor is still visible so redraw it once to hide it
		n.drawCursor()
	}
}

// ShowCursor makes the cursor visible and flashing (if flashing enabled)
func (n *Nimbus) ShowCursor() {
	n.cursorFlashEnabled = true
}
