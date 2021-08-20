package evaluator

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
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
	case *ast.ListStatement:
		return evalListStatement(g, node, env)
	case *ast.ClsStatement:
		return evalClsStatement(g, node, env)
	case *ast.HomeStatement:
		return evalHomeStatement(g, node, env)
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
	case *ast.GotoStatement:
		return evalGotoStatement(g, node, env)
	case *ast.RepeatStatement:
		return evalRepeatStatement(g, node, env)
	case *ast.UntilStatement:
		return evalUntilStatement(g, node, env)
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
		return evalInfixExpression(node.Operator, left, right)
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
	var Brush, Direction, Font, Over, SizeX, SizeY, X, Y int
	var Text string
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if val, ok := obj.(*object.Numeric); ok {
			Brush = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Direction == nil {
		Direction = -255
	} else {
		obj := Eval(g, stmt.Direction, env)
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
	// Handle text and coordinates
	obj := Eval(g, stmt.X, env)
	if val, ok := obj.(*object.Numeric); ok {
		X = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.Y, env)
	if val, ok := obj.(*object.Numeric); ok {
		Y = int(val.Value)
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	obj = Eval(g, stmt.Value, env)
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
	// Execute
	opt := nimgobus.PlotOptions{Brush: Brush, Direction: Direction, Font: Font, SizeX: SizeX, SizeY: SizeY, Over: Over}
	g.Plot(opt, Text, X, Y)
	return nil
}

func evalLineStatement(g *game.Game, stmt *ast.LineStatement, env *object.Environment) object.Object {
	// Handle defaults
	var Brush, Over int
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if val, ok := obj.(*object.Numeric); ok {
			Brush = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Over == nil {
		Over = -255
	} else {
		obj := Eval(g, stmt.Over, env)
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

func evalAreaStatement(g *game.Game, stmt *ast.AreaStatement, env *object.Environment) object.Object {
	// Handle defaults
	var Brush, Over int
	if stmt.Brush == nil {
		Brush = -255
	} else {
		obj := Eval(g, stmt.Brush, env)
		if val, ok := obj.(*object.Numeric); ok {
			Brush = int(val.Value)
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	}
	if stmt.Over == nil {
		Over = -255
	} else {
		obj := Eval(g, stmt.Over, env)
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
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	filename := ""
	if stringVal, ok := obj.(*object.String); ok {
		filename = stringVal.Value
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Add .BAS if necessary
	if !strings.HasSuffix(filename, ".BAS") {
		filename += ".BAS"
	}
	// Save the program
	file, err := os.Create(filename)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure)}
	}
	defer file.Close()
	for _, lineString := range env.Program.List() {
		file.WriteString(fmt.Sprintf("%s\n", lineString))
	}
	return obj
}

func evalLoadStatement(g *game.Game, stmt *ast.LoadStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	filename := ""
	if stringVal, ok := obj.(*object.String); ok {
		filename = stringVal.Value
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
	// Add .BAS if necessary
	if !strings.HasSuffix(filename, ".BAS") {
		filename += ".BAS"
	}
	// Load the program
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.FileOperationFailure)}
	}
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

func evalSetModeStatement(g *game.Game, stmt *ast.SetModeStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
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
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		// We have to restrict the value range depending on screen mode.  RM Basic didn't quite handle it
		// like this so TODO is to implement this properly.
		highestColour := 3
		if g.AskMode() == 40 {
			highestColour = 15
		}
		if val.Value < 0 {
			val.Value = 0
		}
		if val.Value > float64(highestColour) {
			val.Value = float64(highestColour)
		}
		g.SetPaper(int(val.Value))
		return obj
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
}

func evalSetBorderStatement(g *game.Game, stmt *ast.SetBorderStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		// We have to restrict the value range depending on screen mode.  RM Basic didn't quite handle it
		// like this so TODO is to implement this properly.
		highestColour := 3
		if g.AskMode() == 40 {
			highestColour = 15
		}
		if val.Value < 0 {
			val.Value = 0
		}
		if val.Value > float64(highestColour) {
			val.Value = float64(highestColour)
		}
		g.SetBorder(int(val.Value))
		return obj
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
}

func evalSetPenStatement(g *game.Game, stmt *ast.SetPenStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		// We have to restrict the value range depending on screen mode.  RM Basic didn't quite handle it
		// like this so TODO is to implement this properly.
		highestColour := 3
		if g.AskMode() == 40 {
			highestColour = 15
		}
		if val.Value < 0 {
			val.Value = 0
		}
		if val.Value > float64(highestColour) {
			val.Value = float64(highestColour)
		}
		g.SetPen(int(val.Value))
		return obj
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
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
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if _, ok := obj.(*object.Numeric); ok {
		env.Degrees = isTruthy(obj)
		return obj
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
}

func evalSetConfigBootStatement(g *game.Game, stmt *ast.SetConfigBootStatement, env *object.Environment) object.Object {
	obj := Eval(g, stmt.Value, env)
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
	obj := Eval(g, stmt.Value, env)
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	if val, ok := obj.(*object.Numeric); ok {
		if env.Program.Jump(int(val.Value), 0) {
			return obj
		} else {
			return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.LineNumberDoesNotExist), ErrorTokenIndex: stmt.Token.Index + 1}
		}
	} else {
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded), ErrorTokenIndex: stmt.Token.Index + 1}
	}
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
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
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

func evalHomeStatement(g *game.Game, stmt *ast.HomeStatement, env *object.Environment) object.Object {
	g.SetCurpos(0, 0)
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
			if _, ok := obj.(*object.Error); ok {
				return obj
			}
		}
		return NULL
	} else if ie.Alternative != nil {
		for _, stmt := range ie.Alternative.Statements {
			obj := Eval(g, stmt, env)
			if _, ok := obj.(*object.Error); ok {
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
		return &object.Numeric{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
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
