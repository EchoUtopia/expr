package expr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testBool() bool {
	return true
}

func testInt() int64 {
	return 1
}

func testFloat() float64 {
	return 2.2
}

func testString() string {
	return `test`
}

func TestMain(t *testing.M) {
	RegisterFunc(`t_bool`, testBool)
	RegisterFunc(`t_int`, testInt)
	RegisterFunc(`t_float`, testFloat)
	RegisterFunc(`t_string`, testString)
	RegisterFunc(`testArgs`, testArgs)
	t.Run()
}

func TestParseBoolExpression(t *testing.T) {
	invalids := []string{
		`1`,
		`'a'`,
		`1='a'`,
		`t_float()`,
		`t_int()`,
		`t_string()`,
		`t_bool() > 1`,
		`t_float() > t_string()`,
		`'true'`,
		`'false'`,
		`t_string()= 1`,
		`!t_string()`,
		`-$var`,
		`1 in ('a', 1)`,
		`'a' in ('a', 1)`,
		`1 in ('a', 'a')`,
		`'a' in (1, 1)`,
		`'a' = {a}`,
	}
	for _, v := range invalids {
		p := NewParser()
		_, err := p.Parse(v)
		require.NotEqualf(t, nil, err, v)
	}
	valids := []string{
		`$var`,
		`!$_var`,
		`t_bool()`,
		`t_float() > t_float()`,
		`t_float() > t_int()`,
		`a`, // identifier for special use
		`t_string() = 'a'`,
		`t_string() = 'a'`,
		`true and (false or t_bool())`,
		`!t_bool()`,
		`true`,
		//`1 in `
	}
	for _, v := range valids {
		p := NewParser()
		_, err := p.Parse(v)
		require.Equalf(t, nil, err, v)
	}
}

func TestParseExpression(t *testing.T) {
	invalids := []string{
		`-'a' > 1`,
		`1 - 'a' > 1`,
		`1 - -'a' > 1`,
		`-'a' > 1`,
		`-true`,
		`-t_bool() > 1`,
		`-t_string() > 1`,
		`-`,
	}
	for _, v := range invalids {
		p := NewParser()
		_, err := p.Parse(v)
		require.NotEqualf(t, nil, err, v)
	}
	//return
	valids := []string{
		`-t_float() > 1`,
		`-t_float() > -t_int()`,
		`--1 > 1`,
		`1+-1> 2`,
		`(3--1)*-3 > 1`,
		`1+a > 1`,
		`1+(-1)> 2`,
	}
	for _, v := range valids {
		p := NewParser()
		_, err := p.Parse(v)
		require.Equalf(t, nil, err, v)
	}
}

func TestParseCompare(t *testing.T) {
	invalids := []string{
		`1>'a'`,
		`'a'=1`,
		`true > false`,
		`true > 1`,
		`1 > true`,
		`t_float() > t_string()`,
		`1 = 'a'`,
		`(1+5)*3 > 'a'`,
		`(1+a)*3 > 'a'`,
	}
	for _, v := range invalids {
		p := NewParser()
		_, err := p.Parse(v)
		require.NotEqualf(t, nil, err, v)
	}
	//return
	valids := []string{
		`$a > 1`,
		`$a = 'a'`,
		`$a < 1.0`,
		`$a >= $b`,
		`$a >= -$b`,
		`(1+2)*3 > 4`,
		`t_string() > t_string()`,
		`(1+2)*3 > (4)`,
		`1 = 1.0`,
	}
	for _, v := range valids {
		p := NewParser()
		_, err := p.Parse(v)
		require.Nil(t, err, v)
	}
}

func TestParseIn(t *testing.T) {
	invalids := []string{
		`1 in ('a')`,
		`'a' in (1,2,3)`,
		`t_string() in (1,2,3)`,
		`t_int() in ('a','b')`,
		`1 in (1, 'a')`,
		`1+2 in ('a', 'b')`,
		`false in (1,2)`,
		`false in (true, false)`,
	}
	for _, v := range invalids {
		p := NewParser()
		_, err := p.Parse(v)
		require.NotNil(t, err, v)
	}
	//return
	valids := []string{
		// TODO: this expr should be invalid in parser
		`!$var in (1,2,3)`,
		`1 in (1,2,3)`,
		`1 in (1.0, -1)`,
		`$var in (1.0, -2)`,
		`t_int() in (1,2,3)`,
		`t_string() in ('a', 'b')`,
		`1 in (1)`,
		`(1 in (1))`,
		`-$var in (1,2,3)`,
		`(a+b) in (1,2,3)`,
		`(1+2)* 3.0 in (1,2,3)`,
		`t_string() in ('a', 'b')`,
		`1+2 in (1,2,3)`,
	}
	for _, v := range valids {
		p := NewParser()
		_, err := p.Parse(v)
		require.Equalf(t, nil, err, v)
	}
}

func testArgs(a int64, b float64, c bool, d string) bool {
	return true
}

func TestParseFunctionArgs(t *testing.T) {
	invalids := []string{
		`testArgs()`,
		`testArgs(1, 1)`,
		`testArgs(1, 1,1,1)`,
		`testArgs(1, 1.0,true,'a', 1)`,
		`testArgs(1.0, 1,true,'a')`,
		`testArgs(-$a,$a,$a,$a)`,
		`testArgs(!$a,$a,$a,$a)`,
	}
	for _, v := range invalids {
		p := NewParser()
		_, err := p.Parse(v)
		require.NotEqualf(t, nil, err, v)
	}
	//return
	valids := []string{
		`testArgs(1, 1.0,true,'a')`,
		`testArgs(1, 1.0,true,'a')`,
		`testArgs($a,$a,$a,$a)`,
	}
	for _, v := range valids {
		p := NewParser()
		_, err := p.Parse(v)
		require.Equalf(t, nil, err, v)
	}
}

func TestCheckFunction(t *testing.T) {
	invalids := map[string]interface{}{
		`k1`: 1,
		`k2`: func(...interface{}) {},
		`k3`: func(string, int) {},
		`k4`: func() (string, int, error) { return ``, 0, nil },
		`k5`: func() (error, int) { return nil, 0 },
		`k6`: func() int { return 0 },
		`k7`: func() (int, error) { return 0, nil },
		`k8`: func() error { return nil },
	}
	for k, v := range invalids {
		_, _, err := checkFunction(k, v)
		require.NotEqual(t, nil, err, k, v)
	}
}
