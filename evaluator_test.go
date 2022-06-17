package expr

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBoolEvaluator(t *testing.T) {
	vars := map[string]interface{}{
		`true`:   true,
		`false`:  false,
		`int`:    1,
		`float`:  2.2,
		`string`: `str`,
	}
	trues := []string{
		`$true`,
		`$true and false or $true`,
		`$true or false`,
		`!!$true`,
		`!!true`,
		`'str' = $string`,
		`2 = 2.0`,
		`$int = 1.0`,
		`$string in ('1', 'str')`,
		`$string not in ('1', 'abc')`,
		`-t_int() = - $int`,
		`$string > 'a'`,
		`t_int() in (1.0, 2)`,
		`1 in (1.0)`,
		`'a' in ('a', 'b')`,
		`1 not in (2, 3)`,
		`t_bool()`,
		`!1 in (2, 3)`,
		`true != false`,
		`! true =false`,
	}

	falses := []string{
		`t_int() = - $int`,
		`$false and $true`,
		`true = false`,
	}
	for _, v := range trues {
		if v == `` {
			continue
		}
		testEvaluator(t, v, vars, true)
	}
	for _, v := range falses {
		testEvaluator(t, v, vars, false)
	}
}

func TestEvaluator(t *testing.T) {
	vars := map[string]interface{}{
		`true`:   true,
		`false`:  false,
		`int`:    1,
		`float`:  2.2,
		`string`: `str`,
	}
	trues := []string{
		`1+1 > 1`,
		`1+1.0 = 2`,
		`(3+2)*2 > 7`,
		`3+2 * 2 = 7`,
		`$float - $int = 1.2000000000000002`,
		`-t_int() = - $int`,
		`-(1+2)*3 = -9`,
		`-t_int() * -t_int() = t_int() * t_int()`,
		`2.1 * 2.1 > 4 and 2.1 * 2.1 < 5 and 2.1 * 2.1 = 4.41`,
		`-2.1 * 2.1 = -4.41`,
	}

	falses := []string{
		// float lose precision, eg 2.2-1 = 1.2000000000000002
		`$float - $int = 1.2`,
		`t_int() = - $int`,
		`2.1 * 2.1 = 4`,
	}
	for _, v := range trues {
		testEvaluator(t, v, vars, true)
	}
	for _, v := range falses {
		testEvaluator(t, v, vars, false)
	}
}

func TestInEvaluate(t *testing.T) {
	vars := map[string]interface{}{
		`true`:   true,
		`false`:  false,
		`int`:    1,
		`float`:  2.2,
		`string`: `str`,
	}

	invalids := []string{
		`!$var in (1,2,3)`,
		`(a+b) in (1,2,3)`,
	}
	for _, v := range invalids {
		_, err := Evaluate(v, vars)
		require.NotNil(t, err, v)
	}
	//return
	trues := []string{
		`!$float in (1,2,3)`,
		`1 in (1,2,3)`,
		`1 in (1.0, -1)`,
		`$int in (1.0, -2)`,
		`t_int() in (1,2,3)`,
		`t_string() in ('a', 'b', 'test')`,
		`1 in (1)`,
		`(1 in (1))`,
		`-$float in (1,2,3, -2.2)`,
		`(1+2)* 3.0 in (1,2,3, 9)`,
		`1+2 in (1,2,3)`,
	}
	for _, v := range trues {
		testEvaluator(t, v, vars, true)
	}

	falses := []string{
		`t_string() in ('a', 'b')`,
	}

	for _, v := range falses {
		testEvaluator(t, v, vars, false)
	}
}

func TestFunctionEvaluator(t *testing.T) {
	vars := map[string]interface{}{
		`true`:    true,
		`false`:   false,
		`int`:     1,
		`float`:   2.2,
		`string`:  `str`,
		`coord2d`: `1,1`,
	}
	trues := []string{
		`testArgs($int, $float, true, $string)`,
		`testArgs($int, $float, $true, $string)`,
		`contains('abc', 'ab')`,
		`endsWith('abc', 'bc')`,
		`startsWith('abc', 'ab')`,
		`length('abc')=3`,
		`toLower('ABC')='abc'`,
		`toUpper($string)='STR'`,
	}

	for _, v := range trues {
		if v == `` {
			continue
		}
		testEvaluator(t, v, vars, true)
	}

	falses := []string{
		`contains('abc', 'ac')`,
		`endsWith('abc', 'ab')`,
		`startsWith('abc', 'bc')`,
	}

	for _, v := range falses {
		if v == `` {
			continue
		}
		testEvaluator(t, v, vars, false)
	}

	invalids := []string{
		`testArgs($int, $int, $int, $int)`,
		`testArgs($int, $int, $xxx, $int)`,
	}
	for _, expr := range invalids {
		parser := NewParser()
		tree, err := parser.Parse(expr)
		require.Equalf(t, nil, err, expr)
		evaluator, err := NewEvaluatorWithParser(tree, parser, vars)
		require.Equalf(t, nil, err, expr)
		_, err = evaluator.Evaluate()
		require.NotEqualf(t, nil, err, expr)
	}
}

func TestFuncError(t *testing.T) {
	parser := NewParser()
	pn := func() bool { panic(`panic`) }
	if err := parser.RegisterFunc(`panic`, pn); err != nil {
		t.Fatal(err)
	}
	returnErr := func() (bool, error) { return false, errors.New(`error`) }
	if err := parser.RegisterFunc(`error`, returnErr); err != nil {
		t.Fatal(err)
	}
	exprs := []string{
		`panic()`,
		`error()`,
	}
	for _, v := range exprs {
		tree, err := parser.Parse(v)
		if err != nil {
			t.Fatal(err)
		}
		ev, err := NewEvaluatorWithParser(tree, parser, nil)
		if err != nil {
			t.Fatal(err)
		}
		result, err := ev.Evaluate()
		require.Equal(t, false, result, v)
		require.NotNil(t, err, v)
	}
}

func testEvaluator(t *testing.T, expr string, vars map[string]interface{}, expect bool) {
	result, err := Evaluate(expr, vars)
	require.Equal(t, nil, err, expr, err)
	require.Equal(t, expect, result, expr)
}

func TestCheckVariables(t *testing.T) {
	invalids := map[string]interface{}{
		``:   ``,
		`t1`: uint64(1),
		`t2`: struct{}{},
		`t3`: time.Now(),
	}
	for k, v := range invalids {
		require.NotEqual(t, nil, checkSetVariables(map[string]interface{}{
			k: v,
		}))
	}
}

func ExampleEvaluate() {
	expr := `$car in ('bwm', 'byd') and (3 + 2) * 2.0 = 10 and startsWith($car, 'b')`
	variables := map[string]interface{}{`car`: `byd`}
	result, err := Evaluate(expr, variables)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	// true
}

func BenchmarkEvaluate(b *testing.B) {
	expr := `$car in ('bwm', 'byd') and (3 + 2) * 2.0 = 10 and startsWith($car, 'b')`
	variables := map[string]interface{}{`car`: `byd`}
	for i := 0; i < b.N; i++ {
		_, err := Evaluate(expr, variables)
		if err != nil {
			panic(err)
		}

	}

}

func ExampleParseAndEvaluate() {
	expr := `$car in ('bwm', 'byd') and (3 + echo_int(2)) * 2.0 = 10 and startsWith($car, 'b')`
	var tInt = func(a int64) int64 {
		return a
	}
	if err := RegisterFunc(`echo_int`, tInt); err != nil {
		panic(err)
	}

	variables := map[string]interface{}{`car`: `byd`}
	result, err := Evaluate(expr, variables)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	// true
}
