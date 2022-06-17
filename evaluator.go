package expr

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/EchoUtopia/expr/parser"
	"github.com/EchoUtopia/zerror"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type evaluator struct {
	*listenerForParse
	tree      antlr.Tree
	variables map[string]interface{}
}

func NewEvaluatorWithParser(tree antlr.Tree, parser *listenerForParse, vars map[string]interface{}) (*evaluator, error) {
	if err := checkSetVariables(vars); err != nil {
		return nil, err
	}
	parser.stack.stack = parser.stack.stack[:0]
	return &evaluator{
		listenerForParse: parser,
		tree:             tree,
		variables:        vars,
	}, nil
}

// Evaluate evaluate expr with custom variables
func Evaluate(expr string, vars map[string]interface{}) (bool, error) {
	parser := NewParser()
	tree, err := parser.ParseWithCache(expr)
	if err != nil {
		return false, err
	}
	ev, err := NewEvaluatorWithParser(tree, parser, vars)
	if err != nil {
		return false, err
	}
	result, err := ev.Evaluate()
	return result, err
}

func checkSetVariables(vars map[string]interface{}) error {
	for k, v := range vars {
		if k == `` {
			return zerror.BadRequest.WithMsg(`empty function name found`)
		}
		vVal := reflect.ValueOf(v)
		switch vVal.Kind() {
		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
			vars[k] = vVal.Int()
		case reflect.Uint32, reflect.Uint16, reflect.Uint8:
			vars[k] = int64(vVal.Uint())
		case reflect.Float64, reflect.Float32:
			vars[k] = vVal.Float()
		case reflect.Bool, reflect.String:
		default:
			return zerror.BadRequest.Errorf(`variable name: %s, type: %s not supported`, k, vVal.Kind())
		}
	}
	return nil
}

func (e *evaluator) Evaluate() (result bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = zerror.Internal.Errorf(`panic: %v`, recovered)
		}
	}()
	e.walker.Walk(e, e.tree)
	if e.err != nil {
		err = e.err
		return
	}
	result = e.pop().Bool()
	return
}

// ExitNot is called when exiting the Not production.
func (e *evaluator) ExitNot(c *parser.NotContext) {
	val := e.pop().Bool()
	val = !val
	e.push(reflect.ValueOf(val))
}

// ExitOr is called when exiting the Or production.
func (e *evaluator) ExitOr(c *parser.OrContext) {
	right, left := e.pop().Bool(), e.pop().Bool()
	val := left || right
	e.push(reflect.ValueOf(val))
}

func (e *evaluator) callFunc(funcName string, hasArgs bool) (reflect.Value, error) {
	var args []reflect.Value
	if hasArgs {
		args = e.pop().Interface().([]reflect.Value)
	}
	fn, err := e.getFunc(funcName)
	if err != nil {
		return reflect.Value{}, err
	}
	if hasArgs {
		for i := 0; i < fn.argsNumber; i++ {
			expected := fn.fn.Type().In(i)
			got := args[i].Type()
			if expected != got {
				return reflect.Value{}, zerror.BadRequest.Errorf(`func: %s, arg position: %d, expect: %s, got: %s`, funcName, i, expected, got)
			}
		}
	}

	outs := fn.fn.Call(args)
	result := outs[0]
	if len(outs) == 2 {
		out1 := outs[1].Interface()
		if out1 != nil {
			err = out1.(error)
		}
	}
	return result, err
}

func (e *evaluator) ExitIn(c *parser.InContext) {
	rv, lv := e.pop(), e.pop()
	if lv.Kind() == reflect.Bool {
		e.setError(zerror.BadRequest.Errorf(`invalid in operand: (%s)%s`, lv.Kind(), lv.Interface()))
		return
	}
	_, isStringList := rv.Interface().([]string)
	_, isFloatList := rv.Interface().([]float64)
	in := false
	if isStringList {
		right, left := rv.Interface().([]string), lv.String()
		for _, v := range right {
			if left == v {
				in = true
				break
			}
		}
	} else if isFloatList {
		left := lv.Interface()
		if i64, ok := left.(int64); ok {
			left = float64(i64)
		}
		right := rv.Interface().([]float64)
		for _, v := range right {
			if left == v {
				in = true
				break
			}
		}
	} else {
		iLeft := lv.Interface()
		useFloat := false
		if _, ok := iLeft.(float64); ok {
			useFloat = true
		}
		rights := rv.Interface().([]int64)
		var l, r interface{}
		for _, v := range rights {
			l = iLeft
			if useFloat {
				r = float64(v)
			} else {
				r = v
			}
			if l == r {
				in = true
				break
			}
		}
	}

	if c.GetOp().GetTokenType() == parser.ExprParserNOTIN {
		in = !in
	}
	e.push(reflect.ValueOf(in))
}

