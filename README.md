# ![RM BASICx64](https://github.com/adamstimb/rmbasicx64/blob/main/docs/branding/rmbasicx64BannerLarge.png)

_RM BASICx64 is a tribute project and is in no way linked to or endorsed by RM plc._

## About

RM BASICx64 is a re-implementation of RM Basic, a dialect of BASIC developed by Research Machines in 1985 for the [RM Nimbus PC186](https://en.wikipedia.org/wiki/RM_Nimbus), targeting 64-bit Windows and Linux operating systems.  It is written in Go and uses the [Nimbgobus](https://github.com/adamstimb/nimgobus) extension for [ebiten game engine](https://ebiten.org) to simulate the inputs and outputs of the RM Nimbus.

High-level goals:

- Backwards compatibility without the severe resource limitations of the original platform
- Recreate the original RM Basic user interface
- Extend the dialect to handle http requests


## Status

An alpha release is in the works featuring the original RM Basic user interface, just enough commands to write simple programs, and an installer for Windows.  See the [release notes](https://github.com/adamstimb/rmbasicx64/blob/main/docs/releaseNotes.md) for details.

## I want it now!

Ok, but don't say I didn't warn you.  Following the alpha release a Windows installer will be available.  In the meantime here's how to build and run the application on Windows:

### Prerequisites

- Latest version of [Go](https://golang.org/doc/install)
- [Git for Windows](https://gitforwindows.org/) (GITBASH)

### Build

Open GITBASH and clone this repository:

```bash
git clone https://github.com/adamstimb/rmbasicx64.git
```

Change directory and run the test and build scripts:

```bash
cd rmbasicx64/scripts
./test.sh
./build.sh
```

All being well a file called `rmbasic.exe` will appear in the `build\` folder.

Use File Explorer to make a new folder called `rmbasicx64` in `C:\Program Files` and move the `rmbasic.exe` file into it.  Double-click `rmbasic.exe` to run.

If you get a message saying "Windows protected your PC" click "More info" then "Run anyway".

## Screenshot

![editor](https://github.com/adamstimb/rmbasicx64/blob/main/docs/screenshots/editor-screenshot.png)

## Links

- [retrocomputingforum.com](https://retrocomputingforum.com/) - A retrocomputing forum with a [discussion thread](https://retrocomputingforum.com/t/rm-nimbus-basic-revival-64-bits/) on this project
- [Ebiten](https://ebiten.org/) - A dead simple 2D game library for Go
- [nimgobus](https://github.com/adamstimb/nimgobus) - An RM Nimbus-inspired Ebiten extension for building retro apps and games in Go
- [Crafting Interpreters](https://craftinginterpreters.com/) - The scanner code was inspired by the examples in this book
- [Writing an Interpreter in Go](https://interpreterbook.com/) - The parser code was inspired by the examples in this book
- [Facebook](https://www.facebook.com/RMNimbus/) - RM Nimbus facebook group
- [Center for Computing History](http://www.computinghistory.org.uk/) - original RM Nimbus manuals and technical data
- [Center for Computing History - RM Nimbus PC (Later Beige Model)](http://www.computinghistory.org.uk/det/41537/RM-Nimbus-PC-(Later-Beige-Model)/) - online exhibit
- [The Nimbus Museum](https://thenimbus.co.uk/) - online museum that looks like the Welcome Disk!
- [RM Nimbus](https://en.wikipedia.org/wiki/RM_Nimbus) - Wikipedia article
- [mame](https://www.mamedev.org/) - comprehensive retro computer emulation project
- [Nimbusinator](https://github.com/adamstimb/nimbusinator) - the Pythonic predecessor to Nimgobus
- [Ironstone Innovation](https://ironstoneinnovation.eu) - what I do for a living