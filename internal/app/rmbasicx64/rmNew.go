package rmbasicx64

// RmNew represents the NEW command
func (i *Interpreter) RmNew() (ok bool) {
	i.TokenPointer++
	if !i.OnSegmentEnd() {
		return false
	}
	// just initialize interpreter
	i.Init(i.g)
	return true
}
