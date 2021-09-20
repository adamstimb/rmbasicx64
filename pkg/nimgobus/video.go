package nimgobus

import (
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// getPixel waits until the drawQueue is empty gets the colour of a pixel in the video memory
func (n *Nimbus) GetPixel(x, y int) (colour int) {
	drawQueueNotEmpty := true
	for drawQueueNotEmpty {
		n.muDrawQueue.Lock()
		lenDrawQueue := len(n.drawQueue)
		n.muDrawQueue.Unlock()
		if lenDrawQueue > 0 {
			time.Sleep(1 * time.Microsecond)
		} else {
			drawQueueNotEmpty = false
		}
	}
	n.muVideoMemory.Lock()
	colour = n.videoMemory[249-y][x]
	n.muVideoMemory.Unlock()
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

func (n *Nimbus) applyDrawingbox(thisSprite Sprite, d int) (Sprite, bool) {
	box := n.drawingBoxes[d]
	spriteWidth := len(thisSprite.pixels[0])
	spriteHeight := len(thisSprite.pixels)
	xt := box.x1
	yt := box.y1
	// Get max x,y within drawingbox
	maxX := box.x2 - box.x1
	maxY := box.y2 - box.y1
	// Reject off-screen sprites
	if thisSprite.y > maxY || thisSprite.x > maxX || thisSprite.y+spriteHeight < 0 || thisSprite.x+spriteWidth < 0 {
		return thisSprite, false
	}
	// Truncate as necessary
	if thisSprite.x < 0 {
		// truncate left
		chop := thisSprite.x * -1
		newImg := make2dArray(len(thisSprite.pixels[0])-chop, len(thisSprite.pixels))
		for x := 0; x < len(newImg[0]); x++ {
			for y := 0; y < len(newImg); y++ {
				newImg[y][x] = thisSprite.pixels[y][x+chop]
			}
		}
		thisSprite.pixels = newImg
		thisSprite.x = 0
	}
	if thisSprite.y+len(thisSprite.pixels) > maxY+1 {
		// truncate top
		chop := (thisSprite.y + len(thisSprite.pixels)) - (maxY + 1)
		newImg := make2dArray(len(thisSprite.pixels[0]), len(thisSprite.pixels)-chop)
		for x := 0; x < len(newImg[0]); x++ {
			for y := 0; y < len(newImg); y++ {
				newImg[y][x] = thisSprite.pixels[y+chop][x]
			}
		}
		thisSprite.pixels = newImg
	}
	if thisSprite.x+len(thisSprite.pixels[0]) > maxX+1 {
		// truncate right
		chop := thisSprite.x + len(thisSprite.pixels[0]) - (maxX + 1)
		newImg := make2dArray(len(thisSprite.pixels[0])-chop, len(thisSprite.pixels))
		for x := 0; x < len(newImg[0]); x++ {
			for y := 0; y < len(newImg); y++ {
				newImg[y][x] = thisSprite.pixels[y][x]
			}
		}
		thisSprite.pixels = newImg
	}
	if thisSprite.y < 0 {
		// truncate below
		chop := thisSprite.y * -1
		newImg := make2dArray(len(thisSprite.pixels[0]), len(thisSprite.pixels)-chop)
		for x := 0; x < len(newImg[0]); x++ {
			for y := 0; y < len(newImg); y++ {
				newImg[y][x] = thisSprite.pixels[y][x]
			}
		}
		thisSprite.pixels = newImg
		thisSprite.y = 0
	}
	return Sprite{pixels: thisSprite.pixels, x: thisSprite.x + xt, y: thisSprite.y + yt, colour: thisSprite.colour, over: thisSprite.over, fillStyle: thisSprite.fillStyle}, true
}

// drawSprite waits until the drawQueue is unlocked then adds a sprite for drawing to the drawQueue
func (n *Nimbus) drawSprite(thisSprite Sprite) {
	n.redrawComplete = false
	n.muDrawQueue.Lock()
	// add to queue and unlock
	n.drawQueue = append(n.drawQueue, thisSprite)
	n.muDrawQueue.Unlock()
}

// handleColourFlash handles flashing colours and returns the palette slot of the colour that
// should be displayed if flashing
func (n *Nimbus) handleColourFlash(c int) int {
	var retVal int
	n.muColourFlashSettings.Lock()
	flashSpeed := n.colourFlashSettings[c].speed
	flashColour := n.colourFlashSettings[c].flashColour
	n.muColourFlashSettings.Unlock()
	switch flashSpeed {
	case 0:
		// not flashing
		retVal = n.palette[c]
	case 1:
		// slow flash
		if n.colourFlash <= 1 {
			retVal = n.palette[c]
		} else {
			retVal = flashColour
		}
	case 2:
		// fast flash
		if n.colourFlash == 0 || n.colourFlash == 2 {
			retVal = n.palette[c]
		} else {
			retVal = flashColour
		}
	}
	return retVal
}

func overflow(n, max int) int {
	if n <= max {
		return n
	}
	return int(n - (max+1)*int(math.Floor(float64(n)/float64((max+1)))))
}

// handlePattern is used to figure out the colour of a pixel if a pattern is being used instead
// of a solid colour
func (n *Nimbus) handlePattern(x, y, c int, fillStyle FillStyle) int {
	// Handle fill style
	if fillStyle.Style == 2 && c < 128 {
		// use hatching
		x = overflow(x, 15)
		y = overflow(y, 15)
		hatching := n.hatchings[fillStyle.Hatching]
		if hatching[y][x] == 1 {
			return c
		} else {
			if fillStyle.Colour2 == -1 {
				return 0
			} else {
				return fillStyle.Colour2
			}
		}
	}
	if c < 128 {
		// is not a pattern
		return c
	}
	x = overflow(x, 3)
	y = overflow(y, 3)
	pattern := n.patterns[c-128]
	return pattern[y][x]
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
				// don't draw if it's above the screen (y is flipped here, remember)
				if y < 0 {
					spriteY++
					continue
				}
				// colour >= 0 represents a b+w sprite so colourise it with specified colour
				// otherwise use the colour given by the pixel
				if thisSprite.colour >= 0 {
					if thisSprite.pixels[spriteY][spriteX] == 1 {
						n.videoMemory[y][x] = n.handlePattern(x, y, thisSprite.colour, thisSprite.fillStyle)
					}
				} else {
					n.videoMemory[y][x] = n.handlePattern(x, y, thisSprite.pixels[spriteY][spriteX], thisSprite.fillStyle)
				}
				spriteY++
			}
			spriteX++
		}
	} else {
		// over==false, i.e XOR mode -- I think we disregard colour flashing in this mode
		spriteX := 0
		for x := spriteOffsetX; x < xLimit; x++ {
			spriteY := 0
			for y := spriteOffsetY; y < yLimit; y++ {
				// don't draw if it's off the left-hand side
				if x < 0 {
					spriteX++
					continue
				}
				if thisSprite.pixels[spriteY][spriteX] == 1 {
					n.videoMemory[y][x] = n.videoMemory[y][x] ^ n.handlePattern(x, y, thisSprite.colour, thisSprite.fillStyle)
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
	img := image.NewRGBA(image.Rect(0, 0, maxX, maxY))
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			if len(n.palette) <= n.videoMemoryOverlay[y][x] {
				img.Set(x, y, n.basicColours[n.palette[1]])
			} else {
				//img.Set(x, y, n.basicColours[n.palette[n.videoMemoryOverlay[y][x]]])
				//img.Set(x, y, n.basicColours[n.handleColourFlash(n.palette[n.videoMemoryOverlay[y][x]])])
				img.Set(x, y, n.basicColours[n.handleColourFlash(n.videoMemoryOverlay[y][x])])
			}
		}
	}
	n.videoImage = ebiten.NewImageFromImage(img)
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
	// flush drawQueue and update videoImage
	n.drawQueue = []Sprite{}
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
	blankSprite := Sprite{pixels: blank, x: 640, y: 250, colour: 0, over: true}
	n.drawSprite(blankSprite)
	for !n.redrawComplete {
		time.Sleep(1 * time.Microsecond)
	}
}
