# bool expression parser and evaluator

## supported grammars

- `true`, `false`, `number` (float64 and int64), `string`, `identifier`
- compare operations: `=` `>` `>=` `<` `<=` `!=` 
- simple math operations: `+` `-` `*` `/`
- prefix: `-`, `!`
- bracket `()`
- custom variable, like `$car`
- custom none-variadic functions, like : `takeBus()`
- in, like `$car in ('bwm','byd')`

##  grammars limits
- function inputs only supported `int64`, `float64`, `string`, `bool`
- functions must have one or two returns, the first one must be one of: `int64`, `float64`, `string`, `bool`, if the second one exists, it must be error
- variables support numbers, including all ints and uints, except `uint64`
- function does not support variadic args
- if `identifier` exists, then the expression can only be parsed, but cannot be evaluated


## builtin functions:

- `contains`   
- `endsWith`   
- `startsWith`   
- `length`   
- `toLower`   
- `toUpper`   
- `trim`   
- `concat`   
- `geoWithin2d`   
- `hasIntersection`   
- `timestampBefore`   
- `now`   // return time.Now() rfc3339 formatted string

## examples

### simple one

```go
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
```

## full one

```go

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
```