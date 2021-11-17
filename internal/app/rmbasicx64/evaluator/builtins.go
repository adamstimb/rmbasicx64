package evaluator

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/game"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/object"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// don't forget to add builtins to the map in lexer as well (better solution?)
var builtins = map[string]*object.Builtin{
	"LEN": &object.Builtin{
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
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
	"LOOKUP": &object.Builtin{
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			obj := args[0]
			var val string
			if stringVal, ok := obj.(*object.String); ok {
				val = stringVal.Value
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: 0}
			}
			// Don't allow * or ?
			if strings.Contains(val, "*") || strings.Contains(val, "?") {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ExactFilenameIsNeeded), ErrorTokenIndex: 0}
			}
			// Add .BAS if necessary
			if !strings.HasSuffix(strings.ToUpper(val), ".BAS") {
				val += ".BAS"
			}
			// execute
			systemPath := getAbsPath(val)
			// ensuring systemPath is a file and not a directory gets you a true, otherwise false
			fileInfo, err := os.Stat(systemPath)
			result := -1.0
			if err != nil {
				result = 0.0
			} else {
				if fileInfo.IsDir() {
					result = 0.0
				}
			}
			return &object.Numeric{Value: result}
		},
	},
	"PATH$": &object.Builtin{
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 0)
			}
			systemPath := getAbsPath("")
			nimbusPath := strings.ReplaceAll(systemPath[len(g.WorkspacePath):], "/", "\\")
			if nimbusPath == "" {
				nimbusPath = "\\"
			}
			return &object.String{Value: nimbusPath}
		},
	},
	"GET": &object.Builtin{
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
			if len(args) > 1 {
				return newError("wrong number of arguments, got %d, want < %d", len(args), 2)
			}
			// Determine timeout from args.  If no arg passed then timeout is infinite (-1).
			timeout := 100 * time.Millisecond
			useTimeout := false
			if len(args) > 0 {
				switch arg := args[0].(type) {
				case *object.Numeric:
					if arg.Value > 0 {
						timeout = time.Duration(100) * time.Millisecond * time.Duration(int(arg.Value))
						useTimeout = true
					}
					if arg.Value == 0 {
						timeout = time.Duration(0)
					}
				default:
					return newError("argument to `GET` not supported, got %s", args[0].Type())
				}
			}
			// Get
			elapsedTime := 1 * time.Millisecond
			for {
				c := g.Get()
				if c > 0 {
					// normal key pressed
					return &object.Numeric{
						Value: float64(c),
					}
				}
				// Handle zero timeout
				if timeout == time.Duration(0) {
					break
				}
				time.Sleep(100 * time.Millisecond)
				elapsedTime += 100 * time.Millisecond
				if useTimeout && elapsedTime > timeout {
					break
				}
			}
			// if key not pressed before timeout
			return &object.Numeric{
				Value: float64(0),
			}
		},
	},
	"STR$": &object.Builtin{
		Fn: func(env *object.Environment, g *game.Game, args []object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got %d, want %d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Numeric:
				return &object.String{
					Value: fmt.Sprintf("%g", arg.Value),
				}
			default:
				return newError("argument to `STR$` not supported, got %s", args[0].Type())
			}
		},
	},
}
