package evaluator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/ast"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/game"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/lexer"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/object"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/parser"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
	"github.com/adamstimb/rmbasicx64/pkg/nimgobus"
)

// Because null and boolean values never change we can reference them instead of
// creating new objects
var (
	NULL  = &object.Null{}
	TRUE  = &object.Numeric{Value: -1.0}
	FALSE = &object.Numeric{Value: 0}
)

func Eval(g *game.Game, node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.RemStatement:
		return nil
	case *ast.ByeStatement:
		os.Exit(0)
	case *ast.EndStatement:
		env.EndProgram()
		return nil
	case *ast.RunStatement:
		return evalRunStatement(g, node, env)
	case *ast.NewStatement:
		env.Wipe()
		return nil
	case *ast.SaveStatement:
		return evalSaveStatement(g, node, env)
	case *ast.LoadStatement:
		return evalLoadStatement(g, node, env)
	case *ast.FetchStatement:
		return evalFetchStatement(g, node, env)
	case *ast.WriteblockStatement:
		return evalWriteblockStatement(g, node, env)
	case *ast.ReadblockStatement:
		return evalReadblockStatement(g, node, env)
	case *ast.CopyblockStatement:
		return evalCopyblockStatement(g, node, env)
	case *ast.SquashStatement:
		return evalSquashStatement(g, node, env)
	case *ast.ClearblockStatement:
		return evalClearblockStatement(g, node, env)
	case *ast.DelblockStatement:
		return evalDelblockStatement(g, node, env)
	case *ast.KeepStatement:
		return evalKeepStatement(g, node, env)
	case *ast.ListStatement:
		return evalListStatement(g, node, env)
	case *ast.ClsStatement:
		return evalClsStatement(g, node, env)
	case *ast.HomeStatement:
		return evalHomeStatement(g, node, env)
	case *ast.DirStatement:
		return evalDirStatement(g, node, env)
	case *ast.ChdirStatement:
		return evalChdirStatement(g, node, env)
	case *ast.MkdirStatement:
		return evalMkdirStatement(g, node, env)
	case *ast.RmdirStatement:
		return evalRmdirStatement(g, node, env)
	case *ast.EraseStatement:
		return evalEraseStatement(g, node, env)
	case *ast.RenameStatement:
		return evalRenameStatement(g, node, env)
	case *ast.SetMouseStatement:
		return evalSetMouseStatement(g, node, env)
	case *ast.SetModeStatement:
		return evalSetModeStatement(g, node, env)
	case *ast.SetPaperStatement:
		return evalSetPaperStatement(g, node, env)
	case *ast.SetBorderStatement:
		return evalSetBorderStatement(g, node, env)
	case *ast.SetPenStatement:
		return evalSetPenStatement(g, node, env)
	case *ast.SetDegStatement:
		return evalSetDegStatement(g, node, env)
	case *ast.SetColourStatement:
		return evalSetColourStatement(g, node, env)
	case *ast.SetRadStatement:
		return evalSetRadStatement(g, node, env)
	case *ast.SetCurposStatement:
		return evalSetCurposStatement(g, node, env)
	case *ast.SetWritingStatement:
		return evalSetWritingStatement(g, node, env)
	case *ast.SetPatternStatement:
		return evalSetPatternStatement(g, node, env)
	case *ast.SetConfigBootStatement:
		return evalSetConfigBootStatement(g, node, env)
	case *ast.MoveStatement:
		return evalMoveStatement(g, node, env)
	case *ast.PrintStatement:
		return evalPrintStatement(g, node, env)
	case *ast.PlotStatement:
		return evalPlotStatement(g, node, env)
	case *ast.LineStatement:
		return evalLineStatement(g, node, env)
	case *ast.AreaStatement:
		return evalAreaStatement(g, node, env)
	case *ast.SetFillStyleStatement:
		return evalSetFillStyleStatement(g, node, env)
	case *ast.CircleStatement:
		return evalCircleStatement(g, node, env)
	case *ast.PointsStatement:
		return evalPointsStatement(g, node, env)
	case *ast.FloodStatement:
		return evalFloodStatement(g, node, env)
	case *ast.GotoStatement:
		return evalGotoStatement(g, node, env)
	case *ast.EditStatement:
		return evalEditStatement(g, node, env)
	case *ast.DataStatement:
		return evalDataStatement(g, node, env)
	case *ast.SubroutineStatement:
		return evalSubroutineStatement(g, node, env)
	case *ast.GosubStatement:
		return evalGosubStatement(g, node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(g, node, env)
	case *ast.FunctionDeclaration:
		return evalFunctionDeclaration(g, node, env)
	case *ast.ProcedureDeclaration:
		return evalProcedureDeclaration(g, node, env)
	case *ast.ProcedureCallStatement:
		return evalProcedureCallStatement(g, node, env)
	case *ast.ResultStatement:
		return evalResultStatement(g, node, env)
	case *ast.EndfunStatement:
		return evalEndfunStatement(g, node, env)
	case *ast.EndprocStatement:
		return evalEndprocStatement(g, node, env)
	case *ast.LeaveStatement:
		return evalLeaveStatement(g, node, env)
	case *ast.ReadStatement:
		return evalReadStatement(g, node, env)
	case *ast.RestoreStatement:
		return evalRestoreStatement(g, node, env)
	case *ast.RenumberStatement:
		return evalRenumberStatement(g, node, env)
	case *ast.RepeatStatement:
		return evalRepeatStatement(g, node, env)
	case *ast.UntilStatement:
		return evalUntilStatement(g, node, env)
	case *ast.ForStatement:
		return evalForStatement(g, node, env)
	case *ast.NextStatement:
		return evalNextStatement(g, node, env)
	case *ast.DimStatement:
		return evalDimStatement(g, node, env)
	case *ast.AskMouseStatement:
		return evalAskMouseStatement(g, node, env)
	case *ast.AskBlocksizeStatement:
		return evalAskBlocksizeStatement(g, node, env)
	case *ast.Program:
		return evalProgram(g, node, env)
	case *ast.ExpressionStatement:
		return Eval(g, node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(g, node, env)
	case *ast.IfStatement:
		return evalIfStatement(g, node, env)
	case *ast.LetStatement:
		val := Eval(g, node.Value, env)
		if isError(val) {
			return val
		}
		if len(node.Name.Subscripts) > 0 {
			// is Array
			subscripts := make([]int, len(node.Name.Subscripts))
			for i := 0; i < len(node.Name.Subscripts); i++ {
				obj := Eval(g, node.Name.Subscripts[i], env)
				if val, ok := obj.(*object.Numeric); ok {
					subscripts[i] = int(val.Value)
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: node.Token.Index}
				}
			}
			ret, _ := env.SetArray(node.Name.Value, subscripts, val)
			return ret
		} else {
			// is variable
			return env.Set(node.Name.Value, val)
		}
	case *ast.BindStatement:
		val := Eval(g, node.Value, env)
		if isError(val) {
			return val
		}
		if len(node.Name.Subscripts) > 0 {
			// is Array
			subscripts := make([]int, len(node.Name.Subscripts))
			for i := 0; i < len(node.Name.Subscripts); i++ {
				obj := Eval(g, node.Name.Subscripts[i], env)
				if val, ok := obj.(*object.Numeric); ok {
					subscripts[i] = int(val.Value)
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: node.Token.Index}
				}
			}
			ret, _ := env.SetArray(node.Name.Value, subscripts, val)
			return ret
		} else {
			// is variable
			return env.Set(node.Name.Value, val)
		}

	// Expressions
	case *ast.NumericLiteral:
		return &object.Numeric{
			Value: node.Value,
		}
	case *ast.StringLiteral:
		return &object.String{
			Value: node.Value,
		}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(g, node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(g, node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(g, node.Right, env)
		if isError(right) {
			return right
		}
		result := evalInfixExpression(node.Operator, left, right)
		return result
	case *ast.Identifier:
		obj := evalIdentifier(g, node, env)
		// Catch builtin
		if obj.Type() == object.BUILTIN_OBJ {
			args := evalExpressions(g, node.Subscripts, env)
			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}
			return applyFunction(env, g, obj, args)
		}
		// If a warning is returned, print the warning *then* re-run the evaluation and return
		if warningMsg, ok := obj.(*object.Warning); ok {
			g.Print(fmt.Sprintf("Warning: %s", warningMsg.Message))
			g.Put(13)
			return evalIdentifier(g, node, env)
		} else {
			return obj
		}
	case *ast.CallExpression:
		// This is most likely redundant now
		function := Eval(g, node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(g, node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(env, g, function, args)
	}
	return nil
}

func applyFunction(env *object.Environment, g *game.Game, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(g, fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		// If the builtin is a trig function and env.Degrees is true we need to
		// convert the passed angle from degrees to radians,
		if fn == builtins["ATN"] || fn == builtins["COS"] || fn == builtins["SIN"] || fn == builtins["TAN"] {
			if env.Degrees {
				args[0].(*object.Numeric).Value *= (math.Pi / 180)
			}
			return fn.Fn(env, g, args) // args...
		} else {
			return fn.Fn(env, g, args) // args...
		}
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalPrintStatement(g *game.Game, stmt *ast.PrintStatement, env *object.Environment) object.Object {
	printStr := ""
	oldTextBoxSlot, _, _, _, _ := g.AskWriting()
	tempTextBoxSlot := oldTextBoxSlot
	_, curY := g.AskCurpos()
	// Evaluate and handle TextBoxSlot if set
	if stmt.TextBoxSlot != nil {
		obj := Eval(g, stmt.TextBoxSlot, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			tempTextBoxSlot = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	g.SetWriting(tempTextBoxSlot)

	for _, val := range stmt.PrintList {
		// Handle seperator type
		if s, ok := val.(string); ok {
			switch s {
			case "noSpace":
				printStr += ""
			case "nextPrintZone":
				printStr += "     " // TODO: Implement actual print zones in Nimgobus
			case "newLine":
				g.Print(printStr)
				g.Put(13)
				printStr = ""
			}
			continue
		}
		obj := Eval(g, val.(ast.Node), env)
		if isError(obj) {
			if oldTextBoxSlot != tempTextBoxSlot {
				g.SetWriting(oldTextBoxSlot)
				g.SetCurpos(1, curY)
			}
			return obj
		}
		if numericVal, ok := obj.(*object.Numeric); ok {
			printStr += fmt.Sprintf("%g", numericVal.Value)
		}
		if boolVal, ok := obj.(*object.Boolean); ok {
			if boolVal.Value {
				printStr += "TRUE"
			} else {
				printStr += "FALSE"
			}
		}
		if stringVal, ok := obj.(*object.String); ok {
			printStr += stringVal.Value
		}
	}
	g.Print(printStr)
	g.Put(13)
	if oldTextBoxSlot != tempTextBoxSlot {
		g.SetWriting(oldTextBoxSlot)
		g.SetCurpos(1, curY)
	}
	return nil
}

func evalPlotStatement(g *game.Game, stmt *ast.PlotStatement, env *object.Environment) object.Object {
	// Handle defaults
	var Brush, Direction, Font, Over, SizeX, SizeY int
	var Text string
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if g.ValidateColour(int(val.Value)) {
				Brush = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Direction == nil {
		Direction = -255
	} else {
		obj := Eval(g, stmt.Direction, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			Direction = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Font == nil {
		Font = -255
	} else {
		obj := Eval(g, stmt.Font, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			Font = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Over == nil {
		Over = -255
	} else {
		obj := Eval(g, stmt.Over, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			Over = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.SizeX == nil {
		SizeX = -255
	} else {
		obj := Eval(g, stmt.SizeX, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if int(val.Value) > 0 {
				SizeX = int(val.Value)
			} else {
				SizeX = 1
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.SizeY == nil {
		SizeY = -255
	} else {
		obj := Eval(g, stmt.SizeY, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if int(val.Value) > 0 {
				SizeY = int(val.Value)
			} else {
				SizeY = 1
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Handle text string
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		Text = fmt.Sprintf("%g", val.Value)
	}
	if val, ok := obj.(*object.Boolean); ok {
		if val.Value {
			Text = "TRUE"
		} else {
			Text = "FALSE"
		}
	}
	if val, ok := obj.(*object.String); ok {
		Text = val.Value
	}
	// Handle coord list
	var coordList []nimgobus.XyCoord
	var X, Y int
	for i := 0; i < len(stmt.CoordList)-1; i += 2 {
		obj := Eval(g, stmt.CoordList[i], env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			X = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		obj = Eval(g, stmt.CoordList[i+1], env)
		if val, ok := obj.(*object.Numeric); ok {
			Y = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		coordList = append(coordList, nimgobus.XyCoord{X, Y})
	}
	// Execute
	for _, coord := range coordList {
		opt := nimgobus.PlotOptions{Brush: Brush, Direction: Direction, Font: Font, SizeX: SizeX, SizeY: SizeY, Over: Over}
		g.Plot(opt, Text, coord.X, coord.Y)
	}
	return nil
}

func evalLineStatement(g *game.Game, stmt *ast.LineStatement, env *object.Environment) object.Object {
	// Handle defaults
	var Brush, Over int
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if g.ValidateColour(int(val.Value)) {
				Brush = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Over == nil {
		Over = -255
	} else {
		obj := Eval(g, stmt.Over, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			Over = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Handle coord list
	var coordList []nimgobus.XyCoord
	var X, Y int
	for i := 0; i < len(stmt.CoordList)-1; i += 2 {
		obj := Eval(g, stmt.CoordList[i], env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			X = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		obj = Eval(g, stmt.CoordList[i+1], env)
		if val, ok := obj.(*object.Numeric); ok {
			Y = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		coordList = append(coordList, nimgobus.XyCoord{X, Y})
	}
	// Execute
	opt := nimgobus.LineOptions{Brush: Brush, Over: Over}
	g.Line(opt, coordList)
	return nil
}

func evalCircleStatement(g *game.Game, stmt *ast.CircleStatement, env *object.Environment) object.Object {
	// Handle defaults
	var Brush, Over int
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if val.Value >= 0 {
				Brush = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Over == nil {
		Over = -255
	} else {
		obj := Eval(g, stmt.Over, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			Over = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Handle fill style
	var fillStyle nimgobus.FillStyle
	if stmt.FillStyle == nil {
		fillStyle.Style = -1
	} else {
		obj := Eval(g, stmt.FillStyle, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if val.Value >= 0 && val.Value <= 2 {
				fillStyle.Style = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		if fillStyle.Style == 2 {
			obj := Eval(g, stmt.FillHatching, env)
			if isError(obj) {
				return obj
			}
			if val, ok := obj.(*object.Numeric); ok {
				if val.Value >= 0 && val.Value <= 5 {
					fillStyle.Hatching = int(val.Value)
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
				}
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
			}
			if stmt.FillColour2 == nil {
				fillStyle.Colour2 = -1
			} else {
				obj := Eval(g, stmt.FillColour2, env)
				if isError(obj) {
					return obj
				}
				if val, ok := obj.(*object.Numeric); ok {
					if val.Value >= 0 && val.Value <= 5 {
						fillStyle.Colour2 = int(val.Value)
					} else {
						return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
					}
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
				}
			}
		}
	}
	// Handle radius
	var radius int
	obj := Eval(g, stmt.Radius, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		radius = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Handle coord list
	var coordList []nimgobus.XyCoord
	var X, Y int
	for i := 0; i < len(stmt.CoordList)-1; i += 2 {
		obj := Eval(g, stmt.CoordList[i], env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			X = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		obj = Eval(g, stmt.CoordList[i+1], env)
		if val, ok := obj.(*object.Numeric); ok {
			Y = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		coordList = append(coordList, nimgobus.XyCoord{X, Y})
	}
	// Execute
	opt := nimgobus.CircleOptions{Brush: Brush, Over: Over, FillStyle: fillStyle}
	for _, coord := range coordList {
		g.Circle(opt, radius, coord.X, coord.Y)
	}
	return nil
}

func evalPointsStatement(g *game.Game, stmt *ast.PointsStatement, env *object.Environment) object.Object {
	// Handle defaults
	var Brush, Over, Style int
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if g.ValidateColour(int(val.Value)) {
				Brush = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Style == nil {
		Style = -255
	} else {
		obj := Eval(g, stmt.Style, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if g.ValidateStyle(int(val.Value)) {
				Style = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Over == nil {
		Over = -255
	} else {
		obj := Eval(g, stmt.Over, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			Over = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Handle coord list
	var coordList []nimgobus.XyCoord
	var X, Y int
	for i := 0; i < len(stmt.CoordList)-1; i += 2 {
		obj := Eval(g, stmt.CoordList[i], env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			X = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		obj = Eval(g, stmt.CoordList[i+1], env)
		if val, ok := obj.(*object.Numeric); ok {
			Y = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		coordList = append(coordList, nimgobus.XyCoord{X, Y})
	}
	// Execute
	opt := nimgobus.PointsOptions{Brush: Brush, Over: Over, Style: Style}
	g.Points(opt, coordList)
	return nil
}

func evalFloodStatement(g *game.Game, stmt *ast.FloodStatement, env *object.Environment) object.Object {
	// Handle defaults
	var Brush, EdgeColour int
	var UseEdgeColour bool
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if g.ValidateColour(int(val.Value)) {
				Brush = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.UseEdgeColour == nil {
		UseEdgeColour = false
	} else {
		obj := Eval(g, stmt.UseEdgeColour, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			UseEdgeColour = isTruthy(val)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.EdgeColour == nil {
		EdgeColour = -255
	} else {
		obj := Eval(g, stmt.EdgeColour, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if g.ValidateColour(int(val.Value)) {
				EdgeColour = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Handle fill style
	var fillStyle nimgobus.FillStyle
	if stmt.FillStyle == nil {
		fillStyle.Style = -1
	} else {
		obj := Eval(g, stmt.FillStyle, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if val.Value >= 0 && val.Value <= 2 {
				fillStyle.Style = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		if fillStyle.Style == 2 {
			obj := Eval(g, stmt.FillHatching, env)
			if isError(obj) {
				return obj
			}
			if val, ok := obj.(*object.Numeric); ok {
				if val.Value >= 0 && val.Value <= 5 {
					fillStyle.Hatching = int(val.Value)
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
				}
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
			}
			if stmt.FillColour2 == nil {
				fillStyle.Colour2 = -1
			} else {
				obj := Eval(g, stmt.FillColour2, env)
				if isError(obj) {
					return obj
				}
				if val, ok := obj.(*object.Numeric); ok {
					if val.Value >= 0 && val.Value <= 5 {
						fillStyle.Colour2 = int(val.Value)
					} else {
						return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
					}
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
				}
			}
		}
	}
	// Handle coord list
	var coordList []nimgobus.XyCoord
	var X, Y int
	for i := 0; i < len(stmt.CoordList)-1; i += 2 {
		obj := Eval(g, stmt.CoordList[i], env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			X = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		obj = Eval(g, stmt.CoordList[i+1], env)
		if val, ok := obj.(*object.Numeric); ok {
			Y = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		coordList = append(coordList, nimgobus.XyCoord{X, Y})
	}
	// Execute
	opt := nimgobus.FloodOptions{Brush: Brush, UseEdgeColour: UseEdgeColour, EdgeColour: EdgeColour, FillStyle: fillStyle}
	for _, coord := range coordList {
		g.Flood(opt, coord)
	}
	return nil
}

func evalAreaStatement(g *game.Game, stmt *ast.AreaStatement, env *object.Environment) object.Object {
	// Handle defaults
	var Brush, Over int
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if g.ValidateColour(int(val.Value)) {
				Brush = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Over == nil {
		Over = -255
	} else {
		obj := Eval(g, stmt.Over, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			Over = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Handle fill style
	var fillStyle nimgobus.FillStyle
	if stmt.FillStyle == nil {
		fillStyle.Style = -1
	} else {
		obj := Eval(g, stmt.FillStyle, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if val.Value >= 0 && val.Value <= 2 {
				fillStyle.Style = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		if fillStyle.Style == 2 {
			obj := Eval(g, stmt.FillHatching, env)
			if isError(obj) {
				return obj
			}
			if val, ok := obj.(*object.Numeric); ok {
				if val.Value >= 0 && val.Value <= 5 {
					fillStyle.Hatching = int(val.Value)
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
				}
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
			}
			if stmt.FillColour2 == nil {
				fillStyle.Colour2 = -1
			} else {
				obj := Eval(g, stmt.FillColour2, env)
				if isError(obj) {
					return obj
				}
				if val, ok := obj.(*object.Numeric); ok {
					if val.Value >= 0 && val.Value <= 5 {
						fillStyle.Colour2 = int(val.Value)
					} else {
						return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
					}
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
				}
			}
		}
	}
	// Handle coord list
	var coordList []nimgobus.XyCoord
	var X, Y int
	for i := 0; i < len(stmt.CoordList)-1; i += 2 {
		obj := Eval(g, stmt.CoordList[i], env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			X = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		obj = Eval(g, stmt.CoordList[i+1], env)
		if val, ok := obj.(*object.Numeric); ok {
			Y = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		coordList = append(coordList, nimgobus.XyCoord{X, Y})
	}
	// Execute
	opt := nimgobus.AreaOptions{Brush: Brush, Over: Over, FillStyle: fillStyle}
	g.Area(opt, coordList)
	return nil
}

func evalSetFillStyleStatement(g *game.Game, stmt *ast.SetFillStyleStatement, env *object.Environment) object.Object {
	// Handle fill style
	var fillStyle nimgobus.FillStyle
	obj := Eval(g, stmt.FillStyle, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		if val.Value >= 0 && val.Value <= 2 {
			fillStyle.Style = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if fillStyle.Style == 2 {
		obj := Eval(g, stmt.FillHatching, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if val.Value >= 0 && val.Value <= 5 {
				fillStyle.Hatching = int(val.Value)
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		if stmt.FillColour2 == nil {
			fillStyle.Colour2 = -1
		} else {
			obj := Eval(g, stmt.FillColour2, env)
			if isError(obj) {
				return obj
			}
			if val, ok := obj.(*object.Numeric); ok {
				if val.Value >= 0 && val.Value <= 5 {
					fillStyle.Colour2 = int(val.Value)
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
				}
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		}
	}
	// Execute
	g.SetFillStyle(fillStyle.Style, fillStyle.Hatching, fillStyle.Colour2)
	return nil
}

func evalSaveStatement(g *game.Game, stmt *ast.SaveStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	filename := ""
	if stringVal, ok := obj.(*object.String); ok {
		filename = stringVal.Value
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Don't allow * or ?
	if strings.Contains(filename, "*") || strings.Contains(filename, "?") {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ExactFilenameIsNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Add .BAS if necessary
	if !strings.HasSuffix(strings.ToUpper(filename), ".BAS") {
		filename += ".BAS"
	}
	// Preprend workspace folder
	fullpath := filepath.Join(g.WorkspacePath, filename)
	// Don't allow directories
	if isDirectory(fullpath) {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FilenameIsADirectory), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Warn if file already exists
	_, err := ioutil.ReadFile(fullpath)
	if !errors.Is(err, os.ErrNotExist) {
		// Warn user and ask to abort
		g.Print("Named file already exists")
		var abort bool
		var gotAnswer bool
		for !gotAnswer {
			g.Put(13)
			g.Print("Abort command? (Y/N): ")
			waiting := true
			for waiting {
				key := g.Get()
				if key < 0 {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				g.Put(key)
				switch key {
				case 89: // Y
					abort = true
					gotAnswer = true
					waiting = false
				case 121: // y
					abort = true
					gotAnswer = true
					waiting = false
				case 78: // N
					abort = false
					gotAnswer = true
					waiting = false
				case 110: // n
					abort = false
					gotAnswer = true
					waiting = false
				default:
					waiting = false
				}
			}
		}
		g.Put(13)
		if abort {
			return nil
		}
	}
	// Save the program
	file, err := os.Create(fullpath)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure)}
	}
	defer file.Close()
	for _, lineString := range env.Program.List(0, 0, false) {
		file.WriteString(fmt.Sprintf("%s\n", lineString))
	}
	return obj
}

// isDirectory determines if a file represented by `path` is a directory or not
// https://freshman.tech/snippets/go/check-if-file-is-dir/
func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func evalLoadStatement(g *game.Game, stmt *ast.LoadStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	filename := ""
	if stringVal, ok := obj.(*object.String); ok {
		filename = stringVal.Value
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Don't allow * or ?
	if strings.Contains(filename, "*") || strings.Contains(filename, "?") {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ExactFilenameIsNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Add .BAS if necessary
	if !strings.HasSuffix(strings.ToUpper(filename), ".BAS") {
		filename += ".BAS"
	}
	// Preprend workspace folder
	fullpath := filepath.Join(g.WorkspacePath, filename)
	// Don't allow directories
	if isDirectory(fullpath) {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FilenameIsADirectory), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Load the program
	fileBytes, err := ioutil.ReadFile(fullpath)
	// Handle file doesn't exist
	if errors.Is(err, os.ErrNotExist) {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnableToOpenNamedFile), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Handle any other errors
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Committed to load the program so erase any existing program in memory
	env.Program.New()
	// To read into the program space we just pretend the code is being manually keyed it (I think that's how it worked originally)
	sliceData := strings.Split(string(fileBytes), "\n")
	l := &lexer.Lexer{}
	for _, rawLine := range sliceData {
		if g.BreakInterruptDetected {
			break
		}
		l.Scan(rawLine)
		p := parser.New(l, g)
		line := p.ParseLine()
		// Check of parser errors here.  Parser errors are handled just like evaluation errors but
		// obviously we'll skip evaluation if parsing already failed.
		if errorMsg, hasError := p.GetError(); hasError {
			g.Print(errorMsg)
			g.Put(13)
			p.JumpToToken(0)
			g.Print(p.PrettyPrint())
			g.Put(13)
			continue
		}
		// And this is temporary while we're still migrating from Monkey to RM Basic
		if len(p.Errors()) > 0 {
			g.Print("Oops! Some random parsing error occurred. These will be handled properly downstream by for now here's some spewage:")
			g.Put(13)
			p.JumpToToken(0)
			g.Print(p.PrettyPrint())
			g.Put(13)
			for _, msg := range p.Errors() {
				g.Print(msg)
				g.Put(13)
			}
			continue
		}
		// Add new line to stored program
		if line.Statements == nil {
			env.Program.AddLine(line.LineNumber, line.LineString)
			continue
		}
		// Execute each statement in the inputted line.  If an error occurs, print the
		// error message and stop.
		for _, stmt := range line.Statements {
			obj := Eval(g, stmt, env)
			if errorMsg, ok := obj.(*object.Error); ok {
				g.Print(errorMsg.Message)
				g.Put(13)
				p.JumpToToken(0)
				g.Print(p.PrettyPrint())
				g.Put(13)
				break
			}
		}
	}
	return obj
}

func evalFetchStatement(g *game.Game, stmt *ast.FetchStatement, env *object.Environment) object.Object {
	// Block
	var block int
	obj := Eval(g, stmt.Block, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		block = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if block < 0 || block > 99 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Path
	var path string
	obj = Eval(g, stmt.Path, env)
	if val, ok := obj.(*object.String); ok {
		path = val.Value
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Don't allow * or ?
	if strings.Contains(path, "*") || strings.Contains(path, "?") {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ExactFilenameIsNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Preprend workspace folder
	fullpath := filepath.Join(g.WorkspacePath, path)
	// Don't allow directories
	if isDirectory(fullpath) {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FilenameIsADirectory), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Infer file format from extention and fail if not recognised
	if !(strings.HasSuffix(strings.ToUpper(path), ".PNG") || strings.HasSuffix(strings.ToUpper(path), ".JPG") || strings.HasSuffix(strings.ToUpper(path), ".JPEG")) {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnsupportedImageFileFormat), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Handle file doesn't exist
	_, err := ioutil.ReadFile(fullpath)
	if errors.Is(err, os.ErrNotExist) {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnableToOpenNamedFile), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Execute
	ok := g.Fetch(block, fullpath)
	// Return any other errors
	if !ok {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.CouldNotDecodeImageFile), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	return obj
}

func evalWriteblockStatement(g *game.Game, stmt *ast.WriteblockStatement, env *object.Environment) object.Object {
	// Block
	var block int
	obj := Eval(g, stmt.Block, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		block = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if block < 0 || block > 99 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// X
	var x int
	obj = Eval(g, stmt.X, env)
	if val, ok := obj.(*object.Numeric); ok {
		x = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Y
	var y int
	obj = Eval(g, stmt.Y, env)
	if val, ok := obj.(*object.Numeric); ok {
		y = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Over
	var over bool
	obj = Eval(g, stmt.Over, env)
	if val, ok := obj.(*object.Numeric); ok {
		over = isTruthy(val)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Execute
	g.Writeblock(block, x, y, over)
	return nil
}

func evalReadblockStatement(g *game.Game, stmt *ast.ReadblockStatement, env *object.Environment) object.Object {
	// Block
	var block int
	obj := Eval(g, stmt.Block, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		block = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if block < 0 || block > 99 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// X1
	var x1 int
	obj = Eval(g, stmt.X1, env)
	if val, ok := obj.(*object.Numeric); ok {
		x1 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Y1
	var y1 int
	obj = Eval(g, stmt.Y1, env)
	if val, ok := obj.(*object.Numeric); ok {
		y1 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// X2
	var x2 int
	obj = Eval(g, stmt.X2, env)
	if val, ok := obj.(*object.Numeric); ok {
		x2 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Y2
	var y2 int
	obj = Eval(g, stmt.Y2, env)
	if val, ok := obj.(*object.Numeric); ok {
		y2 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Execute
	g.Readblock(block, x1, y1, x2, y2)
	return nil
}

func evalCopyblockStatement(g *game.Game, stmt *ast.CopyblockStatement, env *object.Environment) object.Object {
	// X1
	var x1 int
	obj := Eval(g, stmt.X1, env)
	if val, ok := obj.(*object.Numeric); ok {
		x1 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Y1
	var y1 int
	obj = Eval(g, stmt.Y1, env)
	if val, ok := obj.(*object.Numeric); ok {
		y1 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// X2
	var x2 int
	obj = Eval(g, stmt.X2, env)
	if val, ok := obj.(*object.Numeric); ok {
		x2 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Y2
	var y2 int
	obj = Eval(g, stmt.Y2, env)
	if val, ok := obj.(*object.Numeric); ok {
		y2 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Dx
	var dx int
	obj = Eval(g, stmt.Dx, env)
	if val, ok := obj.(*object.Numeric); ok {
		dx = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Dy
	var dy int
	obj = Eval(g, stmt.Dy, env)
	if val, ok := obj.(*object.Numeric); ok {
		dy = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Over
	var over bool
	obj = Eval(g, stmt.Over, env)
	if val, ok := obj.(*object.Numeric); ok {
		over = isTruthy(val)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Execute
	g.Readblock(0, x1, y1, x2, y2)
	g.Writeblock(0, dx, dy, over)
	g.Delblock(0)
	return nil
}

func evalSquashStatement(g *game.Game, stmt *ast.SquashStatement, env *object.Environment) object.Object {
	// Block
	var block int
	obj := Eval(g, stmt.Block, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		block = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if block < 0 || block > 99 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// X
	var x int
	obj = Eval(g, stmt.X, env)
	if val, ok := obj.(*object.Numeric); ok {
		x = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Y
	var y int
	obj = Eval(g, stmt.Y, env)
	if val, ok := obj.(*object.Numeric); ok {
		y = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Over
	var over bool
	obj = Eval(g, stmt.Over, env)
	if val, ok := obj.(*object.Numeric); ok {
		over = isTruthy(val)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Execute
	g.Squash(block, x, y, over)
	return nil
}

func evalDelblockStatement(g *game.Game, stmt *ast.DelblockStatement, env *object.Environment) object.Object {
	// Block
	var block int
	obj := Eval(g, stmt.Block, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		block = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if block < 0 || block > 99 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Execute
	g.Delblock(block)
	return nil
}

func evalKeepStatement(g *game.Game, stmt *ast.KeepStatement, env *object.Environment) object.Object {
	// Block
	var block int
	obj := Eval(g, stmt.Block, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		block = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if block < 0 || block > 99 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Path
	var path string
	obj = Eval(g, stmt.Path, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.String); ok {
		path = val.Value
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Don't allow * or ?
	if strings.Contains(path, "*") || strings.Contains(path, "?") {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ExactFilenameIsNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Preprend workspace folder
	fullpath := filepath.Join(g.WorkspacePath, path)
	// Don't allow directories
	if isDirectory(fullpath) {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FilenameIsADirectory), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Infer file format from extention and fail if not recognised
	if !(strings.HasSuffix(strings.ToUpper(path), ".PNG") || strings.HasSuffix(strings.ToUpper(path), ".JPG") || strings.HasSuffix(strings.ToUpper(path), ".JPEG")) {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnsupportedImageFileFormat), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	var format string
	if strings.HasSuffix(strings.ToUpper(path), ".PNG") {
		format = "png"
	}
	if strings.HasSuffix(strings.ToUpper(path), ".JPG") || strings.HasSuffix(strings.ToUpper(path), ".JPEG") {
		format = "jpeg"
	}
	// Execute
	err := g.Keep(block, format, fullpath)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	return nil
}

func evalSetModeStatement(g *game.Game, stmt *ast.SetModeStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		g.SetMode(int(val.Value))
		return obj
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
}

func evalSetPaperStatement(g *game.Game, stmt *ast.SetPaperStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		if g.ValidateColour(int(val.Value)) {
			g.SetPaper(int(val.Value))
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	return nil
}

func evalSetBorderStatement(g *game.Game, stmt *ast.SetBorderStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		if g.ValidateColour(int(val.Value)) {
			g.SetBorder(int(val.Value))
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	return nil
}

func evalSetPenStatement(g *game.Game, stmt *ast.SetPenStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		if g.ValidateColour(int(val.Value)) {
			g.SetPen(int(val.Value))
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	return nil
}

func evalSetCurposStatement(g *game.Game, stmt *ast.SetCurposStatement, env *object.Environment) object.Object {
	colVal := 0
	rowVal := 0
	col := Eval(g, stmt.Col, env)
	row := Eval(g, stmt.Row, env)
	// evaluate col
	if _, ok := col.(*object.Error); ok {
		return col
	}
	if val, ok := col.(*object.Numeric); ok {
		colVal = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}

	// evaluate row
	if _, ok := row.(*object.Error); ok {
		return col
	}
	if val, ok := row.(*object.Numeric); ok {
		rowVal = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	g.SetCurpos(colVal, rowVal)
	return nil
}

func evalSetWritingStatement(g *game.Game, stmt *ast.SetWritingStatement, env *object.Environment) object.Object {
	var slot, col1, row1, col2, row2 int
	obj := Eval(g, stmt.Slot, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		slot = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// activating slot 0 is allowed but you can't change it
	minSlot := 1
	if stmt.Col1 == nil {
		minSlot = 0
	}
	if slot < minSlot || slot > 10 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// SET WRITING e1
	if stmt.Col1 == nil {
		g.SetWriting(slot)
		return nil
	}
	// SET WRITING e1 TO e2, e3; e4, e5
	obj = Eval(g, stmt.Col1, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		col1 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if col1 < 1 || col1 > g.AskMode() {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.Row1, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		row1 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if row1 < 1 || row2 > 25 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.Col2, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		col2 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if col2 < 1 || col2 > g.AskMode() {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.Row2, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		row2 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if row2 < 1 || row2 > 25 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	g.SetWriting(slot, col1, row1, col2, row2)
	return nil
}

func evalSetPatternStatement(g *game.Game, stmt *ast.SetPatternStatement, env *object.Environment) object.Object {
	var slot, row, c1, c2, c3, c4 int
	obj := Eval(g, stmt.Slot, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		slot = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.Row, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		row = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.C1, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		c1 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.C2, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		c2 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.C3, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		c3 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.C4, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		c4 = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	g.SetPattern(slot, row, c1, c2, c3, c4)
	return nil
}

func evalMoveStatement(g *game.Game, stmt *ast.MoveStatement, env *object.Environment) object.Object {
	colsIncr := 0
	rowsIncr := 0
	cols := Eval(g, stmt.Cols, env)
	rows := Eval(g, stmt.Rows, env)
	// evaluate col
	if _, ok := cols.(*object.Error); ok {
		return cols
	}
	if val, ok := cols.(*object.Numeric); ok {
		colsIncr = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// evaluate row
	if _, ok := rows.(*object.Error); ok {
		return cols
	}
	if val, ok := rows.(*object.Numeric); ok {
		rowsIncr = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	col, row := g.AskCurpos()
	g.SetCurpos(col+colsIncr, row+rowsIncr)
	return nil
}

func evalSetDegStatement(g *game.Game, stmt *ast.SetDegStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	if _, ok := obj.(*object.Numeric); ok {
		env.Degrees = isTruthy(obj)
		return obj
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
}

func evalSetColourStatement(g *game.Game, stmt *ast.SetColourStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.PaletteSlot, env)
	if isError(obj) {
		return obj
	}
	var paletteSlot, basicColour, flashSpeed, flashColour int
	// paletteSlot
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		if g.ValidateColour(int(val.Value)) {
			paletteSlot = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// basicColour
	obj = Eval(g, stmt.BasicColour, env)
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		basicColour = int(val.Value)
		if basicColour < 0 || basicColour > 15 {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 3}
	}
	// optionals
	if stmt.FlashSpeed == nil || stmt.FlashColour == nil {
		flashSpeed = 0
		flashColour = 0
	} else {
		// flashSpeed
		obj = Eval(g, stmt.FlashSpeed, env)
		if _, ok := obj.(*object.Error); ok {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			flashSpeed = int(val.Value)
			if flashSpeed < 0 || flashSpeed > 2 {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 5} // <-- this is wrong...
		}
		// flashColour
		obj = Eval(g, stmt.FlashColour, env)
		if _, ok := obj.(*object.Error); ok {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			if flashColour < 0 || flashColour > 15 {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumberNotAllowedInRange), ErrorTokenIndex: stmt.Token.Index + 1}
			}
			flashColour = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 7} // <-- so's this!
		}
	}
	// execute
	g.SetColour(paletteSlot, basicColour, flashSpeed, flashColour)
	return nil
}

func evalSetConfigBootStatement(g *game.Game, stmt *ast.SetConfigBootStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if _, ok := obj.(*object.Numeric); ok {
		c, err := g.ReadConf()
		if err != nil {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure), ErrorTokenIndex: 0}
		}
		c.Boot = isTruthy(obj)
		if !g.WriteConf(c) {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure), ErrorTokenIndex: 0}
		}
		return obj
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
}

func evalGotoStatement(g *game.Game, stmt *ast.GotoStatement, env *object.Environment) object.Object {
	// Get line number direct from literal
	val, _ := strconv.ParseFloat(stmt.Linenumber.Literal, 64)
	lineNumber := int(val)
	if lineNumber < 0 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.PositiveValueRequired), ErrorTokenIndex: stmt.Linenumber.Index}
	}
	if env.Program.Jump(lineNumber, 0) {
		return nil
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.LineNumberDoesNotExist), ErrorTokenIndex: stmt.Linenumber.Index}
	}
}

func evalEditStatement(g *game.Game, stmt *ast.EditStatement, env *object.Environment) object.Object {
	// TODO: Handle no line number so try to get line of last error
	if stmt.Linenumber.Literal == "" {
		return nil
	}
	// Get line number direct from literal
	val, _ := strconv.ParseFloat(stmt.Linenumber.Literal, 64)
	lineNumber := int(val)
	if lineNumber < 0 {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.PositiveValueRequired), ErrorTokenIndex: stmt.Linenumber.Index}
	}
	// Edit the line if it exists
	if line, ok := env.Program.GetLineForEditing(lineNumber); ok {
		g.Print(fmt.Sprintf("%d ", lineNumber))
		rawLine := g.Input(line)
		if rawLine == "" {
			// cancel edit
			return nil
		} else {
			// tokenize and save changes
			l := &lexer.Lexer{}
			l.Scan(fmt.Sprintf("%d %s", lineNumber, rawLine))
			p := parser.New(l, g)
			line := p.ParseLine()
			env.Program.AddLine(lineNumber, line.LineString)
		}
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.LineNumberDoesNotExist), ErrorTokenIndex: stmt.Linenumber.Index}
	}
	return nil
}

func evalRestoreStatement(g *game.Game, stmt *ast.RestoreStatement, env *object.Environment) object.Object {
	resumeLine := env.Program.GetLineNumber()
	resumeStatement := env.Program.CurrentStatementNumber
	l := &lexer.Lexer{}
	env.Program.Start()
	// Jump to line number if specified and it exists in program
	if stmt.Linenumber.Literal != "" {
		// Get line number direct from literal
		val, _ := strconv.ParseFloat(stmt.Linenumber.Literal, 64)
		lineNumber := int(val)
		if lineNumber < 0 {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.PositiveValueRequired), ErrorTokenIndex: stmt.Linenumber.Index}
		}
		// Try to jump
		if !env.Program.Jump(lineNumber, 0) {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.LineNumberDoesNotExist), ErrorTokenIndex: stmt.Linenumber.Index}
		}
	}
	// Run through the stored program but only collect data
	env.DeleteData()
	for !env.Program.EndOfProgram() {
		l.Scan(env.Program.GetLine())
		p := parser.New(l, g)
		line := p.ParseLine()
		// Disregard parser errors as these will be handling during execution.
		if _, hasError := p.GetError(); hasError {
			continue
		}
		// Only DATA
		for statementNumber, stmt := range line.Statements {
			env.Program.CurrentStatementNumber = statementNumber
			tokenType := stmt.TokenLiteral()
			if tokenType == token.DATA {
				env.Prerun = true
				obj := Eval(g, stmt, env)
				env.Prerun = false
				if errorMsg, ok := obj.(*object.Error); ok {
					if errorMsg.ErrorTokenIndex != 0 {
						p.ErrorTokenIndex = errorMsg.ErrorTokenIndex
					}
					lineNumber := env.Program.GetLineNumber()
					g.Print(fmt.Sprintf("%s in line %d", errorMsg.Message, lineNumber))
					g.Put(13)
					p.JumpToToken(0)
					g.Print(fmt.Sprintf("%d %s", lineNumber, p.PrettyPrint()))
					g.Put(13)
					return nil
				}
			}
		}
		env.Program.Next()
	}
	env.Program.Jump(resumeLine, resumeStatement)
	env.Program.Next()
	return nil
}

func evalDataStatement(g *game.Game, stmt *ast.DataStatement, env *object.Environment) object.Object {
	// Pass if not in prerun
	if !env.Prerun {
		return nil
	}
	// Add items to data
	for _, item := range stmt.ItemList {
		switch item.TokenType {
		case token.NumericLiteral:
			val, _ := strconv.ParseFloat(item.Literal, 64)
			env.PushData(&object.Numeric{Value: val})
		case token.StringLiteral:
			env.PushData(&object.String{Value: item.Literal})
		case token.IdentifierLiteral:
			env.PushData(&object.String{Value: item.Literal})
		}
	}
	return nil
}

func evalSubroutineStatement(g *game.Game, stmt *ast.SubroutineStatement, env *object.Environment) object.Object {
	// Error if not prerun
	if !env.Prerun {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.CannotExecuteDefinition), ErrorTokenIndex: stmt.Token.Index}
	}
	// Register subroutine
	stmt.LineNumber = env.Program.GetLineNumber()
	stmt.StatementNumber = env.Program.CurrentStatementNumber
	env.PushSubroutine(stmt)
	return nil
}

func evalGosubStatement(g *game.Game, stmt *ast.GosubStatement, env *object.Environment) object.Object {
	// Push gosub statement onto jump stack so RETURN will know where to jump back to
	stmt.LineNumber = env.Program.GetLineNumber()
	stmt.StatementNumber = env.Program.CurrentStatementNumber
	env.JumpStack.Push(stmt)

	// Try to get location of subroutine and jump.  Yield error if not found (run on emulator
	// to find out which one)
	if sub, ok := env.GetSubroutine(stmt.Name.Value); ok {
		env.Program.Jump(sub.LineNumber, sub.StatementNumber)
		env.Program.Next()
	} else {
		// return whichever error
	}
	return nil
}

func evalReturnStatement(g *game.Game, stmt *ast.ReturnStatement, env *object.Environment) object.Object {
	// Pop return stack until we find a gosub statement.  If we don't find one, return the
	// RETURN without any GOSUB error.
	looking := true
	for looking {
		jumpItem := env.JumpStack.Pop()
		// test if this is a gosub and jump back to it if so
		if gosub, ok := jumpItem.(*ast.GosubStatement); ok {
			env.Program.Jump(gosub.LineNumber, gosub.StatementNumber)
			env.Program.Next()
			return nil
		}
	}
	return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ReturnWithoutAnyGosub), ErrorTokenIndex: stmt.Token.Index}
}

func evalFunctionDeclaration(g *game.Game, stmt *ast.FunctionDeclaration, env *object.Environment) object.Object {
	// Error if not prerun
	if !env.Prerun {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.CannotExecuteDefinition), ErrorTokenIndex: stmt.Token.Index}
	}
	// Register function
	stmt.LineNumber = env.Program.GetLineNumber()
	stmt.StatementNumber = env.Program.CurrentStatementNumber
	env.PushFunction(stmt)
	return nil
}

func evalResultStatement(g *game.Game, stmt *ast.ResultStatement, env *object.Environment) object.Object {
	if env.IsBaseScope() {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FunctionExitWithoutCall), ErrorTokenIndex: stmt.Token.Index}
	}
	obj := Eval(g, stmt.ResultValue, env)
	env.ReturnVals = append(env.ReturnVals, obj)
	env.LeaveFunction()
	return nil // Return vals are picked out of the env so this stays as nil
}

func evalEndfunStatement(g *game.Game, stmt *ast.EndfunStatement, env *object.Environment) object.Object {
	return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NeedResultToExitFunction), ErrorTokenIndex: stmt.Token.Index}
}

func evalProcedureDeclaration(g *game.Game, stmt *ast.ProcedureDeclaration, env *object.Environment) object.Object {
	// Error if not prerun
	if !env.Prerun {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.CannotExecuteDefinition), ErrorTokenIndex: stmt.Token.Index}
	}
	// Register function
	stmt.LineNumber = env.Program.GetLineNumber()
	stmt.StatementNumber = env.Program.CurrentStatementNumber
	env.PushProcedure(stmt)
	return nil
}

func evalEndprocStatement(g *game.Game, stmt *ast.EndprocStatement, env *object.Environment) object.Object {
	if env.IsBaseScope() {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ProcedureExitWithoutCall), ErrorTokenIndex: stmt.Token.Index}
	}
	// TODO: Handle return values
	env.ReturnVals = append(env.ReturnVals, nil)
	env.LeaveFunction()
	return nil
}

func evalLeaveStatement(g *game.Game, stmt *ast.LeaveStatement, env *object.Environment) object.Object {
	if env.IsBaseScope() {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ProcedureExitWithoutCall), ErrorTokenIndex: stmt.Token.Index}
	}
	// TODO: Handle return values
	env.ReturnVals = append(env.ReturnVals, nil)
	env.LeaveFunction()
	return nil
}

func evalRenumberStatement(g *game.Game, stmt *ast.RenumberStatement, env *object.Environment) object.Object {
	env.Program.Renumber()
	return nil
}

func evalRepeatStatement(g *game.Game, stmt *ast.RepeatStatement, env *object.Environment) object.Object {
	stmt.LineNumber = env.Program.GetLineNumber()
	stmt.StatementNumber = env.Program.CurrentStatementNumber
	// don't repush the same Repeat
	if repeatStmt, ok := env.JumpStack.Peek().(*ast.RepeatStatement); ok {
		if repeatStmt.LineNumber == stmt.LineNumber && repeatStmt.StatementNumber == stmt.StatementNumber {
			return nil
		}
	}
	env.JumpStack.Push(stmt)
	return nil
}

func evalArraySubscripts(g *game.Game, env *object.Environment, subscripts []ast.Expression) (evaluatedSubscripts []int, obj object.Object, ok bool) {
	for i := 0; i < len(subscripts); i++ {
		obj := Eval(g, subscripts[i], env)
		if val, ok := obj.(*object.Numeric); ok {
			evaluatedSubscripts[i] = int(val.Value)
		} else {
			return evaluatedSubscripts, &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: 1}, false
		}
	}
	return evaluatedSubscripts, nil, true
}

func evalReadStatement(g *game.Game, stmt *ast.ReadStatement, env *object.Environment) object.Object {
	for i := 0; i < len(stmt.VariableList); i++ {
		varName := stmt.VariableList[i].Value
		// evaluate array subscripts, if any
		var subscripts []int
		if subs, obj, ok := evalArraySubscripts(g, env, stmt.VariableList[i].Subscripts); ok {
			// all good
			subscripts = subs
		} else {
			return obj
		}
		// Variable type
		if varName[len(varName)-1] == '%' || varName[len(varName)-1] != '$' {
			// is numeric var
			obj := env.PopData()
			if obj == nil {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NoMoreDataToBeRead), ErrorTokenIndex: stmt.Token.Index}
			}
			// can only accept numeric data
			if val, ok := obj.(*object.Numeric); ok {
				if len(subscripts) > 0 {
					// is array
					if obj, ok := env.SetArray(varName, subscripts, val); ok {
						// all good
					} else {
						// failed
						return obj
					}
				} else {
					// is var
					env.Set(varName, &object.Numeric{Value: val.Value})
				}
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringVariableExpected), ErrorTokenIndex: stmt.Token.Index}
			}
		} else {
			// is string var - can accept numeric or string data
			obj := env.PopData()
			if obj == nil {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NoMoreDataToBeRead), ErrorTokenIndex: stmt.Token.Index}
			}
			// string data can be set directly
			if val, ok := obj.(*object.String); ok {
				if len(subscripts) > 0 {
					// is array
					if obj, ok := env.SetArray(varName, subscripts, val); ok {
						// all good
					} else {
						// failed
						return obj
					}
				} else {
					// is var
					env.Set(varName, &object.String{Value: val.Value})
				}
			}
			if val, ok := obj.(*object.Numeric); ok {
				// convert to string then set
				strVal := &object.String{Value: fmt.Sprintf("%g", val.Value)}
				if len(subscripts) > 0 {
					// is array
					if obj, ok := env.SetArray(varName, subscripts, strVal); ok {
						// all good
					} else {
						// failed
						return obj
					}
				} else {
					// is var
					env.Set(varName, strVal)
				}
			}
		}
	}
	return nil
}

func evalForStatement(g *game.Game, stmt *ast.ForStatement, env *object.Environment) object.Object {
	stmt.LineNumber = env.Program.GetLineNumber()
	stmt.StatementNumber = env.Program.CurrentStatementNumber
	// don't repush the same Repeat
	if forStmt, ok := env.JumpStack.Peek().(*ast.ForStatement); ok {
		if forStmt.LineNumber == stmt.LineNumber && forStmt.StatementNumber == stmt.StatementNumber {
			return nil
		}
	}
	// Otherwise evaluate start, stop, step expressions, bind counting variable, and push to stack
	var start, stop, step float64
	// Start
	obj := Eval(g, stmt.Start, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		start = val.Value
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Stop
	obj = Eval(g, stmt.Stop, env)
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		stop = val.Value
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Step (Default=1)
	if stmt.Step == nil {
		step = 1.0
	} else {
		obj := Eval(g, stmt.Step, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			step = val.Value
			if step < 0 {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.PositiveValueRequired), ErrorTokenIndex: stmt.Token.Index + 1}
			}
			if step == 0 {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StepValueNotLargeEnough), ErrorTokenIndex: stmt.Token.Index + 1}
			}
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Flip step sign if counting down
	if stop < start {
		step *= -1.0
	}
	// Bind counting variable (must be numeric) --- TODO this should be handled by parser!
	if stmt.Name.Value[len(stmt.Name.Value)-1:] == "$" {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericVariableNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	env.Set(stmt.Name.Value, &object.Numeric{Value: start})
	// Update ast with evaluated stop and step values, then push to stack
	stmt.StartValue = start
	stmt.StopValue = stop
	stmt.StepValue = step
	env.JumpStack.Push(stmt)
	return nil
}

func evalNextStatement(g *game.Game, stmt *ast.NextStatement, env *object.Environment) object.Object {
	// Ensure we're inside the FOR loop before evaluating condition
	if forStmt, ok := env.JumpStack.Peek().(*ast.ForStatement); ok {
		// Get value of counter variable
		var counterVal float64
		if obj, ok := env.Get(stmt.Name.Value); ok {
			if val, ok := obj.(*object.Numeric); ok {
				counterVal = val.Value
			} else {
				return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index}
			}
		} else {
			// TODO: Set warning, somehow
			env.Set(stmt.Name.Value, &object.Numeric{Value: 0})
			counterVal = 0
		}
		// Evaluate condition
		conditionMet := false
		if forStmt.StartValue < forStmt.StopValue && counterVal+forStmt.StepValue > forStmt.StopValue {
			conditionMet = true
		}
		if forStmt.StartValue > forStmt.StopValue && counterVal+forStmt.StepValue < forStmt.StopValue {
			conditionMet = true
		}
		if !conditionMet {
			// increment counter and loop again
			counterVal += forStmt.StepValue
			env.Set(stmt.Name.Value, &object.Numeric{Value: counterVal})
			env.Program.Jump(forStmt.LineNumber, forStmt.StatementNumber)
			return nil
		} else {
			// drop through loop
			env.JumpStack.Pop()
			return nil
		}
	}
	return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UntilWithoutAnyRepeat), ErrorTokenIndex: stmt.Token.Index}
}

func evalDimStatement(g *game.Game, stmt *ast.DimStatement, env *object.Environment) object.Object {
	subscripts := make([]int, len(stmt.Subscripts))
	for i := 0; i < len(stmt.Subscripts); i++ {
		obj := Eval(g, stmt.Subscripts[i], env)
		if val, ok := obj.(*object.Numeric); ok {
			subscripts[i] = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index}
		}
	}
	obj, _ := env.NewArray(stmt.Name.Value, subscripts)
	return obj
}

func evalAskMouseStatement(g *game.Game, stmt *ast.AskMouseStatement, env *object.Environment) object.Object {
	// Handle no args
	if stmt.XName == nil && stmt.YName == nil {
		return nil
	}
	// Set X, Y
	env.Set(stmt.XName.Value, &object.Numeric{Value: float64(g.MouseX)})
	env.Set(stmt.YName.Value, &object.Numeric{Value: float64(g.MouseY)})
	// Handle button if required
	if stmt.BName != nil {
		env.Set(stmt.BName.Value, &object.Numeric{Value: float64(g.MouseButton)})
	}
	return nil
}

func evalAskBlocksizeStatement(g *game.Game, stmt *ast.AskBlocksizeStatement, env *object.Environment) object.Object {
	// Handle no args
	if stmt.Block == nil {
		return nil
	}
	// Get block
	var block int
	obj := Eval(g, stmt.Block, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		block = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Execute
	width, height, mode := g.AskBlocksize(block)
	// Set variables
	if stmt.Width != nil {
		env.Set(stmt.Width.Value, &object.Numeric{Value: float64(width)})
	}
	if stmt.Height != nil {
		env.Set(stmt.Height.Value, &object.Numeric{Value: float64(height)})
	}
	if stmt.Mode != nil {
		env.Set(stmt.Mode.Value, &object.Numeric{Value: float64(mode)})
	}
	return nil
}

func evalUntilStatement(g *game.Game, stmt *ast.UntilStatement, env *object.Environment) object.Object {
	// Ensure we're inside a Repeat loop before evaluating condition
	if repeatStmt, ok := env.JumpStack.Peek().(*ast.RepeatStatement); ok {
		condition := Eval(g, stmt.Condition, env)
		if isError(condition) {
			return condition
		}
		if isTruthy(condition) {
			// condition is true to drop through loop
			env.JumpStack.Pop()
			return nil
		} else {
			// condition is false so jump to repeat
			env.Program.Jump(repeatStmt.LineNumber, repeatStmt.StatementNumber)
			return nil
		}
	}
	return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UntilWithoutAnyRepeat), ErrorTokenIndex: stmt.Token.Index}
}

func evalSetRadStatement(g *game.Game, stmt *ast.SetRadStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	if isError(obj) {
		return obj
	}
	if val, ok := obj.(*object.Boolean); ok {
		env.Degrees = !val.Value
		return obj
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
}

func evalListStatement(g *game.Game, stmt *ast.ListStatement, env *object.Environment) object.Object {
	oldTextBoxSlot, _, _, _, _ := g.AskWriting()
	tempTextBoxSlot := oldTextBoxSlot
	_, curY := g.AskCurpos()
	// Evaluate and handle TextBoxSlot if set
	if stmt.TextBoxSlot != nil {
		obj := Eval(g, stmt.TextBoxSlot, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			tempTextBoxSlot = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	g.SetWriting(tempTextBoxSlot)

	fromLinenumber := 0
	toLinenumber := 0
	if stmt.FromLinenumber.Literal != "" {
		val, _ := strconv.ParseFloat(stmt.FromLinenumber.Literal, 64)
		fromLinenumber = int(val)
	}
	if stmt.ToLinenumber.Literal != "" {
		val, _ := strconv.ParseFloat(stmt.ToLinenumber.Literal, 64)
		toLinenumber = int(val)
	}
	listing := env.Program.List(fromLinenumber, toLinenumber, stmt.FromLineOnly)
	if listing == nil {
		if oldTextBoxSlot != tempTextBoxSlot {
			g.SetWriting(oldTextBoxSlot)
			g.SetCurpos(1, curY)
		}
		return nil
	}
	for _, listString := range listing {
		g.Print(listString)
		g.Put(13)
	}
	if oldTextBoxSlot != tempTextBoxSlot {
		g.SetWriting(oldTextBoxSlot)
		g.SetCurpos(1, curY)
	}
	return nil
}

func prerun(g *game.Game, env *object.Environment) bool {
	// Run through the stored program without executing instructions.  Instead
	// register all functions, procedures, subroutines and collect data.
	l := &lexer.Lexer{}
	env.Program.Start()
	env.DeleteData()
	env.DeleteSubroutines()
	env.DeleteFunctions()
	env.DeleteProcedures()
	env.Prerun = true
	for !env.Program.EndOfProgram() {
		//log.Printf("%s", env.Program.GetLine())
		l.Scan(env.Program.GetLine())
		p := parser.New(l, g)
		line := p.ParseLine()
		// Handle parsing error here --> need some tweaks
		if errorMsg, hasError := p.GetError(); hasError {
			g.Print(errorMsg)
			g.Put(13)
			p.JumpToToken(0)
			g.Print(p.PrettyPrint())
			g.Put(13)
			return false
		}
		// Only evaluate the following statements:
		// FUNCTION, PROCEDURE, SUBROUTINE, DATA
		for statementNumber, stmt := range line.Statements {
			env.Program.CurrentStatementNumber = statementNumber
			tokenType := stmt.TokenLiteral()
			// Capture DATA statements, FUNCTION and PROCEDURE statements, and SUBROUTINE statements
			if tokenType == token.DATA || tokenType == token.SUBROUTINE || tokenType == token.FUNCTION || tokenType == token.PROCEDURE {
				obj := Eval(g, stmt, env)
				// Handle eval error
				if errorMsg, ok := obj.(*object.Error); ok {
					if errorMsg.ErrorTokenIndex != 0 {
						p.ErrorTokenIndex = errorMsg.ErrorTokenIndex
					}
					lineNumber := env.Program.GetLineNumber()
					g.Print(fmt.Sprintf("%s in line %d", errorMsg.Message, lineNumber))
					g.Put(13)
					p.JumpToToken(0)
					g.Print(fmt.Sprintf("%d %s", lineNumber, p.PrettyPrint()))
					g.Put(13)
					return false
				}
			}
		}
		env.Program.Next()
	}
	return true
}

func evalRunStatement(g *game.Game, stmt *ast.RunStatement, env *object.Environment) object.Object {
	// Prerun stored program and return if prerun failed
	if !prerun(g, env) {
		return nil
	}
	// Otherwise execute stored program
	l := &lexer.Lexer{}
	env.Prerun = false
	env.Program.Start()
	env.EndProgramSignal = false
	env.LeaveFunctionSignal = false
	// If a line number was passed, attempt to jump to it and return error if this fails
	if stmt.Linenumber.Literal != "" {
		val, _ := strconv.ParseFloat(stmt.Linenumber.Literal, 64)
		lineNumber := int(val)
		if lineNumber < 0 {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.PositiveValueRequired), ErrorTokenIndex: stmt.Linenumber.Index}
		}
		if !env.Program.Jump(lineNumber, 0) {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.LineNumberDoesNotExist), ErrorTokenIndex: stmt.Linenumber.Index}
		} else {
			env.Program.Next()
		}
	}
	// And away we go
	for !env.Program.EndOfProgram() && !g.BreakInterruptDetected && !env.EndProgramSignal {
		//log.Printf("%s", env.Program.GetLine())
		l.Scan(env.Program.GetLine())
		p := parser.New(l, g)
		line := p.ParseLine()
		// Check of parser errors here.  Parser errors are handled just like evaluation errors but
		// obviously we'll skip evaluation if parsing already failed.
		if errorMsg, hasError := p.GetError(); hasError {
			lineNumber := env.Program.GetLineNumber()
			g.Print(fmt.Sprintf("%s in line %d", errorMsg, lineNumber))
			g.Put(13)
			p.JumpToToken(0)
			g.Print(fmt.Sprintf("%d %s", lineNumber, p.PrettyPrint()))
			g.Put(13)
			return nil
		}
		// And this is temporary while we're still migrating from Monkey to RM Basic
		if len(p.Errors()) > 0 {
			g.Print("Oops! Some random parsing error occurred. These will be handled properly downstream by for now here's some spewage:")
			g.Put(13)
			p.JumpToToken(0)
			g.Print(p.PrettyPrint())
			g.Put(13)
			for _, msg := range p.Errors() {
				g.Print(msg)
				g.Put(13)
			}
			return nil
		}
		// Execute each statement in the program line.  If an error occurs, print the
		// error message and stop.  If JumpToStatement is non-zero, all statements in
		// the line will be skipped until i == JumpToStatement.
		for statementNumber, stmt := range line.Statements {
			env.Program.CurrentStatementNumber = statementNumber
			obj := Eval(g, stmt, env)
			if errorMsg, ok := obj.(*object.Error); ok {
				if errorMsg.ErrorTokenIndex != 0 {
					p.ErrorTokenIndex = errorMsg.ErrorTokenIndex
				}
				lineNumber := env.Program.GetLineNumber()
				g.Print(fmt.Sprintf("%s in line %d", errorMsg.Message, lineNumber))
				g.Put(13)
				p.JumpToToken(0)
				g.Print(fmt.Sprintf("%d %s", lineNumber, p.PrettyPrint()))
				g.Put(13)
				return nil
			}
			if g.BreakInterruptDetected {
				break
			}
		}
		env.Program.Next()
	}
	if g.BreakInterruptDetected {
		g.Print(fmt.Sprintf("%s in line %d", syntaxerror.ErrorMessage(syntaxerror.InterruptedByBreakKey), env.Program.GetLineNumber()))
		g.Put(13)
		time.Sleep(150 * time.Millisecond)
	}
	return nil
}

func evalClsStatement(g *game.Game, stmt *ast.ClsStatement, env *object.Environment) object.Object {
	g.Cls()
	g.SetCurpos(1, 1)
	return nil
}

func evalClearblockStatement(g *game.Game, stmt *ast.ClearblockStatement, env *object.Environment) object.Object {
	g.Clearblock()
	return nil
}

func evalSetMouseStatement(g *game.Game, stmt *ast.SetMouseStatement, env *object.Environment) object.Object {
	g.SetMouse(true)
	return nil
}

func evalHomeStatement(g *game.Game, stmt *ast.HomeStatement, env *object.Environment) object.Object {
	g.SetCurpos(0, 0)
	return nil
}

func getAbsPath(p string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %s", err)
	}
	return filepath.Join(wd, filepath.FromSlash(strings.ReplaceAll(p, "\\", "/")))
}

func evalDirStatement(g *game.Game, stmt *ast.DirStatement, env *object.Environment) object.Object {
	oldTextBoxSlot, _, _, _, _ := g.AskWriting()
	tempTextBoxSlot := oldTextBoxSlot
	_, curY := g.AskCurpos()
	// Evaluate and handle TextBoxSlot if set
	if stmt.TextBoxSlot != nil {
		obj := Eval(g, stmt.TextBoxSlot, env)
		if isError(obj) {
			return obj
		}
		if val, ok := obj.(*object.Numeric); ok {
			tempTextBoxSlot = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	g.SetWriting(tempTextBoxSlot)

	// evaluate path if given
	val := ""
	if stmt.Value != nil {
		obj := Eval(g, stmt.Value, env)
		if isError(obj) {
			if oldTextBoxSlot != tempTextBoxSlot {
				g.SetWriting(oldTextBoxSlot)
				g.SetCurpos(1, curY)
			}
			return obj
		}
		if stringVal, ok := obj.(*object.String); ok {
			val = stringVal.Value
		} else {
			if oldTextBoxSlot != tempTextBoxSlot {
				g.SetWriting(oldTextBoxSlot)
				g.SetCurpos(1, curY)
			}
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// add *.BAS if no extension given
	if !strings.Contains(val, ".") {
		if !strings.HasSuffix(val, "\\") && len(val) > 0 {
			val += "\\*.BAS"
		} else {
			val += "*.BAS"
		}
	}
	systemPath := getAbsPath(val)
	nimbusPath := strings.ReplaceAll(systemPath[len(g.WorkspacePath):], "/", "\\")
	files, err := filepath.Glob(systemPath)
	if err != nil {
		if oldTextBoxSlot != tempTextBoxSlot {
			g.SetWriting(oldTextBoxSlot)
			g.SetCurpos(1, curY)
		}
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.DirectoryCannotBeFound), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	g.Print(fmt.Sprintf("Directory of %s", nimbusPath))
	g.Put(13)
	g.Put(13)
	// list subdirs first - get rid of wildcard expression if present
	var subdirsSystemPath string
	pathStrings := strings.Split(systemPath, "*")
	for _, pathString := range pathStrings {
		if strings.HasPrefix(pathString, ".") {
			break
		} else {
			subdirsSystemPath += pathString
		}
	}
	dirs, err := ioutil.ReadDir(subdirsSystemPath)
	if err != nil {
		if oldTextBoxSlot != tempTextBoxSlot {
			g.SetWriting(oldTextBoxSlot)
			g.SetCurpos(1, curY)
		}
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	for _, d := range dirs {
		if d.IsDir() {
			dirName := d.Name()
			dirDate := d.ModTime().Format("2006-01-02 15:04:05")
			dirString := fmt.Sprintf("%16s %6s       %16s", dirName, "<DIR>", dirDate)
			g.Print(dirString)
			g.Put(13)
		}
	}
	// then list files
	for _, f := range files {
		fileInfo, err := os.Stat(f)
		if err != nil {
			if oldTextBoxSlot != tempTextBoxSlot {
				g.SetWriting(oldTextBoxSlot)
				g.SetCurpos(1, curY)
			}
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure), ErrorTokenIndex: stmt.Token.Index + 1}
		}
		fileSize := fileInfo.Size()
		fileName := filepath.Base(f)
		fileDate := fileInfo.ModTime().Format("2006-01-02 15:04:05")
		var dirString string
		if !fileInfo.IsDir() {
			dirString = fmt.Sprintf("%16s %6d Bytes %16s", fileName, fileSize, fileDate)
			g.Print(dirString)
			g.Put(13)
		}
	}
	g.Put(13)
	if oldTextBoxSlot != tempTextBoxSlot {
		g.SetWriting(oldTextBoxSlot)
		g.SetCurpos(1, curY)
	}
	return nil
}

func evalChdirStatement(g *game.Game, stmt *ast.ChdirStatement, env *object.Environment) object.Object {
	// evaluate path if given
	val := ""
	if stmt.Value != nil {
		obj := Eval(g, stmt.Value, env)
		if isError(obj) {
			return obj
		}
		if stringVal, ok := obj.(*object.String); ok {
			val = stringVal.Value
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// execute
	systemPath := getAbsPath(val)
	// Special case of "/" indicating user wants to go back to root (which is really the workspace dir)
	if val == "\\" {
		systemPath = g.WorkspacePath
	}
	err := os.Chdir(systemPath)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.DirectoryCannotBeFound), ErrorTokenIndex: stmt.Token.Index + 1}
	} else {
		return nil
	}
}

func evalMkdirStatement(g *game.Game, stmt *ast.MkdirStatement, env *object.Environment) object.Object {
	// evaluate path if given
	val := ""
	if stmt.Value != nil {
		obj := Eval(g, stmt.Value, env)
		if isError(obj) {
			return obj
		}
		if stringVal, ok := obj.(*object.String); ok {
			val = stringVal.Value
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// execute
	systemPath := getAbsPath(val)
	// Special case of "/" indicating user wants to go back to root (which is really the workspace dir)
	if val == "\\" {
		systemPath = g.WorkspacePath
	}
	err := os.Mkdir(systemPath, 0755)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnableToCreateDirectory), ErrorTokenIndex: stmt.Token.Index + 1}
	} else {
		return nil
	}
}

func evalRmdirStatement(g *game.Game, stmt *ast.RmdirStatement, env *object.Environment) object.Object {
	// evaluate path if given
	val := ""
	if stmt.Value != nil {
		obj := Eval(g, stmt.Value, env)
		if isError(obj) {
			return obj
		}
		if stringVal, ok := obj.(*object.String); ok {
			val = stringVal.Value
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// execute
	systemPath := getAbsPath(val)
	// Special case of "/" indicating user wants to go back to root (which is really the workspace dir)
	if val == "\\" {
		systemPath = g.WorkspacePath
	}
	// ensure systemPath is a directory and not a file
	fileInfo, err := os.Stat(systemPath)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if !fileInfo.IsDir() {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.DirectoryCannotBeFound), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	err = os.Remove(systemPath)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnableToRemoveDirectory), ErrorTokenIndex: stmt.Token.Index + 1}
	} else {
		return nil
	}
}

func evalEraseStatement(g *game.Game, stmt *ast.EraseStatement, env *object.Environment) object.Object {
	// evaluate path if given
	val := ""
	if stmt.Value != nil {
		obj := Eval(g, stmt.Value, env)
		if isError(obj) {
			return obj
		}
		if stringVal, ok := obj.(*object.String); ok {
			val = stringVal.Value
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Don't allow * or ?
	if strings.Contains(val, "*") || strings.Contains(val, "?") {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ExactFilenameIsNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Add .BAS if necessary
	if !strings.HasSuffix(strings.ToUpper(val), ".BAS") {
		val += ".BAS"
	}
	// execute
	systemPath := getAbsPath(val)
	// Special case of "/" indicating user wants to go back to root (which is really the workspace dir)
	if val == "\\" {
		systemPath = g.WorkspacePath
	}
	// ensure systemPath is a file and not a directory
	fileInfo, err := os.Stat(systemPath)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnableToEraseTheFile), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if fileInfo.IsDir() {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FilenameIsADirectory), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	err = os.Remove(systemPath)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnableToEraseTheFile), ErrorTokenIndex: stmt.Token.Index + 1}
	} else {
		return nil
	}
}

func evalRenameStatement(g *game.Game, stmt *ast.RenameStatement, env *object.Environment) object.Object {
	// evaluate filename1
	val1 := ""
	if stmt.Value1 != nil {
		obj := Eval(g, stmt.Value1, env)
		if isError(obj) {
			return obj
		}
		if stringVal, ok := obj.(*object.String); ok {
			val1 = stringVal.Value
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Don't allow * or ?
	if strings.Contains(val1, "*") || strings.Contains(val1, "?") {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ExactFilenameIsNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Add .BAS if necessary
	if !strings.HasSuffix(strings.ToUpper(val1), ".BAS") {
		val1 += ".BAS"
	}
	// ensure val1 is a file and not a directory
	fileInfo, err := os.Stat(val1)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnableToRenameTheFile), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	if fileInfo.IsDir() {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FilenameIsADirectory), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// evaluate filename2
	val2 := ""
	if stmt.Value2 != nil {
		obj := Eval(g, stmt.Value2, env)
		if isError(obj) {
			return obj
		}
		if stringVal, ok := obj.(*object.String); ok {
			val2 = stringVal.Value
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	// Don't allow * or ?
	if strings.Contains(val2, "*") || strings.Contains(val2, "?") {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.ExactFilenameIsNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Add .BAS if necessary
	if !strings.HasSuffix(strings.ToUpper(val2), ".BAS") {
		val2 += ".BAS"
	}
	// execute
	systemPath1 := getAbsPath(val1)
	systemPath2 := getAbsPath(val2)
	// rename file 1 to file 2
	err = os.Rename(systemPath1, systemPath2)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnableToRenameTheFile), ErrorTokenIndex: stmt.Token.Index + 1}
	} else {
		return nil
	}
}

func evalExpressions(g *game.Game, exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(g, e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// canObjectCastToIdentifierType checks if the object can be cast to the identifier type.  If so it
// returns the cast value and true, otherwise it return an error object and false.
func canObjectCastToIdentifierType(obj object.Object, identifierName string) (object.Object, bool) {
	if identifierName[len(identifierName)-1:] == "$" {
		// object must be string type
		if obj.Type() != object.STRING_OBJ {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: 0}, false
		} else {
			return obj, true
		}
	}
	if identifierName[len(identifierName)-1:] == "%" {
		// object must be numeric type and cast to integer
		if obj.Type() != object.STRING_OBJ {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: 0}, false
		} else {
			val := int(obj.(*object.Numeric).Value)
			return &object.Numeric{Value: float64(val)}, true
		}
	}
	// object must be numeric but not casting to integer
	if obj.Type() != object.NUMERIC_OBJ {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: 0}, false
	} else {
		return obj, true
	}
}

// execute runs the function or procedure code and returns the return vals
func executeFunction(g *game.Game, env *object.Environment, startLine int, statementNumber int, skipFirstStatement bool) []object.Object {
	// jump to position and execute
	env.Program.Jump(startLine, statementNumber)
	env.Program.Next()
	l := &lexer.Lexer{}
	env.Prerun = false
	for !env.Program.EndOfProgram() && !g.BreakInterruptDetected && !env.LeaveFunctionSignal && !env.EndProgramSignal {
		l.Scan(env.Program.GetLine())
		p := parser.New(l, g)
		line := p.ParseLine()
		// Check of parser errors here.  Parser errors are handled just like evaluation errors but
		// obviously we'll skip evaluation if parsing already failed.
		if errorMsg, hasError := p.GetError(); hasError {
			lineNumber := env.Program.GetLineNumber()
			g.Print(fmt.Sprintf("%s in line %d", errorMsg, lineNumber))
			g.Put(13)
			p.JumpToToken(0)
			g.Print(fmt.Sprintf("%d %s", lineNumber, p.PrettyPrint()))
			g.Put(13)
			return nil
		}
		// And this is temporary while we're still migrating from Monkey to RM Basic
		if len(p.Errors()) > 0 {
			g.Print("Oops! Some random parsing error occurred. These will be handled properly downstream by for now here's some spewage:")
			g.Put(13)
			p.JumpToToken(0)
			g.Print(p.PrettyPrint())
			g.Put(13)
			for _, msg := range p.Errors() {
				g.Print(msg)
				g.Put(13)
			}
			return nil
		}
		// Execute each statement in the program line.  If an error occurs, print the
		// error message and stop.  If JumpToStatement is non-zero, all statements in
		// the line will be skipped until i == JumpToStatement.
		for statementNumber, stmt := range line.Statements {
			if skipFirstStatement {
				skipFirstStatement = false
				continue
			}
			env.Program.CurrentStatementNumber = statementNumber
			obj := Eval(g, stmt, env)
			if errorMsg, ok := obj.(*object.Error); ok {
				if errorMsg.ErrorTokenIndex != 0 {
					p.ErrorTokenIndex = errorMsg.ErrorTokenIndex
				}
				lineNumber := env.Program.GetLineNumber()
				g.Print(fmt.Sprintf("%s in line %d", errorMsg.Message, lineNumber))
				g.Put(13)
				p.JumpToToken(0)
				g.Print(fmt.Sprintf("%d %s", lineNumber, p.PrettyPrint()))
				g.Put(13)
				return nil
			}
			if g.BreakInterruptDetected {
				break
			}
		}
		env.Program.Next()
	}
	return env.ReturnVals
}

func evalProcedureCallStatement(g *game.Game, stmt *ast.ProcedureCallStatement, env *object.Environment) object.Object {
	if proc, ok := env.GetProcedure(stmt.Name.Value); ok {
		args := make([]object.Object, len(stmt.Args))
		for i := 0; i < len(stmt.Args); i++ {
			args[i] = Eval(g, stmt.Args[i], env)
			if isError(args[i]) {
				return args[i]
			}
		}
		newEnv := object.NewEnvironment()
		newEnv.Copy(env.Dump())
		newEnv.NewScope()
		for i := 0; i < len(proc.ReceiveArgs); i++ {
			obj := newEnv.Set(proc.ReceiveArgs[i].Value, args[i])
			if isError(obj) {
				return obj
			}
		}
		retVal := executeFunction(g, newEnv, proc.LineNumber, proc.StatementNumber, true)[0]
		if newEnv.EndProgramSignal {
			env.EndProgram()
		}
		if isError(retVal) {
			return retVal
		}
		for i := 0; i < len(stmt.ReceiveArgs); i++ {
			val, _ := newEnv.Get(proc.ReturnArgs[i].Value)
			obj := env.Set(stmt.ReceiveArgs[i].Value, val)
			if isError(obj) {
				return obj
			}
		}
		return nil
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.UnknownCommandProcedure), ErrorTokenIndex: stmt.Token.Index}
	}
}

func evalIdentifier(g *game.Game, node *ast.Identifier, env *object.Environment) object.Object {
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	if len(node.Subscripts) > 0 {
		if fun, ok := env.GetFunction(node.Value); ok {
			// Handle function
			subscripts := make([]object.Object, len(node.Subscripts))
			for i := 0; i < len(node.Subscripts); i++ {
				subscripts[i] = Eval(g, node.Subscripts[i], env)
				if isError(subscripts[i]) {
					return subscripts[i]
				}
			}
			newEnv := object.NewEnvironment()
			newEnv.Copy(env.Dump())
			newEnv.NewScope()
			for i := 0; i < len(node.Subscripts); i++ {
				obj := newEnv.Set(fun.ReceiveArgs[i].Value, subscripts[i])
				if isError(obj) {
					return obj
				}
			}
			retVal := executeFunction(g, newEnv, fun.LineNumber, fun.StatementNumber, true)[0]
			if newEnv.EndProgramSignal {
				env.EndProgram()
			}
			if isError(retVal) {
				return retVal
			} else {
				obj, _ := canObjectCastToIdentifierType(retVal, fun.Name.Value)
				return obj
			}
		} else {
			// handle array
			subscripts := make([]int, len(node.Subscripts))
			for i := 0; i < len(subscripts); i++ {
				obj := Eval(g, node.Subscripts[i], env)
				if val, ok := obj.(*object.Numeric); ok {
					subscripts[i] = int(val.Value)
				} else {
					return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: node.Token.Index}
				}
			}
			// Error handling here?
			val, _ := env.GetArray(node.Value, subscripts)
			return val
		}
	}
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	// Create a new variable with null value and return warning.  It's then up to the caller
	// to print the warning and do env.Get again to get the value.
	if node.Value[len(node.Value)-1:] == "$" {
		env.Set(node.Value, &object.String{Value: ""})
	} else {
		env.Set(node.Value, &object.Numeric{Value: 0})
	}
	return &object.Warning{Message: "Variable without any value"}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalBlockStatement(g *game.Game, block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(g, statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalIfStatement(g *game.Game, ie *ast.IfStatement, env *object.Environment) object.Object {
	condition := Eval(g, ie.Condition, env)
	if isError(condition) {
		return condition
	}
	var returnObject object.Object
	if isTruthy(condition) {
		for _, stmt := range ie.Consequence.Statements {
			obj := Eval(g, stmt, env)
			returnObject = obj
			if isError(obj) {
				return obj
			}
		}
		//return NULL
		return returnObject
	} else if ie.Alternative != nil {
		for _, stmt := range ie.Alternative.Statements {
			obj := Eval(g, stmt, env)
			returnObject = obj
			if isError(obj) {
				return obj
			}
		}
	}
	//return NULL
	return returnObject
}

func isTruthy(obj object.Object) bool {
	val := obj.(*object.Numeric).Value
	if val == -1.0 {
		return true
	} else {
		return false
	}
	//switch obj {
	//case NULL:
	//	return false
	//case TRUE:
	//	return true
	//case FALSE:
	//	return false
	//default:
	//	return true
	//}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.NUMERIC_OBJ && right.Type() == object.NUMERIC_OBJ:
		return evalNumericInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "=":
		return nativeBoolToBooleanObject(left == right)
	case operator == "=":
		return nativeBoolToBooleanObject(left == right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("%s (type mismatch: %s %s %s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), left.Type(), operator, right.Type())
	default:
		return newError("%s (unknown operator: %s %s %s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case operator == "+":
		// concatenation
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	case operator == "=":
		// exactly equal, i.e. case sensitive comparison
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		if leftVal == rightVal {
			return &object.Numeric{Value: -1.0}
		} else {
			return &object.Numeric{Value: 0}
		}
	case operator == "==":
		// "interestingly equal", i.e. case-insensitive comparison
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		if strings.EqualFold(leftVal, rightVal) {
			return &object.Numeric{Value: -1.0}
		} else {
			return &object.Numeric{Value: 0}
		}
	default:
		return newError("%s (unknown operator: %s %s %s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), left.Type(), operator, right.Type())
	}

}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Numeric).Value
	rightVal := right.(*object.Numeric).Value

	switch operator {
	case "AND":
		return &object.Numeric{Value: float64(int(leftVal) & int(rightVal))}
	case "OR":
		return &object.Numeric{Value: float64(int(leftVal) | int(rightVal))}
	case "XOR":
		return &object.Numeric{Value: float64(int(leftVal) ^ int(rightVal))}
	default:
		return newError("%s (unknown operator: %s %s %s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), left.Type(), operator, right.Type())
	}
}

func evalNumericInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Numeric).Value
	rightVal := right.(*object.Numeric).Value

	switch operator {
	case "+":
		return &object.Numeric{Value: leftVal + rightVal}
	case "-":
		return &object.Numeric{Value: leftVal - rightVal}
	case "*":
		return &object.Numeric{Value: leftVal * rightVal}
	case "/":
		// catch divide by zero
		if rightVal == 0 {
			return newError(syntaxerror.ErrorMessage(syntaxerror.TryingToDivideByZero))
		} else {
			return &object.Numeric{Value: leftVal / rightVal}
		}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "=<":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "<>":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "><":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "=>":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "=":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "AND":
		return &object.Numeric{Value: float64(int(leftVal) & int(rightVal))}
	case "OR":
		return &object.Numeric{Value: float64(int(leftVal) | int(rightVal))}
	case "XOR":
		return &object.Numeric{Value: float64(int(leftVal) ^ int(rightVal))}
	case "MOD":
		// catch divide by zero
		if rightVal == 0 {
			return newError(syntaxerror.ErrorMessage(syntaxerror.TryingToDivideByZero))
		} else {
			return &object.Numeric{Value: math.Mod(leftVal, rightVal)}
		}
	default:
		return newError("%s (unknown operator: %s %s %s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "NOT":
		return evalNotOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("%s (unknown operator: %s%s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), operator, right.Type())
	}
}

func evalNotOperatorExpression(right object.Object) object.Object {
	val := right.(*object.Numeric).Value
	return &object.Numeric{Value: float64(^int(val))}
	//switch right {
	//case TRUE:
	//	return FALSE
	//case FALSE:
	//	return TRUE
	//case NULL:
	//	return TRUE
	//default:
	//	return FALSE
	//}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMERIC_OBJ {
		return newError("%s (unknown operator: -%s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), right.Type())
	}
	value := right.(*object.Numeric).Value
	return &object.Numeric{Value: -value}
}

func nativeBoolToBooleanObject(input bool) *object.Numeric {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...)}
}

func evalProgram(g *game.Game, program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(g, statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}
