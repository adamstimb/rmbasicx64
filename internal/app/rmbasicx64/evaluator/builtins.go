package evaluator

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/object"
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
	"ABS": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: math.Abs(arg.Value),
				}
			default:
				return newError("argument to `ABS` not supported, got %s", args[0].Type())
			}
		},
	},
	"ATN": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: math.Abs(arg.Value),
				}
			default:
				return newError("argument to `ATN` not supported, got %s", args[0].Type())
			}
		},
	},
	"COS": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: math.Cos(arg.Value),
				}
			default:
				return newError("argument to `COS` not supported, got %s", args[0].Type())
			}
		},
	},
	"SIN": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				log.Printf("SIN(%g)\n", arg.Value)
				return &object.Numeric{
					Value: math.Sin(arg.Value),
				}
			default:
				return newError("argument to `SIN` not supported, got %s", args[0].Type())
			}
		},
	},
	"EXP": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: math.Exp(arg.Value),
				}
			default:
				return newError("argument to `EXP` not supported, got %s", args[0].Type())
			}
		},
	},
	"INT": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: float64(int64(arg.Value)), // RM Basic truncated numbers down, rather than round them up
				}
			default:
				return newError("argument to `INT` not supported, got %s", args[0].Type())
			}
		},
	},
	"LN": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: math.Log(arg.Value),
				}
			default:
				return newError("argument to `LN` not supported, got %s", args[0].Type())
			}
		},
	},
	"LOG": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: math.Log10(arg.Value),
				}
			default:
				return newError("argument to `LOG` not supported, got %s", args[0].Type())
			}
		},
	},
	"RND": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				// In RM Basic, any negative number reseeds the random number generator
				retValue := float64(0)
				if arg.Value < 0 {
					rand.Seed(time.Now().UnixNano())
					retValue = 0
				} else {
					if arg.Value <= 1.0 {
						// generate random float between 0 and 1
						retValue = rand.Float64()
					} else {
						// generate random integer between 0 and arg.Value
						retValue = float64(rand.Intn(int(arg.Value)))
					}
				}

				return &object.Numeric{
					Value: retValue,
				}
			default:
				return newError("argument to `RND` not supported, got %s", args[0].Type())
			}
		},
	},
	"SGN": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				retValue := 0.0
				if arg.Value < 0 {
					retValue = -1.0
				}
				if arg.Value > 0 {
					retValue = 1.0
				}
				return &object.Numeric{
					Value: retValue,
				}
			default:
				return newError("argument to `SGN` not supported, got %s", args[0].Type())
			}
		},
	},
	"SQR": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: math.Sqrt(arg.Value),
				}
			default:
				return newError("argument to `SQR` not supported, got %s", args[0].Type())
			}
		},
	},
	"TAN": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.Numeric{
					Value: math.Tan(arg.Value),
				}
			default:
				return newError("argument to `TAN` not supported, got %s", args[0].Type())
			}
		},
	},
	//"STR$": &object.Builtin{
	//	Fn: func(args ...object.Object) object.Object {
	//		if len(args) != 1 {
	//			return newError("wrong number of arguments, got %d, want %d", len(args), 1)
	//		}
	//
	//		switch arg := args[0].(type) {
	//		case *object.Numeric:
	//			return &object.String{
	//				Value: fmt.Sprintf("%g", arg.Value),
	//			}
	//		default:
	//			return newError("argument to `STR$` not supported, got %s", args[0].Type())
	//		}
	//	},
	//},
}
