package expr

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"reflect"
	"testing"
)

func parse(withCache bool) bool {
	input := `(3+2)*2 >= 1`
	//input = `'s' in ('a')`
	parser := NewParser()
	var tree antlr.Tree
	var err error
	if withCache {
		tree, err = parser.ParseWithCache(input)
	} else {
		tree, err = parser.Parse(input)
	}
	if err != nil {
		panic(err)
	}
	evaluator, _ := NewEvaluatorWithParser(tree, parser, nil)
	result, err := evaluator.Evaluate()
	if err != nil {
		panic(err)
	}
	return result
}

func BenchmarkEvaluator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parse(false)
	}
}

func BenchmarkEvaluateWithCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parse(true)
	}
}

func BenchmarkReflectValue(b *testing.B) {
	a := 1
	for i := 0; i < b.N; i++ {
		va := reflect.ValueOf(a)
		_ = va.Int()
	}
}
