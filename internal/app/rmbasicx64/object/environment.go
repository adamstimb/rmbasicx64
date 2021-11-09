package object

import (
	"fmt"
	"sort"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/ast"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
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
func (p *program) GetLineForEditing(lineNumber int) (string, bool) {
	if p.Jump(lineNumber, 0) {
		return p.lines[p.sortedIndex[p.curLineIndex+1]], true
	} else {
		return "", false
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
func (p *program) Renumber() {
	// TODO: Handle line numbers in GOTO statements
	// TODO: Handle optional params
	firstLineNumber := 10
	increment := 10
	p.Sort()
	p.Start()
	newLineNumber := firstLineNumber
	// populate new line map then replace the old line map
	newLines := make(map[int]string)
	for !p.EndOfProgram() {
		newLines[newLineNumber] = p.GetLine()
		newLineNumber += increment
		p.Next()
	}
	p.lines = make(map[int]string)
	p.lines = newLines
	p.Sort()
}

// Dump and Copy are used to transfer the program from one env to another
func (p *program) Dump() (sortedIndex []int, lines map[int]string) {
	sortedIndex = p.sortedIndex
	lines = p.lines
	return
}
func (p *program) Copy(sortedIndex []int, lines map[int]string) {
	p.lines = lines
	p.sortedIndex = sortedIndex
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
		return j.items[len(j.items)-1]
	}
	return nil
}
func (j *jumpStack) Push(item interface{}) {
	j.items = append(j.items, item)
}
func (j *jumpStack) Pop() interface{} {
	if len(j.items) > 0 {
		item := j.items[len(j.items)-1]
		j.items = j.items[:len(j.items)-1]
		return item
	}
	return nil
}

type storeKey struct {
	Scope  int
	Name   string
	Global bool
}

type Environment struct {
	store               map[storeKey]Object
	globals             []string
	scope               int
	Degrees             bool
	outer               *Environment
	Program             program
	JumpStack           jumpStack
	Prerun              bool
	dataItems           []Object
	subroutines         []*ast.SubroutineStatement
	functions           []*ast.FunctionDeclaration
	LeaveFunctionSignal bool
	ReturnVals          []Object
}

// Dump and Copy are used to transfer global data, including the program itself, from one env to another
func (e *Environment) Dump() (store map[storeKey]Object, globals []string, scope int, degrees bool, outer *Environment, program program, jumpStack jumpStack, prerun bool, dataItems []Object, subroutines []*ast.SubroutineStatement, functions []*ast.FunctionDeclaration) {
	store = e.store
	globals = e.globals
	scope = e.scope
	degrees = e.Degrees
	outer = e.outer
	program = e.Program
	dataItems = e.dataItems
	subroutines = e.subroutines
	functions = e.functions
	return
}
func (e *Environment) Copy(store map[storeKey]Object, globals []string, scope int, degrees bool, outer *Environment, program program, jumpStack jumpStack, prerun bool, dataItems []Object, subroutines []*ast.SubroutineStatement, functions []*ast.FunctionDeclaration) {
	e.store = store
	e.globals = globals
	e.scope = scope
	e.Degrees = degrees
	e.outer = outer
	e.Program = program
	e.JumpStack = jumpStack
	e.Prerun = prerun
	e.dataItems = dataItems
	e.subroutines = subroutines
	e.functions = functions
}
func (e *Environment) NewScope() {
	e.LeaveFunctionSignal = false
	e.scope++
}
func (e *Environment) LeaveFunction() {
	e.LeaveFunctionSignal = true
}
func (e *Environment) KillScope() {
	// Remove all vars from the store with current scope and then
	// go one scope towards global if not already in global.
	if e.scope > 0 {
		toDelete := []storeKey{}
		for k, _ := range e.store {
			if k.Scope == e.scope {
				toDelete = append(toDelete, k)
			}
		}
		for _, k := range toDelete {
			delete(e.store, k)
		}
		e.scope--
	}
}

func (e *Environment) Global(name string) bool {
	key := storeKey{Scope: e.scope, Name: name}
	if _, ok := e.store[key]; ok {
		// variable used as local error
		return false
	} else {
		// add variable to global scope (-1) and to list of globals
		currentScope := e.scope
		e.scope = -1
		e.Set(name, nil)
		e.scope = currentScope
		e.globals = append(e.globals, name)
		return true
	}
}

func (e *Environment) IsGlobal(name string) bool {
	for i := 0; i < len(e.globals); i++ {
		if name == e.globals[i] {
			return true
		}
	}
	return false
}

// TODO: Globalize arrays

func (e *Environment) PushData(obj Object) {
	e.dataItems = append(e.dataItems, obj)
}

func (e *Environment) PopData() Object {
	if len(e.dataItems) > 0 {
		returnObj := e.dataItems[0]
		e.dataItems = e.dataItems[1:]
		return returnObj
	} else {
		return nil
	}
}

func (e *Environment) DeleteData() {
	e.dataItems = []Object{}
}

func (e *Environment) PushSubroutine(sub *ast.SubroutineStatement) {
	e.subroutines = append(e.subroutines, sub)
}

func (e *Environment) GetSubroutine(name string) (*ast.SubroutineStatement, bool) {
	for _, sub := range e.subroutines {
		if sub.Name.Value == name {
			return sub, true
		}
	}
	return nil, false
}

func (e *Environment) DeleteSubroutines() {
	e.subroutines = []*ast.SubroutineStatement{}
}

func (e *Environment) PushFunction(fun *ast.FunctionDeclaration) {
	e.functions = append(e.functions, fun)
}

func (e *Environment) GetFunction(name string) (*ast.FunctionDeclaration, bool) {
	for _, fun := range e.functions {
		if fun.Name.Value == name {
			return fun, true
		}
	}
	return nil, false
}

func (e *Environment) DeleteFunctions() {
	e.functions = []*ast.FunctionDeclaration{}
}

func (e *Environment) NewArray(name string, subscripts []int) (Object, bool) {
	_, ok := e.store[storeKey{Name: name, Scope: e.scope}]
	if ok {
		return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.ArrayAlreadyDimensioned), ErrorTokenIndex: 0}, false
	}
	maxIndex := 1
	for _, subscript := range subscripts {
		maxIndex *= subscript
	}
	// initialize items according to type
	items := make([]Object, maxIndex)
	if name[len(name)-1:] != "$" {
		for i := 0; i < len(items); i++ {
			items[i] = &Numeric{Value: 0}
		}
	} else {
		for i := 0; i < len(items); i++ {
			items[i] = &String{Value: ""}
		}
	}
	e.store[storeKey{Name: name, Scope: e.scope}] = &Array{Items: items, Subscripts: subscripts}
	return e.store[storeKey{Name: name, Scope: e.scope}], true
}

