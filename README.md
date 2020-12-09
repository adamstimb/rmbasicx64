# RM BASICx64

A backwards-compatible RM Basic interpreter for 64-bit machines.

_RM BASICx64 is a tribute project and is in no way linked to or endorsed by RM plc._

## About

The plan is to re-implement the RM Basic interpeter and code editor (originally implemented on the [RM Nimbus PC186](https://en.wikipedia.org/wiki/RM_Nimbus) in the 1980s) for modern 64-bit machines.  RM BASICx64 is written in Go and uses the [Nimbgobus](https://github.com/adamstimb/nimgobus) extension for [ebiten game engine](https://ebiten.org) to simulate (as opposed to emulate) the inputs and outputs of the RM Nimbus.

High-level goals:

- Backwards compatibility - enable crusties like me to run their old RM Basic programs on a modern computer without emulation and the severe resource limitations of the original platform
- Get on the internet - extend the dialect to handle http requests
- Authentic 1980s user experience - recreate the original code editor, but also support easy-to-use modern editors such as VSCode

### Screenshot
# ![editor](editor.png)
