package nimgobus

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// getPixel waits until the drawQueue is empty gets the colour of a pixel in the video memory
func (n *Nimbus) GetPixel(x, y int) (colour int) {
	n.muDrawQueue.Lock()
	colour = n.videoMemory[x][y]
	n.muDrawQueue.Unlock()
	return colour
}

// resizeSprite does a nearest-neighbour resize of a sprite
func (n *Nimbus) resizeSprite(thisSprite Sprite, newWidth, newHeight int) Sprite {
	img := thisSprite.pixels
	newImg := make2dArray(newWidth, newHeight)
	imgWidth := len(img[0])
	imgHeight := len(img)
	xScale := float64(imgWidth) / float64(newWidth)
	yScale := float64(imgHeight) / float64(newHeight)
	for y2 := 0; y2 < newHeight; y2++ {
		for x2 := 0; x2 < newWidth; x2++ {
			x1 := int(math.Floor((float64(x2) + 0.5) * xScale))
			y1 := int(math.Floor((float64(y2) + 0.5) * yScale))
			newImg[y2][x2] = img[y1][x1]
		}
	}
	return Sprite{pixels: newImg, x: thisSprite.x, y: thisSprite.y, colour: thisSprite.colour, over: thisSprite.over}
}

// rotateSprite90 rotates a sprite 90 degrees counterclockwise
func (n *Nimbus) rotateSprite90(thisSprite Sprite) Sprite {
	img := thisSprite.pixels
	imgWidth := len(img[0])
	imgHeight := len(img)
	newWidth := imgHeight
	newHeight := imgWidth
	newImg := make2dArray(newWidth, newHeight)
	for x1 := 0; x1 < imgWidth; x1++ {
		for y1 := 0; y1 < imgHeight; y1++ {
			x2 := y1
			y2 := (newHeight - 1) - x1
			newImg[y2][x2] = img[y1][x1]
		}
	}
	return Sprite{pixels: newImg, x: thisSprite.x, y: thisSprite.y, colour: thisSprite.colour, over: thisSprite.over}
}

// rotateSprite rotates a sprite 90 degress counterclockwise r times
func (n *Nimbus) rotateSprite(thisSprite Sprite, r int) Sprite {
	for i := 0; i < r; i++ {
		thisSprite = n.rotateSprite90(thisSprite)
	}
	return thisSprite
}

// applyDrawingbox applies a drawing box to the sprite (translate + truncate)
func (n *Nimbus) applyDrawingbox(thisSprite Sprite, d int) Sprite {

	// Get translation vector
	xt := n.drawingBoxes[d].x1
	yt := n.drawingBoxes[d].y1
	// Truncate sprite
	img := thisSprite.pixels
	imgWidth := len(img[0])
	imgHeight := len(img)
	newImg := make2dArray(imgWidth, imgHeight)
	for x := 0; x < len(img[0]); x++ {
		for y := 0; y < len(img); y++ {
			absoluteX := thisSprite.x + x + xt
			absoluteY := thisSprite.y + y + yt
			if absoluteX >= xt && absoluteY >= yt && absoluteX <= n.drawingBoxes[d].x2 && absoluteY <= n.drawingBoxes[d].y2 {
				newImg[y][x] = img[y][x]
			}
		}
	}
	// Return truncated sprite with applied translation
	return Sprite{pixels: newImg, x: thisSprite.x + xt, y: thisSprite.y + yt, colour: thisSprite.colour, over: thisSprite.over}
}

// drawSprite waits until the drawQueue is unlocked then adds a sprite for drawing to the drawQueue
func (n *Nimbus) drawSprite(thisSprite Sprite) {
	n.redrawComplete = false
	n.muDrawQueue.Lock()
	// add to queue and unlock
	n.drawQueue = append(n.drawQueue, thisSprite)
	n.muDrawQueue.Unlock()
}

