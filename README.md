# RM BASICx64

_RM BASICx64 is a tribute project and is in no way linked to or endorsed by RM plc._

## About

RM BASICx64 is a re-implementation of RM Basic, a dialect of BASIC developed by Research Machines in 1985 for the [RM Nimbus PC186](https://en.wikipedia.org/wiki/RM_Nimbus), targeting 64-bit Windows and Linux operating systems.  It is written in Go and uses the [Nimbgobus](https://github.com/adamstimb/nimgobus) extension for [ebiten game engine](https://ebiten.org) to simulate the inputs and outputs of the RM Nimbus.

## Status

A preliminary release is available featuring the original RM Basic user interface, expression evaluation and some simple text commands.  Click the "Watch" button on the top right to get notifications of new releases!

## I want it now!

Ok, but don't say I didn't warn you.  Here's how to build RM BASICx64 from source on Windows or Linux:

### Prerequisites

- Latest version of [Go](https://golang.org/doc/install)
- [Git for Windows](https://gitforwindows.org/) (GITBASH) (if using Windows)

### Build

Open (GIT)BASH and clone this repository:

```bash
git clone https://github.com/adamstimb/rmbasicx64.git
```

Change directory and run the build script for your operating system:

```bash
cd rmbasicx64
cd scripts
./build-linux.sh    # To build a Linux executable, or...
./build-windows.sh  # ... to build a Windows .exe
```

If you're running Linux, you can run the executable straight away:

```bash
../build/rmbasicx64
```

If you're running Windows a file called `rmbasicx64.exe` will appear in the `build\` folder.

Use File Explorer to make a new folder called `rmbasicx64` in `C:\Program Files` and move the `rmbasicx64.exe` file into it.  Double-click `rmbasicx64.exe` to run.

If you get a message saying "Windows protected your PC" click "More info" then "Run anyway".

## Screenshots

![RM BASICx64 running on Ubuntu](https://github.com/adamstimb/rmbasicx64/blob/main/docs/screenshots/interpreter-loaded.png)

RM BASICx64 running on Ubuntu

![Some lyrics from "Bike" by The Pink Floyd printed with string evaluation](https://github.com/adamstimb/rmbasicx64/blob/main/docs/screenshots/bike-lyrics.png)

Some lyrics from "Bike" by The Pink Floyd printed with string evaluation

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