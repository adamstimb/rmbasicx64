package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-sysinfo"
)

// bToGb converts bytes to Gb
func bToGb(b uint64) uint64 {
	return b / 1024 / 1024 / 1024
}

func repl(i *Interpreter) {

	// Until *all* the fundamentals of the language have been implemented we'll use
	// this simple text-based REPL as the UI.

	// Collect system info
	host, err := sysinfo.Host()
	if err != nil {
		panic("Could not detect system information")
	}
	memInfo, err := host.Memory()
	if err != nil {
		panic("Could not detect host memory information")
	}

	// Welcome screen and initialize interpreter
	fmt.Printf("\nRM NIMBUS\n\n")
	fmt.Printf("This is a tribute project and is in no way linked to or endorsed by RM plc.\n\n")
	fmt.Printf("RM BASICx64 Version 0.00 23rd March 2021\n")
	fmt.Printf("%dG bytes workspace available.\n", bToGb(memInfo.Available))
	i.Init()

	// REPL loop
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(":")
		code, _ := reader.ReadString('\n')
		code = strings.TrimSpace(code)
		ok := i.RunLine(code)
		if !ok {
			fmt.Printf("Syntax error: %s\n", i.message)
			fmt.Printf("%s\n", i.FormatCode(code, i.badTokenIndex, false))
		}
	}
}

func main() {
	repl(&Interpreter{})
}