// writeSprite writes a sprite directly to videoMemory
func (n *Nimbus) writeSprite(thisSprite Sprite) {
	// assumes drawQueue is locked!
	n.muVideoMemory.Lock()
	// convert coordinates and get dimensions
	imageMemorySize := n.videoImage.Bounds() // should get size of videoMemory instead
	spriteOffsetX := thisSprite.x
	spriteOffsetY := imageMemorySize.Max.Y - thisSprite.y - len(thisSprite.pixels)
	spriteWidth := len(thisSprite.pixels[0])
	spriteHeight := len(thisSprite.pixels)
	// make sure we truncate the sprite if it doesn't fit on screen
	var xLimit, yLimit int
	if spriteOffsetX+spriteWidth >= 640 {
		xLimit = 640
	} else {
		xLimit = spriteOffsetX + spriteWidth
	}
	if spriteOffsetY+spriteHeight >= 250 {
		yLimit = 250
	} else {
		yLimit = spriteOffsetY + spriteHeight
	}
	// write the sprite in over==true mode (default) or over==false mode (XOR)
	if thisSprite.over {
		// over==true
		spriteX := 0
		for x := spriteOffsetX; x < xLimit; x++ {
			spriteY := 0
			// don't draw if it's off the left-hand side
			if x < 0 {
				spriteX++
				continue
			}
			for y := spriteOffsetY; y < yLimit; y++ {
				// don't draw if it's below the screen
				if y < 0 {
					spriteY++
					continue
				}
				// colour >= 0 represents a b+w sprite so colourise it with specified colour
				// otherwise use the colour given by the pixel
				if thisSprite.colour >= 0 {
					if thisSprite.pixels[spriteY][spriteX] == 1 {
						n.videoMemory[y][x] = thisSprite.colour
					}
				} else {
					n.videoMemory[y][x] = thisSprite.pixels[spriteY][spriteX]
				}
				spriteY++
			}
			spriteX++
		}
	} else {
		// over==false, i.e XOR mode
		spriteX := 0
		for x := spriteOffsetX; x < xLimit; x++ {
			spriteY := 0
			for y := spriteOffsetY; y < yLimit; y++ {
				if thisSprite.pixels[spriteY][spriteX] == 1 {
					n.videoMemory[y][x] = n.videoMemory[y][x] ^ thisSprite.colour
				}
				spriteY++
			}
			spriteX++
		}
	}
	n.muVideoMemory.Unlock()
}

// writeSpriteToOverlay writes a sprite directly to videoMemoryOverlay
func (n *Nimbus) writeSpriteToOverlay(thisSprite Sprite) {
	// assumes drawQueue is locked!
	n.muVideoMemoryOverlay.Lock()
	// convert coordinates and get dimensions
	imageMemoryOverlaySize := n.videoImage.Bounds() // should get size of videoMemory instead
	spriteOffsetX := thisSprite.x
	spriteOffsetY := imageMemoryOverlaySize.Max.Y - thisSprite.y - len(thisSprite.pixels)
	spriteWidth := len(thisSprite.pixels[0])
	spriteHeight := len(thisSprite.pixels)
	// make sure we truncate the sprite if it doesn't fit on screen
	var xLimit, yLimit int
	if spriteOffsetX+spriteWidth >= 640 {
		xLimit = 640
	} else {
		xLimit = spriteOffsetX + spriteWidth
	}
	if spriteOffsetY+spriteHeight >= 250 {
		yLimit = 250
	} else {
		yLimit = spriteOffsetY + spriteHeight
	}
	// write the sprite in over==true mode (default) or over==false mode (XOR)
	if thisSprite.over {
		// over==true
		spriteX := 0
		for x := spriteOffsetX; x < xLimit; x++ {
			spriteY := 0
			// don't draw if it's off the left-hand side
			if x < 0 {
				spriteX++
				continue
			}
			for y := spriteOffsetY; y < yLimit; y++ {
				// don't draw if it's below the screen
				if y < 0 {
					spriteY++
					continue
				}
				// colour > 0 represents a b+w sprite so colourise it with specified colour
				// otherwise use the colour given by the pixel
				if thisSprite.colour >= 0 {
					if thisSprite.pixels[spriteY][spriteX] == 1 {
						n.videoMemoryOverlay[y][x] = thisSprite.colour
					}
				} else {
					n.videoMemoryOverlay[y][x] = thisSprite.pixels[spriteY][spriteX]
				}
				spriteY++
			}
			spriteX++
		}
	} else {
		// over==false, i.e XOR mode
		spriteX := 0
		for x := spriteOffsetX; x < xLimit; x++ {
			spriteY := 0
			// don't draw if it's off the left-hand side
			if x < 0 {
				spriteX++
				continue
			}
			for y := spriteOffsetY; y < yLimit; y++ {
				// don't draw if it's below the screen
				if y < 0 {
					spriteY++
					continue
				}
				if thisSprite.pixels[spriteY][spriteX] == 1 {
					n.videoMemoryOverlay[y][x] = n.videoMemoryOverlay[y][x] ^ thisSprite.colour
				}
				spriteY++
			}
			spriteX++
		}
	}
	n.muVideoMemoryOverlay.Unlock()
}

