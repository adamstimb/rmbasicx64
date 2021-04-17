package nimgobus

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// SetMode sets the screen mode. 40 is low-resolution, high-colour mode (320x250) and
// 80 is high-resolutions, low-colour mode (640x250)
func (n *Nimbus) SetMode(columns int) {
	if columns != 40 && columns != 80 {
		// RM Basic would just pass if an invalid column number was given so
		// we'll do the same
		return
	}
	if columns == 40 {
		// low-resolution, high-colour mode (320x250)
		n.paper = ebiten.NewImage(320, 250)
		n.paperColour = 0
		n.borderColour = 0
		n.penColour = 15
		n.palette = n.defaultLowResPalette
		// reinit border image and fill with new colour
		n.borderImage = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
		n.borderImage.Fill(n.convertColour(n.borderColour))
	}
	if columns == 80 {
		// high-resolutions, low-colour mode (640x250)
		n.paper = ebiten.NewImage(640, 250)
		n.palette = n.defaultHighResPalette
		n.paperColour = 0
		n.borderColour = 0
		n.penColour = 3
		// reinit border image and fill with new colour
		n.borderImage = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
		n.borderImage.Fill(n.convertColour(n.borderColour))
	}
	n.cursorPosition = colRow{1, 1}              // Relocate cursor
	n.paper.Fill(n.convertColour(n.paperColour)) // Apply paper colour
	// Redefine textboxes, imageBlocks and clear screen
	for i := 0; i < 10; i++ {
		n.textBoxes[i] = textBox{1, 1, columns, 25}
	}
	n.imageBlocks = [16]*ebiten.Image{}
	n.Cls()
}

// SetBorder sets the border colour
func (n *Nimbus) SetBorder(c int) {
	n.validateColour(c)
	n.borderColour = c
}

// PlonkChar plots a character on the paper at a given location
func (n *Nimbus) PlonkChar(c, x, y, colour int) {
	n.drawChar(n.paper, c, x, y, colour, n.charset)
}

// Mode returns the current screen mode (40 column or 80 column)
func (n *Nimbus) Mode() int {
	width, _ := n.paper.Size()
	if width == 320 {
		return 40 // low-res mode 40
	}
	if width == 640 {
		return 80 // high-res mode 80
	}
	return 0 // this never happens
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
		// and return
		if p[0] < 0 || p[0] > 9 {
			panic("SetWriting index out of range")
		}
		n.selectedTextBox = p[0]
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
	maxColumns := n.Mode()
	if p[1] > maxColumns || p[3] > maxColumns {
		panic("Column value out of range for this screen mode")
	}
	// Validate passed - set the textbox
	n.textBoxes[p[0]] = textBox{p[1], p[2], p[3], p[4]}
	return
}

// SetPaper sets the paper colour
func (n *Nimbus) SetPaper(c int) {
	n.validateColour(c)
	n.paperColour = c
}

// Cls clears the selected textbox if no parameters are passed, or clears another
// textbox if one parameter is passed
func (n *Nimbus) Cls(p ...int) {
	// Validate number of parameters
	if len(p) != 0 && len(p) != 1 {
		// invalid
		panic("Cls accepts either 0 or 1 parameters")
	}
	// Pick the textbox
	var box textBox
	if len(p) == 0 {
		// No parameters passed so clear currently selected textbox
		box = n.textBoxes[n.selectedTextBox]
	} else {
		// One parameter passed so chose another textbox
		box = n.textBoxes[p[0]]
	}
	// Define bounding rectangle for the textbox
	x1, y1 := n.convertColRow(colRow{box.col1, box.row1})
	x2, y2 := n.convertColRow(colRow{box.col2, box.row2})
	x2 += 8
	y2 += 10
	// Create temp image and fill it with paper colour, then paste on the
	// paper
	img := ebiten.NewImage(int(x2-x1), int(y2-y1))
	img.Fill(n.convertColour(n.paperColour))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x1, y1)
	n.paper.DrawImage(img, op)
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

// SetPen sets the pen colour
func (n *Nimbus) SetPen(c int) {
	n.validateColour(c)
	n.penColour = c
}

// SetColour assigns one of the basic colours to a slot in the current palette
func (n *Nimbus) SetColour(paletteSlot, basicColour int) {
	// Validate basicColour and paletteSlot
	if basicColour < 0 || basicColour > 16 {
		panic("basicColour is out of range")
	}
	n.validateColour(paletteSlot)
	// Validation passed, assign colour
	n.palette[paletteSlot] = basicColour
}

// SetCharset selected either the default Nimbus charset (0) or the alternative
// charset (1)
func (n *Nimbus) SetCharset(i int) {
	// Validate
	if i != 0 && i != 1 {
		panic("Invalid charset")
	}
	// set charset
	n.charset = i
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
	n.cursorMode = p[0]
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
