package object

import (
	"fmt"
	"sort"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/ast"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/lexer"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
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
	p.Indent()
}
func (p *program) EndOfProgram() bool {
	if p.curLineIndex >= len(p.lines) {
		// end of program
		return true
	} else {
		return false
	}
}
func (p *program) List(fromLinenumber, toLinenumber int, fromLineOnly bool) []string {
	if len(p.lines) == 0 {
		return nil
	}
	listing := []string{}

	for i := 0; i < len(p.lines); i++ {
		// list entire program
		if fromLinenumber == 0 && toLinenumber == 0 {
			listing = append(listing, fmt.Sprintf("%d %s", p.sortedIndex[i], p.lines[p.sortedIndex[i]]))
		}
		// use lower and upper boundary to get the lines
		if fromLinenumber != 0 && toLinenumber != 0 {
			if p.sortedIndex[i] >= fromLinenumber && p.sortedIndex[i] <= toLinenumber {
				listing = append(listing, fmt.Sprintf("%d %s", p.sortedIndex[i], p.lines[p.sortedIndex[i]]))
			}
		}
		// use lower boundary to get the lines, and get only first line if fromLineOnly
		if fromLinenumber != 0 && toLinenumber == 0 {
			if p.sortedIndex[i] >= fromLinenumber {
				listing = append(listing, fmt.Sprintf("%d %s", p.sortedIndex[i], p.lines[p.sortedIndex[i]]))
				if fromLineOnly {
					break
				}
			}
		}
		// use upper boundary to get the lines
		if fromLinenumber == 0 && toLinenumber != 0 {
			if p.sortedIndex[i] <= toLinenumber {
				listing = append(listing, fmt.Sprintf("%d %s", p.sortedIndex[i], p.lines[p.sortedIndex[i]]))
			}
		}
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

// Indent is used to tidy the code and make it easier to read
func (p *program) Indent() {
	newProg := make(map[int]string)
	indent := ""
	nextIndent := ""
	for i := 0; i < len(p.lines); i++ {
		line := p.lines[p.sortedIndex[i]]
		// new version:
		nextIndent = indent
		l := &lexer.Lexer{}
		tokens := l.Scan(line)
		for _, toke := range tokens {
			tokenType := toke.TokenType
			switch tokenType {
			case token.FOR, token.PROCEDURE, token.FUNCTION, token.REPEAT:
				nextIndent += "  "
			case token.NEXT, token.ENDPROC, token.ENDFUN, token.UNTIL:
				nextIndent = strings.TrimPrefix(nextIndent, "  ")
			}
		}
		indent = nextIndent
		newProg[p.sortedIndex[i]] = indent + line + "\n"

		//// old version:
		//line = strings.TrimSpace(line)
		//fields := strings.Fields(line)
		//var firstWord string
		//if len(fields) > 0 {
		//	firstWord = fields[0]
		//} else {
		//	continue
		//}
		//switch firstWord {
		//case "FOR", "PROCEDURE", "FUNCTION", "REPEAT":
		//	newProg[p.sortedIndex[i]] = indent + line + "\n"
		//	indent += "  "
		//case "NEXT", "ENDPROC", "ENDFUNC", "UNTIL":
		//	indent = strings.TrimPrefix(indent, "  ")
		//	newProg[p.sortedIndex[i]] = indent + line + "\n"
		//default:
		//	newProg[p.sortedIndex[i]] = indent + line + "\n"
		//}
	}
	p.lines = newProg
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
	GlobalEnv           *Environment
	scope               int
	Degrees             bool
	outer               *Environment
	Program             program
	JumpStack           jumpStack
	Prerun              bool
	dataItems           []Object
	subroutines         []*ast.SubroutineStatement
	functions           []*ast.FunctionDeclaration
	procedures          []*ast.ProcedureDeclaration
	LeaveFunctionSignal bool
	EndProgramSignal    bool
	ReturnVals          []Object
}

// Dump and Copy are used to transfer global data, including the program itself, from a parent env to a child env
func (e *Environment) Dump() (scope int, degrees bool, outer *Environment,
	program program, jumpStack jumpStack, prerun bool,
	dataItems []Object, subroutines []*ast.SubroutineStatement, functions []*ast.FunctionDeclaration,
	procedures []*ast.ProcedureDeclaration, leaveFunctionSignal bool, endProgramSignal bool,
	returnVals []Object) {
	scope = e.scope
	degrees = e.Degrees
	outer = e.outer
	program = e.Program
	dataItems = e.dataItems
	subroutines = e.subroutines
	functions = e.functions
	procedures = e.procedures
	return
}
func (e *Environment) Copy(scope int, degrees bool, outer *Environment,
	program program, jumpStack jumpStack, prerun bool,
	dataItems []Object, subroutines []*ast.SubroutineStatement, functions []*ast.FunctionDeclaration,
	procedures []*ast.ProcedureDeclaration, leaveFunctionSignal bool, endProgramSignal bool,
	returnVals []Object) {
	e.scope = scope
	e.Degrees = degrees
	e.outer = outer
	e.Program = program
	e.JumpStack = jumpStack
	e.Prerun = prerun
	e.dataItems = dataItems
	e.subroutines = subroutines
	e.functions = functions
	e.procedures = procedures
}
func (e *Environment) NewScope() {
	e.LeaveFunctionSignal = false
	e.scope++
}
func (e *Environment) IsBaseScope() bool {
	if e.scope == 0 {
		return true
	} else {
		return false
	}
}
func (e *Environment) EndProgram() {
	e.EndProgramSignal = true
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
		// variable already defined in this scope
		return false
	} else {
		// Register variable in globals list
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

func (e *Environment) DeleteStore() {
	e.store = make(map[storeKey]Object)
	e.globals = []string{}
	e.GlobalEnv.store = make(map[storeKey]Object)
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

func (e *Environment) PushProcedure(proc *ast.ProcedureDeclaration) {
	e.procedures = append(e.procedures, proc)
}

func (e *Environment) GetProcedure(name string) (*ast.ProcedureDeclaration, bool) {
	for _, proc := range e.procedures {
		if proc.Name.Value == name {
			return proc, true
		}
	}
	return nil, false
}

func (e *Environment) DeleteProcedures() {
	e.procedures = []*ast.ProcedureDeclaration{}
}

func calculateAddressFromArraySubscripts(bounds []int, subscripts []int) int {
	// Special case of 1D array
	if len(bounds) == 1 {
		return subscripts[0] //+ 1
	}
	// Use method for row-major calculation: https://www.geeksforgeeks.org/calculating-the-address-of-an-element-in-an-n-dimensional-array/
	internalSequence := []int{}
	subOffset := 0
	boundsOffset := 1
	for i := 0; i < len(subscripts)-1; i++ {
		if i == 0 {
			en := subscripts[i] - subOffset
			snp1 := bounds[i+1] - boundsOffset
			enp1 := subscripts[i+1] - subOffset
			internalSequence = append(internalSequence, (en*snp1)+(enp1))
		} else {
			prev := internalSequence[len(internalSequence)-1]
			snp1 := bounds[i+1] - boundsOffset
			enp1 := subscripts[i+1] - subOffset
			internalSequence = append(internalSequence, (prev*snp1)+enp1)
		}
	}
	w := len(bounds)
	addr := w * internalSequence[len(internalSequence)-1]
	return addr
}

func (e *Environment) NewArray(name string, subscripts []int) (Object, bool) {
	//_, ok := e.store[storeKey{Name: name, Scope: e.scope}]
	//if ok {
	//	return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.ArrayAlreadyDimensioned), ErrorTokenIndex: 0}, false
	//}
	key := storeKey{Scope: 0, Name: name}
	// Don't redimension existing arrays
	if e.IsGlobal(name) {
		if _, ok := e.GlobalEnv.store[key]; ok {
			return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.ArrayAlreadyDimensioned), ErrorTokenIndex: 0}, false
		}
	} else {
		if _, ok := e.store[key]; ok {
			return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.ArrayAlreadyDimensioned), ErrorTokenIndex: 0}, false
		}
	}
	// Go ahead and dimension the array
	maxIndex := calculateAddressFromArraySubscripts(subscripts, subscripts) + 1
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
	//e.store[storeKey{Name: name, Scope: e.scope}] = &Array{Items: items, Subscripts: subscripts}
	//return e.store[storeKey{Name: name, Scope: e.scope}], true
	var newArray Array
	if e.IsGlobal(name) {
		newArray = Array{Items: items, Subscripts: subscripts}
		e.GlobalEnv.store[key] = &newArray
	} else {
		newArray = Array{Items: items, Subscripts: subscripts}
		e.store[key] = &newArray
	}
	return &newArray, true
}

