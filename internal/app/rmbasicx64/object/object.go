package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/ast"
)

const (
	NUMERIC_OBJ      = "NUMERIC"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL" // RM Basic didn't have null...I think it just created a new var with zero or "" value...maybe?
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ERROR_OBJ        = "ERROR"
	WARNING_OBJ      = "WARNING"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("FUNCTION")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}
func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Numeric struct {
	Value float64
}

func (i *Numeric) Inspect() string {
	return fmt.Sprintf("%f", i.Value)
}
func (i *Numeric) Type() ObjectType {
	return NUMERIC_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}
func (s *String) Inspect() string {
	return s.Value
}

type Error struct {
	Message         string
	ErrorTokenIndex int
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Warning struct {
	Message string
}

func (w *Warning) Type() ObjectType {
	return ERROR_OBJ
}
func (w *Warning) Inspect() string {
	return "WARNING: " + w.Message
}
