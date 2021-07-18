package object

type Environment struct {
	store map[string]Object
	outer *Environment
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

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{
		store: s,
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
