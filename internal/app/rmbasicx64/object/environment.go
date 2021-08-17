package object

import (
	"fmt"
	"sort"
)

type program struct {
	lines                  map[int]string
	sortedIndex            []int
	curLineIndex           int
	JumpToStatement        int
	CurrentStatementNumber int
}

func (p *program) New() {
	p.lines = make(map[int]string)
	p.sortedIndex = []int{}
	p.curLineIndex = 0
	p.JumpToStatement = 0
	p.CurrentStatementNumber = 0
}
func (p *program) Sort() {
	keys := []int{}
	for k := range p.lines {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	p.sortedIndex = keys
}
func (p *program) Start() {
	p.curLineIndex = 0
	p.JumpToStatement = 0
	p.CurrentStatementNumber = 0
}
func (p *program) Next() {
	p.curLineIndex += 1
	p.JumpToStatement = 0
	p.CurrentStatementNumber = 0
}

// Jump is used to resume program execution from a specific line number and statement within
// that line.  Jump allows GOTO, GOSUB, FOR, WHILE and PROCEDURE/FUNCTION calls to be implemented.
func (p *program) Jump(lineNumber int, statementIndex int) bool {
	// Go to top of program and search for the required lineNumber and set JumpToStatement
	// if found
	currentLocation := p.curLineIndex
	p.Start()
	for !p.EndOfProgram() {
		if p.GetLineNumber() == lineNumber {
			// Found lineNumber, so back up the current location and we're done
			p.curLineIndex -= 1
			p.JumpToStatement = statementIndex
			return true
		} else {
			// Try next line
			p.Next()
		}
	}
	// Failed to find lineNumber.  Return to original location and
	// return false.
	p.curLineIndex = currentLocation
	return false
}
func (p *program) GetLineNumber() int {
	if len(p.lines) > 0 {
		return p.sortedIndex[p.curLineIndex]
	} else {
		return -1
	}
}
func (p *program) GetLine() string {
	if len(p.lines) > 0 {
		return p.lines[p.sortedIndex[p.curLineIndex]]
	} else {
		return ""
	}
}
func (p *program) AddLine(lineNumber int, line string) {
	if line == "" {
		// delete line if it exists
		delete(p.lines, lineNumber)
		p.Sort()
	} else {
		p.lines[lineNumber] = line
		p.Sort()
	}
}
func (p *program) EndOfProgram() bool {
	if p.curLineIndex >= len(p.lines) {
		// end of program
		return true
	} else {
		return false
	}
}
func (p *program) List() []string {
	if len(p.lines) == 0 {
		return nil
	}
	listing := []string{}
	for i := 0; i < len(p.lines); i++ {
		listing = append(listing, fmt.Sprintf("%d %s", p.sortedIndex[i], p.lines[p.sortedIndex[i]]))
	}
	return listing
}

// JumpStack is used to store all the return points and parameters for loops and function/procedure calls
type jumpStack struct {
	items []interface{}
}

func (j *jumpStack) New() {
	j.items = make([]interface{}, 0)
}
func (j *jumpStack) Peek() interface{} {

	if len(j.items) > 0 {
		return j.items[0]
	}
	return nil
}
func (j *jumpStack) Push(item interface{}) {
	j.items = append(j.items, item)
}
func (j *jumpStack) Pop() interface{} {
	if len(j.items) > 0 {
		item := j.items[0]
		j.items = j.items[:len(j.items)-1]
		return item
	}
	return nil
}

type Environment struct {
	store     map[string]Object
	Degrees   bool
	outer     *Environment
	Program   program
	JumpStack jumpStack
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	// Don't allow string val to bind to numeric variable
	if val.Type() == STRING_OBJ && name[len(name)-1:] != "$" {
		return &Error{Message: "Numeric expression needed"}
	}
	// Don't allow numeric val to bind to string variable
	if val.Type() != STRING_OBJ && name[len(name)-1:] == "$" {
		return &Error{Message: "String expression needed"}
	}
	// If a float value is bound to an integer variable (name ends with %) it is rounded-down first (manual 3.7)
	if val.Type() == NUMERIC_OBJ && name[len(name)-1:] == "%" {
		val = &Numeric{Value: float64(int64(val.(*Numeric).Value))}
	}
	e.store[name] = val
	return val
}
func (e *Environment) Wipe() {
	e.Program.New()
	e.store = make(map[string]Object)
	e.Degrees = true
	e.outer = nil
}
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	p := &program{}
	j := &jumpStack{}
	p.New()
	j.New()
	return &Environment{
		store:     s,
		Degrees:   true,
		outer:     nil,
		Program:   *p,
		JumpStack: *j,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
