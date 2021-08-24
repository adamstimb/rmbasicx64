package nimgobus

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// SetMouse turns the mouse monitor on (true) or turns it off (false)
func (n *Nimbus) SetMouse(mouseOn bool) {
	if mouseOn {
		go n.mouseMonitor()
	} else {
		n.MouseOn = false
	}
}

// mouseMonitor updates the mouse status.  It should be called be a go command.
func (n *Nimbus) mouseMonitor() {
	n.MouseOn = true
	for n.MouseOn {
		x, y := ebiten.CursorPosition()
		var b int
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			b = 1
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b = 2
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b = 3
		}
		// Scale x, y to Nimbus screen
		x -= n.PaddingX
		y -= n.PaddingY
		x = int(float64(x) / n.Scale)
		y = int(float64(y) / n.Scale)
		x -= n.borderSize
		y -= n.borderSize
		videoWidth, _ := n.videoImage.Size()
		if videoWidth == 640 {
			y = y / 2
		} else {
			x = x / 2
			y = y / 2
		}
		y = 250 - y // Flip vertical
		// Clamp values
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}
		if x > videoWidth {
			x = videoWidth
		}
		if y > 250 {
			y = 250
		}
		// Update values on Nimbogbus
		n.MouseX = x
		n.MouseY = y
		n.MouseButton = b
		time.Sleep(50 * time.Millisecond)
	}
}
