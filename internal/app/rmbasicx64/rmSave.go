package rmbasicx64

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// rmSave represents the SAVE command
func (i *Interpreter) RmSave() (ok bool) {
	i.TokenPointer++
	if i.EndOfTokens() {
		// No filename passed
		i.ErrorCode = syntaxerror.StringExpressionNeeded
		i.BadTokenIndex = 1
		return false
	}
	// Get filename
	filename, ok := i.AcceptAnyString()
	if ok {
		// Don't accept wildcards
		if strings.Contains(filename, "*") {
			i.ErrorCode = syntaxerror.ExactFilenameIsNeeded
			i.BadTokenIndex = 1
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
				i.ErrorCode = syntaxerror.FilenameIsADirectory
				i.BadTokenIndex = 1
				return false
			}
		}
	} else {
		i.BadTokenIndex = 1
		return false
	}
	// Pass through if no program
	if len(i.Program) == 0 {
		return true
	}
	// Save program to file
	f, err := os.Create(filename)
	if err != nil {
		i.ErrorCode = syntaxerror.FileOperationFailure
		i.BadTokenIndex = 0
		return false
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	lineOrder := i.GetLineOrder()
	for _, lineNumber := range lineOrder {
		w.WriteString(fmt.Sprintf("%d %s\n", lineNumber, i.Program[lineNumber]))
	}
	w.Flush()
	return true
}
