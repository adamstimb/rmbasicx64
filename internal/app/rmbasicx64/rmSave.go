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
	// Get filename
	val, ok := i.OnExpression("string")
	if !ok {
		return false
	}
	filename := val.(string)
	// No more params
	if !i.OnSegmentEnd() {
		return false
	}
	// Execute
	// Don't accept wildcards
	if strings.Contains(filename, "*") {
		i.ErrorCode = syntaxerror.ExactFilenameIsNeeded
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
			return false
		}
	}
	// Pass through if no program
	if len(i.Program) == 0 {
		return true
	}
	// Save program to file
	f, err := os.Create(filename)
	if err != nil {
		i.ErrorCode = syntaxerror.FileOperationFailure
		i.TokenPointer = 0
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
