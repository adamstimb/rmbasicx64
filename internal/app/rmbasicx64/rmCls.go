package rmbasicx64

// rmCls represents the CLS command
// TODO: parameters
func (i *Interpreter) RmCls() (ok bool) {
	// Ensure no parameters
	i.TokenPointer++
	if !i.OnSegmentEnd() {
		return false
	}
	// execute
	i.g.Cls()
	i.g.SetCurpos(1, 1)
	return true
}