func (e *Environment) GetArray(name string, subscripts []int) (Object, bool) {
	objArray, ok := e.store[storeKey{Name: name, Scope: e.scope}]
	if !ok && e.outer != nil {
		objArray, ok = e.outer.GetArray(name, subscripts)
	}
	if !ok {
		return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.FunctionArrayNotFound), ErrorTokenIndex: 0}, false
	}
	arr, ok := objArray.(*Array)
	if !ok {
		return arr, ok
	}
	// Validate subscripts
	if len(subscripts) != len(arr.Subscripts) {
		// Wrong number of subscripts error
		return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.WrongNumberOfSubscripts), ErrorTokenIndex: 0}, false
	}
	for i := 0; i < len(subscripts); i++ {
		if subscripts[i] >= arr.Subscripts[i] || subscripts[i] < 0 {
			// Subscript out of range error
			return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.ArraySubscriptIsWrong), ErrorTokenIndex: 0}, false
		}
	}
	// Resolve index
	index := subscripts[len(subscripts)-1]
	if len(arr.Subscripts) > 1 {
		for j := len(arr.Subscripts) - 2; j >= 0; j-- {
			index += subscripts[j] * arr.Subscripts[j]
		}
	}
	// Return item obj
	return arr.Items[index], true
}

func (e *Environment) SetArray(name string, subscripts []int, val Object) (Object, bool) {
	objArray, ok := e.store[storeKey{Name: name, Scope: e.scope}]
	if !ok && e.outer != nil {
		objArray, ok = e.outer.Get(name)
	}
	if !ok {
		return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.FunctionArrayNotFound), ErrorTokenIndex: 0}, false
	}
	// Don't allow string val to bind to numeric variable
	if val.Type() == STRING_OBJ && name[len(name)-1:] != "$" {
		return &Error{Message: "Numeric expression needed"}, false
	}
	// Don't allow numeric val to bind to string variable
	if val.Type() != STRING_OBJ && name[len(name)-1:] == "$" {
		return &Error{Message: "String expression needed"}, false
	}
	// If a float value is bound to an integer variable (name ends with %) it is rounded-down first (manual 3.7)
	if val.Type() == NUMERIC_OBJ && name[len(name)-1:] == "%" {
		val = &Numeric{Value: float64(int64(val.(*Numeric).Value))}
	}
	arr, _ := objArray.(*Array)
	// Validate subscripts
	if len(subscripts) != len(arr.Subscripts) {
		// Wrong number of subscripts error
		return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.WrongNumberOfSubscripts), ErrorTokenIndex: 0}, false
	}
	for i := 0; i < len(subscripts); i++ {
		if subscripts[i] > arr.Subscripts[i] || subscripts[i] < 0 {
			// Subscript out of range error
			return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.ArraySubscriptIsWrong), ErrorTokenIndex: 0}, false
		}
	}
	// Resolve index
	index := subscripts[len(subscripts)-1]
	if len(arr.Subscripts) > 1 {
		for j := len(arr.Subscripts) - 2; j >= 0; j-- {
			index += subscripts[j] * arr.Subscripts[j]
		}
	}
	// Set item obj
	arr.Items[index] = val
	e.store[storeKey{Name: name, Scope: e.scope}] = arr
	return arr, true
}

func (e *Environment) Get(name string) (Object, bool) {
	// Use current scope if local or global scope if global
	key := storeKey{Name: name, Scope: e.scope}
	if e.IsGlobal(name) {
		key.Scope = -1
	}
	obj, ok := e.store[key]
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
	// Use current scope if local or global scope if global
	key := storeKey{Name: name, Scope: e.scope}
	if e.IsGlobal(name) {
		key.Scope = -1
	}
	e.store[key] = val
	return val
}
func (e *Environment) Wipe() {
	e.Program.New()
	e.store = make(map[storeKey]Object)
	e.DeleteData()
	e.Degrees = true
	e.outer = nil
	e.scope = 0
}

func NewEnvironment() *Environment {
	s := make(map[storeKey]Object)
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
		dataItems: []Object{},
		scope:     0,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
