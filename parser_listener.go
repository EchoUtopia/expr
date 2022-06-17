package expr

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/EchoUtopia/expr/parser"
	"github.com/EchoUtopia/zerror"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type stack struct {
	stack []reflect.Value
}

func (l *stack) push(i reflect.Value) {
	l.stack = append(l.stack, i)
}

func (l *stack) pop() reflect.Value {
	if len(l.stack) < 1 {
		panic("stack is empty unable to pop")
	}
	result := l.stack[len(l.stack)-1]
	l.stack = l.stack[:len(l.stack)-1]
	return result
}

type ErrorListener struct {
	err error
	*antlr.DefaultErrorListener
	walker *ParseTreeWalker
}

func (el *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	el.err = errors.New(fmt.Sprintln("line " + strconv.Itoa(line) + ":" + strconv.Itoa(column) + " " + msg))
}

func (el *ErrorListener) setError(err error) {
	if el.err == nil {
		el.err = err
		el.walker.Stop()
	}
}

type argType struct{}

var argForParse = argType{}
var reflectArgVal = reflect.ValueOf(argForParse)
var reflectArgType = reflect.TypeOf(argForParse)

var boolRVal = reflect.ValueOf(true)
var intRVal = reflect.ValueOf(int64(1))
var floatRVal = reflect.ValueOf(float64(2.2))
var stringRVal = reflect.ValueOf(`str`)

type listenerForParse struct {
	*ErrorListener
	*parser.BaseExprListener
	funcs  map[string]*function
	walker *ParseTreeWalker
	*stack
	// TODO: check if one arg used as many types
	argFirstTypes map[string]reflect.Kind
}

func NewParser() *listenerForParse {
	walker := &ParseTreeWalker{}
	parser := &listenerForParse{
		ErrorListener:    &ErrorListener{walker: walker},
		BaseExprListener: &parser.BaseExprListener{},
		funcs:            map[string]*function{},
		walker:           walker,
		stack:            &stack{},
		//argFirstTypes:    map[string]reflect.Kind{},
	}

	return parser
}

// also used by evaluator
func (l *listenerForParse) ExitNumber(c *parser.NumberContext) {
	var val interface{}
	text := c.GetText()
	if c.FLOAT() != nil {
		f, err := strconv.ParseFloat(text, 64)
		if err != nil {
			panic(err)
		}
		val = f
	} else {
		i, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			panic(err)
		}
		val = i
	}
	l.push(reflect.ValueOf(val))
}

// also used by evaluator
func (l *listenerForParse) ExitString(c *parser.StringContext) {
	l.push(convertTextToVal(c.GetText()))
}

// also used by evaluator
func (l *listenerForParse) ExitStringList(c *parser.StringListContext) {
	ls := []string{}
	for _, v := range c.AllSTRING() {
		ls = append(ls, convertText(v.GetText()))
	}
	l.push(reflect.ValueOf(ls))
}

// also used by evaluator
func (l *listenerForParse) ExitNumberList(c *parser.NumberListContext) {
	useFloat := false
	list := []reflect.Value{}
	for i := 0; i < len(c.AllNumber()); i++ {
		vn := l.pop()
		list = append(list, vn)
		if vn.Kind() == reflect.Float64 {
			useFloat = true
		}
	}
	if useFloat {
		fl := make([]float64, 0, len(list))
		for _, v := range list {
			if v.Kind() == reflect.Float64 {
				fl = append(fl, v.Interface().(float64))
			} else {
				fl = append(fl, float64(v.Interface().(int64)))
			}
		}
		l.push(reflect.ValueOf(fl))
	} else {
		il := make([]int64, 0, len(list))
		for _, v := range list {
			if v.Kind() == reflect.Int64 {
				il = append(il, v.Interface().(int64))
			} else {
				il = append(il, int64(v.Interface().(float64)))
			}
		}
		l.push(reflect.ValueOf(il))
	}
}

// also used by evaluator
func (l *listenerForParse) ExitBoolean(c *parser.BooleanContext) {
	txt := c.GetText()
	var val bool
	if txt == `true` {
		val = true
	} else {
		val = false
	}
	l.push(reflect.ValueOf(val))
}

// check for function returns
func checkParseMathOperand(args ...reflect.Value) error {
	for _, n := range args {
		k := n.Kind()
		if k != reflect.Float64 && k != reflect.Int64 && n.Type() != reflectArgType {
			return zerror.BadRequest.Errorf(`can not do math operation with type: (%s)%v`, n.Kind(), n.Interface())
		}
	}
	return nil
}

func (l *listenerForParse) reset() {
	l.stack.stack = l.stack.stack[:0]
	l.err = nil
}

func (l *listenerForParse) ExitIdentifier(c *parser.IdentifierContext) {
	l.push(reflectArgVal)
}

func (l *listenerForParse) ExitBoolIdentifier(c *parser.BoolIdentifierContext) {
	l.push(boolRVal)
}

// ExitVariable is called when exiting the Variable production.
func (l *listenerForParse) ExitVariable(c *parser.VariableContext) {
	l.push(reflectArgVal)
}

func (l *listenerForParse) ExitBoolVariable(c *parser.BoolVariableContext) {
	l.push(boolRVal)
}

func (l *listenerForParse) ExitCompare(c *parser.CompareContext) {
	rv, lv := l.pop(), l.pop()
	rk, lk := rv.Kind(), lv.Kind()
	if (rk == reflect.String || lk == reflect.String) && rk != lk && rv.Type() != reflectArgType && lv.Type() != reflectArgType {
		l.setError(zerror.BadRequest.Errorf(`can not compare between (%s)%s and (%s)%v`, rk, rv, lk, lv))
		return
	}
	l.push(boolRVal)
}