// ExitAnd is called when exiting the And production.
func (e *evaluator) ExitAnd(c *parser.AndContext) {
	right, left := e.pop().Bool(), e.pop().Bool()
	val := left && right
	e.push(reflect.ValueOf(val))
}

// ExitIdentifier is called when exiting the Identifier production.
func (e *evaluator) ExitIdentifier(c *parser.IdentifierContext) {
	e.setError(fmt.Errorf(`can not evaluate expression consists of identifier`))
}

func (e *evaluator) ExitVariable(c *parser.VariableContext) {
	val, err := e.getAndCheckVar(c.GetText(), false)
	if err != nil {
		e.setError(err)
		return
	}
	e.push(val)
}

func (e *evaluator) ExitBoolVariable(c *parser.BoolVariableContext) {
	val, err := e.getAndCheckVar(c.GetText(), true)
	if err != nil {
		e.setError(err)
		return
	}
	e.push(val)
}

// push args number
func (e *evaluator) ExitArgs(c *parser.ArgsContext) {
	argNum := len(c.AllArg())
	args := make([]reflect.Value, argNum)
	for i := 1; i <= argNum; i++ {
		ni := e.pop()
		args[argNum-i] = ni
	}
	e.push(reflect.ValueOf(args))
}

// ExitArg is called when exiting the arg production.
func (e *evaluator) ExitArg(c *parser.ArgContext) {
	var val interface{}
	if c.VAR() != nil {
		val, err := e.getVar(c.GetText())
		if err != nil {
			e.setError(err)
			return
		}
		e.push(val)
		return
	} else if c.BOOLEAN() != nil {
		if c.GetText() == `true` {
			val = true
		} else {
			val = false
		}
	} else if c.STRING() != nil {
		s := c.GetText()
		val := convertTextToVal(s)
		e.push(val)
		return
	}
	if val != nil {
		e.push(reflect.ValueOf(val))
	}
}

// ExitAddSub is called when exiting the AddSub production.
func (e *evaluator) ExitAddSub(c *parser.AddSubContext) {
	rv, lv := e.pop(), e.pop()
	if err := checkMathOperand(rv, lv); err != nil {
		e.setError(err)
	}
	var val interface{}
	isFloat, li, ri := cleanOperands(lv, rv)
	switch c.GetOp().GetTokenType() {
	case parser.ExprParserADD:
		if isFloat {
			val = li.(float64) + ri.(float64)
		} else {
			val = li.(int64) + ri.(int64)
		}
	case parser.ExprParserSUB:
		if isFloat {
			val = li.(float64) - ri.(float64)
		} else {
			val = li.(int64) - ri.(int64)
		}
	}
	e.push(reflect.ValueOf(val))
}

func (e *evaluator) getVar(name string) (reflect.Value, error) {
	name = name[1:]
	val, ok := e.variables[name]
	if !ok {
		return reflect.Value{}, zerror.BadRequest.Errorf(`var: %s not found`, name)
	}
	//rVal := reflect.ValueOf(val)
	//kind := rVal.Kind()
	//if operatorToken != nil {
	//	tokenTxt := operatorToken.GetText()
	//	if tokenTxt == `!` {
	//		if kind != reflect.Bool {
	//			return reflect.Value{}, zerror.BadRequest.Errorf(`expect $%s bool, but got %s`, name, kind)
	//		}
	//		rVal.SetBool(!rVal.Bool())
	//	} else {
	//		if kind != reflect.Int64 && kind != reflect.Float64 {
	//			return reflect.Value{}, zerror.BadRequest.Errorf(`expect $%s float64 or int64, but got %s`, name, kind)
	//		}
	//		if kind == reflect.Float64 {
	//			rVal.SetFloat(-rVal.Float())
	//		} else {
	//			rVal.SetInt(-rVal.Int())
	//		}
	//	}
	//}
	return reflect.ValueOf(val), nil
}

func (e *evaluator) ExitSubExpression(c *parser.SubExpressionContext) {
	v := e.pop()
	var val interface{}
	if v.Kind() == reflect.Float64 {
		val = -v.Float()
	} else {
		val = -v.Int()
	}
	e.push(reflect.ValueOf(val))
}

