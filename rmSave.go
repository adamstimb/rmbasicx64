package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// rmSave represents the SAVE command
func (i *Interpreter) rmSave() (ok bool) {
	i.tokenPointer++
	if i.EndOfTokens() {
		// No filename passed
		i.errorCode = StringExpressionNeeded
		i.message = errorMessage(StringExpressionNeeded)
		i.badTokenIndex = 1
		return false
	}
	// Get filename
	filename, ok := i.AcceptAnyString()
	if ok {
		// Don't accept wildcards
		if strings.Contains(filename, "*") {
			i.errorCode = ExactFilenameIsNeeded
			i.message = errorMessage(ExactFilenameIsNeeded)
			i.badTokenIndex = 1
			return false
		}
		// If it doesn't have .BAS extension then add it
		if !strings.HasSuffix(strings.ToUpper(filename), ".BAS") {
			filename += ".BAS"
		}
		// Don't accept a directory
		info, err := os.Stat(filename)
		if !os.IsNotExist(err) {
			if info.IsDir() {
				i.errorCode = FilenameIsADirectory
				i.message = errorMessage(FilenameIsADirectory)
				i.badTokenIndex = 1
				return false
			}
		}
	} else {
		i.badTokenIndex = 1
		return false
	}
	// Pass through if no program
	if len(i.program) == 0 {
		return true
	}
	// Save program to file
	file, err := os.Create(filename)
	if err != nil {
		i.errorCode = FileOperationFailure
		i.message = errorMessage(FileOperationFailure)
		i.badTokenIndex = 0
		return false
	}
	w := bufio.NewWriter(file)
	lineOrder := i.GetLineOrder()
	for _, lineNumber := range lineOrder {
		w.WriteString(fmt.Sprintf("%d %s\n", lineNumber, i.program[lineNumber]))
	}
	w.Flush()
	file.Close()
	return true
}
