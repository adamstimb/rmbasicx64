package evaluator

import (
	"github.com/adamstimb/rmbasicx64yar/internal/app/rmbasicx64yar/object"
)

// don't forget to add builtins to the map in lexer as well (better solution?)
var builtins = map[string]*object.Builtin{
	"LEN": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Numeric{
					Value: float64(len(arg.Value)),
				}
			default:
				return newError("argument to `LEN` not supported, got %s", args[0].Type())
			}
		},
	},
}