// putSpriteOn2dArray copies a sprite to a 2d array
func (n *Nimbus) putSpriteOn2dArray(targetArray [][]int, thisSprite Sprite) {
	spriteOffsetX := thisSprite.x
	spriteOffsetY := thisSprite.y
	spriteWidth := len(thisSprite.pixels[0])
	spriteHeight := len(thisSprite.pixels)
	targetArrayWidth := len(targetArray[0])
	targetArrayHeight := len(targetArray)
	// make sure we truncate the sprite if it doesn't fit on the array
	var xLimit, yLimit int
	if spriteOffsetX+spriteWidth > targetArrayWidth-1 {
		xLimit = targetArrayWidth - 1
	} else {
		xLimit = spriteOffsetX + spriteWidth
	}
	if spriteOffsetY+spriteHeight > targetArrayHeight-1 {
		yLimit = targetArrayHeight - 1
	} else {
		yLimit = spriteOffsetY + spriteHeight
	}
	// write the sprite in over==true mode (default) or over==false mode (XOR)
	if thisSprite.over {
		// over==true
		spriteX := 0
		for x := spriteOffsetX; x <= xLimit; x++ {
			spriteY := 0
			for y := spriteOffsetY; y <= yLimit; y++ {
				if thisSprite.pixels[spriteY][spriteX] == 1 {
					targetArray[y][x] = 1
				}
				spriteY++
			}
			spriteX++
		}
	} else {
		// over==false, i.e XOR mode
		spriteX := 0
		for x := spriteOffsetX; x <= xLimit; x++ {
			spriteY := 0
			for y := spriteOffsetY; y <= yLimit; y++ {
				targetArray[y][x] = n.videoMemory[y][x] ^ thisSprite.pixels[spriteY][spriteX]
				spriteY++
			}
			spriteX++
		}
	}
}

// updateVideoImage regenerates the ebiten Image from the videoMemory
func (n *Nimbus) updateVideoImage() {
	// assume the drawQueue is locked!
	maxX, maxY := n.videoImage.Size()
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			n.videoImage.Set(x, y, n.basicColours[n.palette[n.videoMemoryOverlay[y][x]]])
		}
	}
}

// updateVideoMemory writes all sprites in the drawQueue to videoMemory
func (n *Nimbus) updateVideoMemory() {
	n.muDrawQueue.Lock()
	for _, thisSprite := range n.drawQueue {
		n.writeSprite(thisSprite)
	}
	// update video overlay
	for y := 0; y < 250; y++ {
		n.videoMemoryOverlay[y] = n.videoMemory[y]
	}
	// draw cursor on overlay if enabled
	if n.cursorFlashEnabled {
		n.drawCursor()
	}
	// and update videoImage
	n.updateVideoImage()
	n.muDrawQueue.Unlock()
}

// redraw redraws the monitor
func (n *Nimbus) redraw() {
	n.redrawComplete = false
	// rescale videoImage, overlay it on the border then update monitor
	sizeX, sizeY := n.videoImage.Size()
	scaleX := 640.0 / float64(sizeX)
	scaleY := 500.0 / float64(sizeY)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(n.borderSize), float64(n.borderSize))
	n.borderImage.DrawImage(n.videoImage, op)
	op = &ebiten.DrawImageOptions{}
	n.Monitor.DrawImage(n.borderImage, op)
	n.redrawComplete = true
}

// ForceRedraw forces the monitor to redraw in the case of, for example, change of border colour
func (n *Nimbus) ForceRedraw() {
	n.redrawComplete = false
	// send an arbitrary sprite to be drawn off-screen
	blank := make2dArray(10, 10)
	blankSprite := Sprite{blank, 640, 250, 0, true}
	n.drawSprite(blankSprite)
	for !n.redrawComplete {
		time.Sleep(1 * time.Microsecond)
	}
}
