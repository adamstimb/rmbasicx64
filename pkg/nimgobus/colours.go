package nimgobus

import "image/color"

var (
	basicColours = []color.RGBA{
		{0x00, 0x00, 0x00, 0xff}, // black
		{0x00, 0x00, 0xaa, 0xff}, // dark blue
		{0xaa, 0x00, 0x00, 0xff}, // dark red
		{0xaa, 0x00, 0xaa, 0xff}, // purple
		{0x00, 0xaa, 0x00, 0xff}, // dark green
		{0x00, 0xaa, 0xaa, 0xff}, // dark cyan
		{0xaa, 0x54, 0x00, 0xff}, // brown
		{0xaa, 0xaa, 0xaa, 0xff}, // light grey
		{0x54, 0x54, 0x54, 0xff}, // dark grey
		{0x54, 0x54, 0xff, 0xff}, // light blue
		{0xff, 0x54, 0x54, 0xff}, // light red
		{0xff, 0x54, 0xff, 0xff}, // light purple
		{0x54, 0xff, 0x54, 0xff}, // light green
		{0x54, 0xff, 0xff, 0xff}, // light cyan
		{0xff, 0xff, 0x54, 0xff}, // yellow
		{0xff, 0xff, 0xff, 0xff}, // white
	}
	defaultLowResPalette = []int{
		0,  // black
		1,  // dark blue
		2,  // dark red
		3,  // purple
		4,  // dark green
		5,  // dark cyan
		6,  // brown
		7,  // light grey
		8,  // dark grey
		9,  // light blue
		10, // light red
		11, // light purple
		12, // light green
		13, // light cyan
		14, // yellow
		15, // white
	}
	defaultHighResPalette = []int{
		1,  // dark blue
		4,  // dark green
		10, // light red
		15, // white
	}
	defaultPointsStyles = [][][]int{
		{
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{1, 1, 1, 1, 1, 1, 1, 1, 1},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 1, 0, 0, 1, 0, 0, 1, 0},
			{0, 0, 1, 0, 1, 0, 1, 0, 0},
			{0, 0, 0, 1, 1, 1, 0, 0, 0},
			{1, 1, 1, 1, 1, 1, 1, 1, 1},
			{0, 0, 0, 1, 1, 1, 0, 0, 0},
			{0, 0, 1, 0, 1, 0, 1, 0, 0},
			{0, 1, 0, 0, 1, 0, 0, 1, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
		},
		{
			{0, 0, 0, 1, 1, 1, 0, 0, 0},
			{0, 0, 1, 0, 0, 0, 1, 0, 0},
			{0, 1, 0, 0, 0, 0, 0, 1, 0},
			{1, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 1},
			{0, 1, 0, 0, 0, 0, 0, 1, 0},
			{0, 0, 1, 0, 0, 0, 1, 0, 0},
			{0, 0, 0, 1, 1, 1, 0, 0, 0},
		},
		{
			{1, 0, 0, 0, 0, 0, 0, 0, 1},
			{0, 1, 0, 0, 0, 0, 0, 1, 0},
			{0, 0, 1, 0, 0, 0, 1, 0, 0},
			{0, 0, 0, 1, 0, 1, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 1, 0, 1, 0, 0, 0},
			{0, 0, 1, 0, 0, 0, 1, 0, 0},
			{0, 1, 0, 0, 0, 0, 0, 1, 0},
			{1, 0, 0, 0, 0, 0, 0, 0, 1},
		},
	}
)

// initialize colour flash settings so that nothing flashes and the flash colour is the same
// as the normal colour
func (n *Nimbus) initColourFlashSettings() {
	n.muColourFlashSettings.Lock()
	n.colourFlashSettings = []colourFlashSetting{}
	for i := 0; i < 16; i++ {
		n.colourFlashSettings = append(n.colourFlashSettings, colourFlashSetting{speed: 0, flashColour: i})
	}
	n.muColourFlashSettings.Unlock()
}
