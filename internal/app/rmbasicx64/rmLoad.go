package rmbasicx64

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// RmLoad represents the LOAD command
func (i *Interpreter) RmLoad() (ok bool) {
	i.TokenPointer++
	// Get required filename
	filename := ""
	val, ok := i.OnExpression("string")
	if ok {
		filename = val.(string)
		// Don't accept wildcards
		if strings.Contains(filename, "*") {
			i.ErrorCode = syntaxerror.ExactFilenameIsNeeded
			return false
		}
		// If it doesn't have .BAS extension then add it
		if !strings.HasSuffix(strings.ToUpper(filename), ".BAS") {
			filename += ".BAS"
		}
		// Don't accept a directory or nonexistant file
		info, err := os.Stat(filename)
		if !os.IsNotExist(err) {
			if info.IsDir() {
				i.ErrorCode = syntaxerror.FilenameIsADirectory
				return false
			}
		} else {
			i.ErrorCode = syntaxerror.UnableToOpenNamedFile
			return false
		}
	} else {
		return false
	}
	// Load program using ImmediateInput
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		i.ErrorCode = syntaxerror.FileOperationFailure
		i.TokenPointer = 1
		i.g.Print("Badness")
		return false
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		_ = i.ImmediateInput(line)
		if i.ErrorCode != syntaxerror.Success {
			return false
		}
	}
	return true
}
