package evaluator

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	case *ast.AskMouseStatement:
		return evalAskMouseStatement(g, node, env)
	case *ast.Program:
		return evalProgram(g, node, env)
	case *ast.ExpressionStatement:
		return Eval(g, node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(g, node, env)
	case *ast.IfStatement:
		return evalIfStatement(g, node, env)
	case *ast.ReturnStatement:
		val := Eval(g, node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{
			Value: val,
		}
	case *ast.LetStatement:
		val := Eval(g, node.Value, env)
		if isError(val) {
			return val
		}
		return env.Set(node.Name.Value, val)
	case *ast.BindStatement:
		val := Eval(g, node.Value, env)
		if isError(val) {
			return val
		}
		return env.Set(node.Name.Value, val)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
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
		// If a warning is returned, print the warning *then* re-run the evaluation and return
		obj := evalIdentifier(node, env)
		if warningMsg, ok := obj.(*object.Warning); ok {
			g.Print(fmt.Sprintf("Warning: %s", warningMsg.Message))
			g.Put(13)
			return evalIdentifier(node, env)
		} else {
			return obj
		}
	case *ast.CallExpression:
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
	opt := nimgobus.CircleOptions{Brush: Brush, Over: Over}
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
	opt := nimgobus.FloodOptions{Brush: Brush, UseEdgeColour: UseEdgeColour, EdgeColour: EdgeColour}
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
	opt := nimgobus.AreaOptions{Brush: Brush, Over: Over}
	g.Area(opt, coordList)
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
	for _, lineString := range env.Program.List() {
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
	obj = Eval(g, stmt.X, env)
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
	obj = Eval(g, stmt.X, env)
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
	listing := env.Program.List()
	if listing == nil {
		return nil
	}
	for _, listString := range listing {
		g.Print(listString)
		g.Put(13)
	}
	return nil
}

func evalRunStatement(g *game.Game, stmt *ast.RunStatement, env *object.Environment) object.Object {
	l := &lexer.Lexer{}
	env.Program.Start()
	for !env.Program.EndOfProgram() && !g.BreakInterruptDetected {
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

func evalDirStatement(g *game.Game, stmt *ast.DirStatement, env *object.Environment) object.Object {
	// TODO: Handle select different path
	files, err := ioutil.ReadDir(g.WorkspacePath)
	if err != nil {
		return nil // TODO: io error
	}
	g.Print(fmt.Sprintf("Directory of %s", g.WorkspacePath))
	g.Put(13)
	g.Put(13)
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".BAS") {
			continue
		}
		timeString := fmt.Sprintf("%s", f.ModTime().Round(time.Second))[:19]
		g.Print(fmt.Sprintf("%16s %6d Bytes %16s", f.Name(), f.Size(), timeString))
		g.Put(13)
	}
	g.Put(13)
	return nil
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	// Create a new variable with null value and return warning.  It's then up to the caller
	// to print the warning and do env.Get again to get the value.
	//name[len(name)-1:] != "$"
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
	if isTruthy(condition) {
		for _, stmt := range ie.Consequence.Statements {
			obj := Eval(g, stmt, env)
			if isError(obj) {
				return obj
			}
		}
		return NULL
	} else if ie.Alternative != nil {
		for _, stmt := range ie.Alternative.Statements {
			obj := Eval(g, stmt, env)
			if isError(obj) {
				return obj
			}
		}
	}
	return NULL
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
