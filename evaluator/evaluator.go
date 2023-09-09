package evaluator

import (
	//	"math"
	"fmt"
	"math"

	"github.com/hellracer2007/webCalc/calculator/ast"
	"github.com/hellracer2007/webCalc/calculator/object"
)

func Eval(node ast.Node) object.Object {	
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement :
		return Eval(node.Expression)
	case *ast.Procedure :
		body := Eval(node.Body)
		return evalProcedure(node.Func ,body)
	case *ast.PostfixExpression:
		left := Eval(node.Left)
		return evalPostFixExpression(node.Token.Literal, left)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Token.Literal, right)
	case *ast.InfixExpression:
		left:= Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		l, r := normalizeExpr(left, right)
		return evalInfixExpression(node.Token.Literal, l, r)
	}
	return nil
} 

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch v := right.(type){
	case *object.Integer :
		v.Value = -v.Value
		return v
	case *object.Float :
		v.Value = -v.Value
		return v
	}
	return nil
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalInfixFloatExpression(operator, left, right)
	}
	return nil
}

func evalPostFixExpression(operator string, left object.Object) object.Object {
	switch {
	case operator == "!":
		return evalFactorial(left)	
	}
	return nil
}

func evalFactorial(left object.Object) object.Object{
	val :=  left.(*object.Integer).Value 
	result := val
	for i := val; i >= 2; i-- {
		result *= (i-1)
	}
	return &object.Integer{Value: result}
}

func evalInfixIntegerExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	
	var result object.Integer
	switch operator {
	case "+":
		result.Value = leftVal + rightVal
	case "-":
		result.Value = leftVal - rightVal
	case "*":
		result.Value = leftVal * rightVal
	case "/":
		result.Value = leftVal / rightVal
	case "E":
		res := solveExp(left, right)
		sol, ok := res.(*object.Integer)
		if ok {
			result.Value = sol.Value 
		} else {
			sol := res.(*object.Float)
			return &object.Float{Value: float64(sol.Value)}
		}
	case "^" :
		res := &object.Float{Value: math.Pow(float64(leftVal), float64(rightVal))}
		solution, ok := normalizeNumber(res).(*object.Integer)
		if ok {
			return solution
		}
		return res
	case "√":
		fmt.Println("we mafe it here")
		res := math.Pow(float64(rightVal), 1.0/float64(leftVal))
		return normalizeNumber(&object.Float{Value: res})
	}	
	return &result
}

func evalInfixFloatExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	
	var result object.Float
	switch operator {
	case "+":
		result.Value = leftVal + rightVal
	case "-":
		result.Value = leftVal - rightVal
	case "*":
		result.Value = leftVal * rightVal
	case "/":
		result.Value = leftVal / rightVal
	case "E":
		result.Value = solveExp(left, right).(*object.Float).Value
	case "^":
		result.Value = math.Pow(leftVal, rightVal)
	case "√":
		res := math.Pow(rightVal, 1.0/leftVal)
		return normalizeNumber(&object.Float{Value: res})
	}	

	return normalizeNumber(&result)
}

func solveExp(value, exp object.Object) object.Object {
	expo, ok := exp.(*object.Integer)
	if !ok {
		return newError("exp should be of type INT, instead got %d of type %s", expo.Value, expo.Type())
		
	}
	val, ok := value.(*object.Integer)
	if ok {
		value = &object.Float{Value: float64(val.Value)}
	}
	exp = &object.Float{Value: float64(expo.Value)}
	
	left := value.(*object.Float).Value
	right := exp.(*object.Float).Value
	result := &object.Float{Value: left}
	for right > 0 {
		result.Value *= 10
		right--
	}
	fmt.Println(result.Value)
	return normalizeNumber(result)

}

func normalizeExpr(left, right object.Object) (object.Object, object.Object) {
	fmt.Println("normalize")
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		right = &object.Float{Value: float64(right.(*object.Integer).Value)}
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		left = &object.Float{Value: float64(left.(*object.Integer).Value)}
	}

	return left, right
}



func evalProcedure(proc string, body object.Object) object.Object {
	var result object.Object
	switch {
	case proc == "sin":
		result = resolveDegreesProc(body, math.Sin)
	case proc == "tan":
		result = resolveDegreesProc(body, math.Tan)
	case proc == "cos":
		result = resolveDegreesProc(body, math.Cos)
	case proc == "log":
		result = resolveProc(body, math.Log10)
	case proc == "arcsin":
		result = resolveDegreesProc(body, math.Asin)
	case proc == "arccos":
		result = resolveDegreesProc(body, math.Acos)
	case proc == "arctan":
		result = resolveDegreesProc(body, math.Atan)
	case proc == "ln":
		result = resolveProc(body, math.Log)
	case proc == "√":
		result = resolveProc(body, math.Sqrt)
	}
	return result
}

func resolveProc(body object.Object, proc func(float64) float64) object.Object {
	vali, ok := body.(*object.Integer)
	var res object.Object
	if ok {
		result := proc(float64(vali.Value))
		res = &object.Float{Value: result}
	}
	valf, ok := body.(*object.Float)
	if ok {
		result := proc(valf.Value)
		res = &object.Float{Value: result}
	}
	fmt.Println(res)
	return normalizeNumber(res)
}

func resolveDegreesProc(body object.Object, proc func(float64) float64) object.Object{
	vali, ok := body.(*object.Integer)
	var res object.Object
	if ok {
	result := proc(float64(vali.Value)*math.Pi/180)
	res = &object.Float{Value: result}
	}
	valf, ok := body.(*object.Float)
	if ok {
	result := proc(valf.Value*math.Pi/180)
	res = &object.Float{Value: result}
	}

	return normalizeNumber(res)
}


func normalizeNumber(number object.Object)object.Object{
	result, ok := number.(*object.Float)
	if ok && result.Value == float64(int(result.Value)) {
		return &object.Integer{Value: int64(result.Value)}
	}
	return number
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a ...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