func (l *listenerForParse) ExitMulDiv(c *parser.MulDivContext) {
	l.exitMathOperand()
}

func (l *listenerForParse) exitMathOperand() {
	rv, lv := l.pop(), l.pop()
	if err := checkParseMathOperand(rv, lv); err != nil {
		l.setError(err)
	}
	useFloat := false
	if rv.Kind() == reflect.Float64 || lv.Kind() == reflect.Float64 {
		useFloat = true
	}
	if useFloat {
		l.push(floatRVal)
	} else if lv.Type() == reflectArgType && rv.Type() == reflectArgType {
		l.push(lv)
	} else {
		l.push(intRVal)
	}
}

func (l *listenerForParse) ExitAddSub(c *parser.AddSubContext) {
	l.exitMathOperand()
}

func convertTextToVal(s string) reflect.Value {
	s = convertText(s)
	return reflect.ValueOf(s)
}

func (l *listenerForParse) ExitSubExpression(c *parser.SubExpressionContext) {
	v := l.pop()
	if v.Kind() == reflect.String {
		l.setError(zerror.BadRequest.Errorf(`invalid expression: %s`, c.GetText()))
		return
	}
	l.push(v)
}

func convertText(s string) string {
	s = s[1 : len(s)-1]
	s = strings.ReplaceAll(s, `\'`, `'`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	return s
}

func (l *listenerForParse) RegisterFunc(name string, fn interface{}) error {
	_, ok := l.funcs[name]
	if ok {
		return zerror.AlreadyExists.WithMsg(name)
	}
	returnType, fnVal, err := checkFunction(name, fn)
	if err != nil {
		return err
	}
	fnType := fnVal.Type()
	l.funcs[name] = &function{
		fn:         fnVal,
		returnType: returnType,
		argsNumber: fnType.NumIn(),
		isVariadic: fnType.IsVariadic(),
		name:       name,
	}
	return nil
}

func (l *listenerForParse) getFunc(name string) (*function, error) {
	globalFn, ok := globalFuncs.get(name)
	if ok {
		return globalFn, nil
	}
	fn, ok := l.funcs[name]
	if !ok {
		return nil, zerror.BadRequest.Errorf(`func: %s not found`, name)
	}
	return fn, nil
}

func (l *listenerForParse) checkFuncCallReturnTrue(name string, returnBool bool) error {
	fn, err := l.getFunc(name)
	if err != nil {
		return err
	}
	if fn.returnType.Kind() == reflect.Bool != returnBool {
		not := ``
		if !returnBool {
			not = `does not`
		}
		return zerror.BadRequest.Errorf(`expect func: %s %s return bool `, name, not)
	}
	return nil
}

func (l *listenerForParse) getCheckFuncCall(name string, hasArgs bool) (*function, error) {
	fn, err := l.getFunc(name)
	if err != nil {
		return nil, err
	}
	var gotArgNumber int
	if hasArgs {
		gotArgNumber = int(l.pop().Int())
	}
	if !fn.isVariadic && gotArgNumber != fn.argsNumber || fn.isVariadic && gotArgNumber < fn.argsNumber-1 {
		return nil, zerror.BadRequest.Errorf(`func: %s expect %d args, %d got`, name, fn.argsNumber, gotArgNumber)
	}
	for i := gotArgNumber - 1; i >= 0; i-- {
		arg := l.pop()
		if arg.Type() != reflectArgType && arg.Type() != fn.fn.Type().In(i) {
			return nil, zerror.BadRequest.Errorf(`func: %s, arg position: %d expect %s, got %s`, name, i, fn.fn.Type().In(i), arg.Type())
		}
	}
	// TODO: check if got args types match functions args types
	return fn, nil
}

func (l *listenerForParse) ExitExpressionFunction(c *parser.ExpressionFunctionContext) {
	fnName := c.Function().GetName().GetText()
	if err := l.checkFuncCallReturnTrue(fnName, false); err != nil {
		l.setError(err)
	}
}

func (l *listenerForParse) ExitBoolFunction(c *parser.BoolFunctionContext) {
	fnName := c.Function().GetName().GetText()
	if err := l.checkFuncCallReturnTrue(fnName, true); err != nil {
		l.setError(err)
	}
}

func (l *listenerForParse) ExitFunction(c *parser.FunctionContext) {
	fnName := c.GetName().GetText()
	fn, err := l.getCheckFuncCall(fnName, c.GetFnargs() != nil)
	if err != nil {
		l.setError(err)
		return
	}
	l.push(reflect.New(fn.returnType).Elem())
}

// ExitArg is called when exiting the arg production.
func (l *listenerForParse) ExitArg(c *parser.ArgContext) {
	if c.VAR() != nil {
		l.push(reflectArgVal)
	} else if c.BOOLEAN() != nil {
		l.push(boolRVal)
	} else if c.STRING() != nil {
		l.push(stringRVal)
	}
}

func (l *listenerForParse) ExitArgs(c *parser.ArgsContext) {
	argNum := len(c.AllArg())
	l.push(reflect.ValueOf(argNum))
}

func (l *listenerForParse) ExitIn(c *parser.InContext) {
	rv, lv := l.pop(), l.pop()
	lk := lv.Kind()
	_, isStringList := rv.Interface().([]string)
	if lv.Type() != reflectArgType &&
		(isStringList && lv.Kind() != reflect.String || !isStringList && lv.Kind() == reflect.String) {
		l.setError(zerror.BadRequest.Errorf(`invalid expression: %s, left operand type: %s`, c.GetText(), lk))
	} else {
		l.push(boolRVal)
	}

}

func (l *listenerForParse) ExitEveryRule(ctx antlr.ParserRuleContext) {
	if l.err != nil {
		l.walker.Stop()
	}
}
