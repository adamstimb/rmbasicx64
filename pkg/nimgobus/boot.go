package nimgobus

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/elastic/go-sysinfo"
	"github.com/shirou/gopsutil/host"
)

// randDelay delays for a random number of milliseconds within limits
func randDelay(min, max int) {
	delay := time.Duration(rand.Intn(max-min)+min) * time.Millisecond
	time.Sleep(delay)
}

// Boot simulates the RM Nimbus "Welcome" boot screen and operating system
// loading workflow.  The original Nimbus would also display system info, such
// as firmware version, serial number, memory, etc.  Nimgobus immitates this
// using the Go compiler version as the firmware version, and displays the
// actual physical and virtual memory size.  Serial number is a string constant
// as is the serial number of the Nimbus that provided the ROM dump for the
// emulation on MAME, from which various bits and pieces were reversed
// engineering for nimgobus.
func (n *Nimbus) Boot() {
	drawBackground(n)
	plotOpts := PlotOptions{
		Font:      1,
		Brush:     3,
		SizeX:     1,
		SizeY:     1,
		Direction: 0,
		Over:      -1,
	}
	n.Plot(plotOpts, "Please supply an operating system", 188, 100)
	randDelay(1000, 2000)
	plotOpts.Brush = 1
	n.Plot(plotOpts, "Please supply an operating system", 188, 100)
	plotOpts.Brush = 3
	n.Plot(plotOpts, "Looking for an operating system - please wait", 140, 100)
	randDelay(2500, 3500)
	plotOpts.Brush = 1
	n.Plot(plotOpts, "Looking for an operating system - please wait", 140, 100)
	plotOpts.Brush = 3
	n.Plot(plotOpts, "Loading operating system", 224, 100)
	randDelay(2500, 3500)
	n.SetMode(80)
	n.SetColour(0, 0, 0, 0)
	n.SetCursor(1)
	// Mimick an old MS-DOS boot
	// Print platform information
	platform, _, version, _ := host.PlatformInformation()
	randDelay(1000, 2000)
	info := fmt.Sprintf("      %s - Version %s      ", platform, version)
	copyright := ""
	if strings.Contains(platform, "darwin") {
		copyright = "Copyright (c) Apple Inc. All rights reserved."
	}
	n.SetPaper(3)
	n.Print(info)
	n.Put(13)
	n.Print(info)
	n.Put(13)
	n.SetPen(0)
	n.Print(info)
	n.Put(13)
	n.SetPen(3)
	n.Print(info)
	n.Put(13)
	n.Print(info)
	n.SetPaper(0)
	n.Put(13)
	n.Put(13)
	n.Print(copyright)
	n.Put(13)
	n.Put(13)
	randDelay(2000, 3000)
	n.Print("C>")
	randDelay(500, 1000)
	// Try to get directory of executable (doesn't matter if it doesn't)
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Printf("Error resolving directory of executable: %v", err)
		exeDir = "rmbasicx64"
	}
	n.Print(fmt.Sprintf("cd %s", exeDir))
	n.Put(13)
	randDelay(1000, 2000)
	n.Print("C>")
	randDelay(500, 1000)
	n.Print("rmbasicx64")
	n.Put(13)
	randDelay(3000, 4000)
}

// convert bytes to Mb
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// drawBackground draws the background of the Welcome screen
func drawBackground(n *Nimbus) {

	// Collect system info
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	firmwareVersion := runtime.Version() // Use Go version instead
	firmwareVersion = firmwareVersion[2:]
	if len(firmwareVersion) > 8 {
		firmwareVersion = firmwareVersion[:8]
	}
	host, err := sysinfo.Host()
	if err != nil {
		panic("Could not detect system information")
	}
	firmwareVersion = fmt.Sprintf("Firmware version: %s", firmwareVersion)
	serialNumber := "Serial number:  21/06809" // In honour of whichever physical machine donate its ROM to MAME
	memInfo, err := host.Memory()
	mainMemSize := fmt.Sprintf("main    memory size %7d Mbytes", bToMb(memInfo.Available))
	virtualMemSize := fmt.Sprintf("virtual memory size %7d Mbytes", bToMb(memInfo.VirtualTotal))
	totalMemSize := fmt.Sprintf("total   memory size %7d Mbytes", bToMb(memInfo.Available+memInfo.VirtualTotal))

	// Red frame, light blue paper, Nimbus logo in a red frame
	n.SetMode(80)
	n.SetCursor(0)
	n.SetColour(0, 0, 0, 0)
	n.SetSound(true)
	n.Bell()
	n.SetColour(1, 9, 0, 9)
	n.SetPaper(1)
	n.SetBorder(1)
	n.Cls()

	areaOpts := AreaOptions{
		Brush: 2,
		Over:  -1,
	}
	n.Area(areaOpts, []XyCoord{{0, 0}, {639, 0}, {639, 249}, {0, 249}, {0, 0}})
	areaOpts.Brush = 1
	n.Area(areaOpts, []XyCoord{{3, 2}, {636, 2}, {636, 247}, {3, 247}, {3, 2}})
	xl := 10
	yl := 212
	n.PlonkLogo(xl, yl)
	lineOpts := LineOptions{
		Brush: 2,
		Over:  -1,
	}
	//n.Line(lineOpts, xl, yl, xl+304, yl, xl+304, yl+32, xl, yl+32, xl, yl)
	n.Line(lineOpts, []XyCoord{{xl, yl}, {xl + 304, yl}, {xl + 304, yl + 32}, {xl, yl + 32}, {xl, yl}})
	plotOpts := PlotOptions{
		SizeX: 3, SizeY: 3, Font: 1, Direction: 0, Over: -1,
	}

	// Welcome
	n.Plot(plotOpts, "Welcome", 238, 145)
	plotOpts.Brush = 2
	n.Plot(plotOpts, "Welcome", 236, 147)

	// Firmware version and serial number
	areaOpts.Brush = 2
	//n.Area(areaOpts, 393, 4, 632, 4, 632, 30, 393, 30, 393, 4)
	n.Area(areaOpts, []XyCoord{{393, 4}, {632, 4}, {632, 30}, {393, 30}, {393, 4}})
	areaOpts.Brush = 3
	//n.Area(areaOpts, 395, 5, 629, 5, 629, 29, 395, 29, 395, 5)
	n.Area(areaOpts, []XyCoord{{395, 5}, {629, 5}, {629, 29}, {395, 29}, {395, 5}})
	plotOpts.Brush = 0
	plotOpts.SizeX = 1
	plotOpts.SizeY = 1
	n.Plot(plotOpts, firmwareVersion, 400, 17)
	n.Plot(plotOpts, serialNumber, 400, 7)

	// Memory
	plotOpts.Brush = 0
	n.Plot(plotOpts, mainMemSize, 15, 25)
	n.Plot(plotOpts, virtualMemSize, 15, 15)
	n.Plot(plotOpts, totalMemSize, 15, 5)
}
