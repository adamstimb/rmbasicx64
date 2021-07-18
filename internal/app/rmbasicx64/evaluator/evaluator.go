package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/ast"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/game"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/object"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// Because null and boolean values never change we can reference them instead of
// creating new objects
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(g *game.Game, node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.ByeStatement:
		os.Exit(0)
	case *ast.ClsStatement:
		return evalClsStatement(g, node, env)
	case *ast.SetModeStatement:
		return evalSetModeStatement(g, node, env)
	case *ast.SetPaperStatement:
		return evalSetPaperStatement(g, node, env)
	case *ast.SetBorderStatement:
		return evalSetBorderStatement(g, node, env)
	case *ast.SetPenStatement:
		return evalSetPenStatement(g, node, env)
	case *ast.PrintStatement:
		return evalPrintStatement(g, node, env)
	case *ast.Program:
		return evalProgram(g, node, env)
	case *ast.ExpressionStatement:
		return Eval(g, node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(g, node, env)
	case *ast.IfExpression:
		return evalIfExpression(g, node, env)
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
		val := fn.Fn(args...)
		// If the builtin is a trig function we need to catch the result and convert
		// to deg if env.Degrees is true
		if env.Degrees {
			// check if trig and convert radians to deg
			if fn == builtins["ATN"] {
				val.(*object.Numeric).Value *= (180 / math.Pi)
			}
		}
		return val
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
	obj := Eval(g, stmt.Value, env)
	// return error if evaluation failed
	if _, ok := obj.(*object.Error); ok {
		return obj
	}
	printStr := ""
	if numericVal, ok := obj.(*object.Numeric); ok {
		printStr = fmt.Sprintf("%g", numericVal.Value)
	}
	if boolVal, ok := obj.(*object.Boolean); ok {
		if boolVal.Value {
			printStr = "TRUE"
		} else {
			printStr = "FALSE"
		}
	}
	if stringVal, ok := obj.(*object.String); ok {
		printStr = stringVal.Value
	}
	g.Print(printStr)
	g.Put(13)
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
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)}
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
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)}
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
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)}
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
		return &object.Error{Message: syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)}
	}
}

func evalClsStatement(g *game.Game, stmt *ast.ClsStatement, env *object.Environment) object.Object {
	g.Cls()
	g.SetCurpos(1, 1)
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
	//return newError("identifier not found: " + node.Value)
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

func evalIfExpression(g *game.Game, ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(g, ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(g, ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(g, ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.NUMERIC_OBJ && right.Type() == object.NUMERIC_OBJ:
		return evalNumericInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "=":
		return nativeBoolToBooleanObject(left == right)
		// TODO: case operator == "!=": etc.
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
		return &object.Boolean{Value: leftVal == rightVal}
	case operator == "==":
		// "interestingly equal", i.e. case-insensitive comparison
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.Boolean{Value: strings.EqualFold(leftVal, rightVal)}
	default:
		return newError("%s (unknown operator: %s %s %s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), left.Type(), operator, right.Type())
	}

}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "AND":
		return &object.Boolean{Value: leftVal && rightVal}
	case "OR":
		return &object.Boolean{Value: leftVal || rightVal}
	case "XOR":
		// Ummm.... https://stackoverflow.com/questions/23025694/is-there-no-xor-operator-for-booleans-in-golang
		return &object.Boolean{Value: leftVal != rightVal}
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
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMERIC_OBJ {
		return newError("%s (unknown operator: -%s)", syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound), right.Type())
	}
	value := right.(*object.Numeric).Value
	return &object.Numeric{Value: -value}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
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