func (e *evaluator) getAndCheckVar(name string, expectBool bool) (reflect.Value, error) {
	rVal, err := e.getVar(name)
	if err != nil {
		return reflect.Value{}, err
	}
	if (rVal.Kind() == reflect.Bool) != expectBool {
		not := ``
		if !expectBool {
			not = `not`
		}
		return rVal, zerror.BadRequest.Errorf(`%s expect %s bool, got: %s`, name, not, rVal.Kind())
	}

	return rVal, nil
}

// ExitMulDivMod is called when exiting the MulDivMod production.
func (e *evaluator) ExitMulDiv(c *parser.MulDivContext) {

	rv, lv := e.pop(), e.pop()
	if err := checkMathOperand(rv, lv); err != nil {
		e.setError(err)
	}

	var val interface{}
	isFloat, li, ri := cleanOperands(lv, rv)
	switch c.GetOp().GetTokenType() {
	case parser.ExprParserMUL:
		if isFloat {
			val = li.(float64) * ri.(float64)
		} else {
			val = li.(int64) * ri.(int64)
		}
	case parser.ExprParserDIV:
		if isFloat {
			val = li.(float64) / ri.(float64)
		} else {
			val = li.(int64) / ri.(int64)
		}
	}
	e.push(reflect.ValueOf(val))
}

func (e *evaluator) ExitFunction(c *parser.FunctionContext) {
	result, err := e.callFunc(c.GetName().GetText(), c.GetFnargs() != nil)
	if err != nil {
		e.setError(err)
		return
	}
	e.push(result)
}

func (e *evaluator) ExitBoolIdentifier(c *parser.BoolIdentifierContext) {
	e.setError(fmt.Errorf(`can not evaluate expression consists of identifier`))
}

// ExitBoolCompare is called when exiting the BoolCompare production.
func (e *evaluator) ExitBoolCompare(c *parser.BoolCompareContext) {
	right, left := e.pop().Bool(), e.pop().Bool()
	var val interface{}
	if c.EQ() != nil {
		val = left == right
	} else {
		val = left != right
	}
	e.push(reflect.ValueOf(val))
}

func cleanOperands(lv, rv reflect.Value) (isFloat bool, left, right interface{}) {
	var il, ir int64
	var fl, fr float64
	lk, rk := lv.Kind(), rv.Kind()
	if lk == reflect.Float64 || rk == reflect.Float64 {
		isFloat = true
		if lk == reflect.Float64 {
			fl = lv.Float()
		} else {
			fl = float64(lv.Int())
		}
		if rk == reflect.Float64 {
			fr = rv.Float()
		} else {
			fr = float64(rv.Int())
		}
	} else {
		il, ir = lv.Int(), rv.Int()
	}
	if isFloat {
		left, right = fl, fr
	} else {
		left, right = il, ir
	}
	return
}

// ExitCompare is called when exiting the Compare production.
func (e *evaluator) ExitCompare(c *parser.CompareContext) {
	rv, lv := e.pop(), e.pop()
	if !isComparable(lv.Kind(), rv.Kind()) {
		e.setError(zerror.BadRequest.Errorf(`(%s)%s and (%s)%s is not comparable`, lv.Kind(), lv.Interface(), rv.Kind(), rv.Interface()))
		return
	}

	var compareVal float64
	if lv.Kind() == reflect.String {
		compareVal = float64(strings.Compare(lv.String(), rv.String()))
	} else {
		isFloat, li, ri := cleanOperands(lv, rv)
		if isFloat {
			left, right := li.(float64), ri.(float64)
			compareVal = left - right
		} else {
			left, right := li.(int64), ri.(int64)
			compareVal = float64(left - right)
		}
	}
	var result bool
	switch c.GetOp().GetTokenType() {
	case parser.ExprParserEQ:
		result = compareVal == 0
	case parser.ExprParserNEQ:
		result = compareVal != 0
	case parser.ExprParserGT:
		result = compareVal > 0
	case parser.ExprParserGTE:
		result = compareVal >= 0
	case parser.ExprParserLT:
		result = compareVal < 0
	case parser.ExprParserLTE:
		result = compareVal <= 0
	}
	e.push(reflect.ValueOf(result))
}

// only for string, int64, float64
func isComparable(a, b reflect.Kind) bool {
	stringKind := reflect.String
	if a == stringKind && b == stringKind {
		return true
	}
	if a == stringKind || b == stringKind {
		return false
	}
	return true
}

func checkMathOperand(args ...reflect.Value) error {
	for _, n := range args {
		k := n.Kind()
		if k != reflect.Float64 && k != reflect.Int64 {
			return zerror.BadRequest.Errorf(`can not do math operation with type: %v, value: %s`, n.Interface(), k)
		}
	}
	return nil
}
