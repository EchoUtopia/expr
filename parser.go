package expr

import (
	"sync"

	"github.com/EchoUtopia/expr/parser"
	"github.com/EchoUtopia/zerror"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func init() {

	for k, v := range builtinFuncs {
		if err := RegisterFunc(k, v); err != nil {
			panic(err)
		}
	}
}

type parsedPools struct {
	sync.Mutex
	trees map[string]antlr.Tree
}

var pools = &parsedPools{trees: map[string]antlr.Tree{}}

type funcPools struct {
	sync.RWMutex
	funcs map[string]*function
}

var globalFuncs = funcPools{funcs: map[string]*function{}}

func (fp *funcPools) get(name string) (*function, bool) {
	fp.RLock()
	defer fp.RUnlock()
	fn, ok := globalFuncs.funcs[name]
	return fn, ok
}

func (fp *funcPools) set(name string, fn *function) {
	fp.Lock()
	defer fp.Unlock()
	globalFuncs.funcs[name] = fn
}

// Parse can check some errors before evaluate stage
func Parse(input string) (tree antlr.Tree, err error) {
	parser := NewParser()
	//tree, err := parser.Parse(expr)
	tree, err = parser.ParseWithCache(input)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func RegisterFunc(name string, fn interface{}) error {

	_, ok := globalFuncs.get(name)
	if ok {
		return zerror.AlreadyExists.WithMsg(name)
	}
	returnType, fnVal, err := checkFunction(name, fn)
	if err != nil {
		return err
	}
	fnType := fnVal.Type()
	f := &function{
		fn:         fnVal,
		returnType: returnType,
		argsNumber: fnType.NumIn(),
		isVariadic: fnType.IsVariadic(),
		name:       name,
	}
	globalFuncs.set(name, f)
	return nil
}

func (l *listenerForParse) Parse(input string) (tree antlr.Tree, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = zerror.Internal.Errorf(`panic: %v`, recovered)
		}
	}()
	l.reset()
	is := antlr.NewInputStream(input)

	// Create the Lexer
	lexer := parser.NewExprLexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(l.ErrorListener)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	// Create the Parser
	p := parser.NewExprParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(l.ErrorListener)
	tree = p.Start()
	l.walker.Walk(l, tree)
	err = l.err
	return
}

// for performance improvement
func (l *listenerForParse) ParseWithCache(input string) (antlr.Tree, error) {
	pools.Lock()
	defer pools.Unlock()
	tree, ok := pools.trees[input]
	if ok {
		return tree, nil
	}
	tree, err := l.Parse(input)
	if err != nil {
		return nil, err
	}
	pools.trees[input] = tree
	return tree, nil
}