func (e *Environment) GetArray(name string, subscripts []int) (Object, bool) {
	//objArray, ok := e.store[storeKey{Name: name, Scope: e.scope}]
	//if !ok && e.outer != nil {
	//	objArray, ok = e.outer.GetArray(name, subscripts)
	//}
	key := storeKey{Scope: 0, Name: name}
	var arr *Array
	var ok bool
	if e.IsGlobal(name) {
		arr, ok = e.GlobalEnv.store[key].(*Array)
	} else {
		arr, ok = e.store[key].(*Array)
	}
	if !ok {
		return &Error{Message: syntaxerror.ErrorMessage(syntaxerror.FunctionArrayNotFound), ErrorTokenIndex: 0}, false
	}
	//arr, ok := objArray.(*Array)
	//if !ok {
	//	return arr, ok
	//}
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
	index := calculateAddressFromArraySubscripts(arr.Subscripts, subscripts)
	// Return item obj
	return arr.Items[index], true
}

func (e *Environment) SetArray(name string, subscripts []int, val Object) (Object, bool) {
	key := storeKey{Scope: 0, Name: name}
	var arr *Array
	var ok bool
	if e.IsGlobal(name) {
		arr, ok = e.GlobalEnv.store[key].(*Array)
	} else {
		arr, ok = e.store[key].(*Array)
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
	//arr, _ := objArray.(*Array)
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
	index := calculateAddressFromArraySubscripts(arr.Subscripts, subscripts)
	//log.Printf("Setting array %s[%v] with index %v to %v", name, subscripts, index, val)
	// Set item obj
	arr.Items[index] = val
	if e.IsGlobal(name) {
		e.GlobalEnv.store[key] = arr
	} else {
		e.store[key] = arr
	}
	//e.store[storeKey{Name: name, Scope: e.scope}] = arr
	return arr, true
}

func (e *Environment) Get(name string) (Object, bool) {

	// Use current scope if local or global scope if global
	key := storeKey{Name: name, Scope: 0}
	if e.IsGlobal(name) {
		obj, ok := e.GlobalEnv.store[key]
		return obj, ok
	} else {
		obj, ok := e.store[key]
		return obj, ok
	}
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
	key := storeKey{Name: name, Scope: 0}
	if e.IsGlobal(name) {
		e.GlobalEnv.store[key] = val
	} else {
		e.store[key] = val
	}
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

func NewEnvironment(GlobalEnv *Environment) *Environment {
	s := make(map[storeKey]Object)
	p := &program{}
	j := &jumpStack{}
	p.New()
	j.New()
	return &Environment{
		store:     s,
		GlobalEnv: GlobalEnv,
		Degrees:   true,
		outer:     nil,
		Program:   *p,
		JumpStack: *j,
		dataItems: []Object{},
		scope:     0,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment(nil)
	env.outer = outer
	return env
}
