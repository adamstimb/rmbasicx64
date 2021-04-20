package nimgobus

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// getPixel waits until the drawQueue is empty gets the colour of a pixel in the video memory
func (n *Nimbus) GetPixel(x, y int) (colour int) {
	n.muDrawQueue.Lock()
	colour = n.videoMemory[x][y]
	n.muDrawQueue.Unlock()
	return colour
}

// GetPixels waits until the drawQueue is empty then returns a copy of the video memory
func (n *Nimbus) GetPixels() (videoMemoryCopy [250][640]int) {
	n.muDrawQueue.Lock()
	copy(videoMemoryCopy[:], n.videoMemory[:])
	n.muDrawQueue.Unlock()
	return videoMemoryCopy
}

// drawSprite waits until the drawQueue is unlocked then adds a sprite for drawing to the drawQueue
func (n *Nimbus) drawSprite(thisSprite Sprite) {
	n.muDrawQueue.Lock()
	// add to queue and unlock
	n.drawQueue = append(n.drawQueue, thisSprite)
	n.muDrawQueue.Unlock()
}

// writeSprite writes a sprite directly to videoMemory
func (n *Nimbus) writeSprite(thisSprite Sprite) {
	// assumes drawQueue is locked!
	// convert coordinates and get dimensions
	imageMemorySize := n.videoImage.Bounds()
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
				// colour > 0 represents a b+w sprite so colourise it with specified colour
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
			n.videoImage.Set(x, y, n.basicColours[n.palette[n.videoMemory[y][x]]])
		}
	}
}

// updateVideoMemory writes all sprites in the drawQueue to videoMemory
func (n *Nimbus) updateVideoMemory() {
	// assume the drawQueue is locked!
	for _, thisSprite := range n.drawQueue {
		n.writeSprite(thisSprite)
	}
	// and update videoImage
	n.updateVideoImage()
}

// redraw redraws the monitor
func (n *Nimbus) redraw() {
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
}

// ForceRedraw forces the monitor to redraw in the case of, for example, change of border colour
func (n *Nimbus) ForceRedraw() {
	// send an arbitrary sprite to be drawn off-screen
	blank := make2dArray(10, 10)
	blankSprite := Sprite{blank, 640, 250, 0, true}
	n.drawSprite(blankSprite)
}
