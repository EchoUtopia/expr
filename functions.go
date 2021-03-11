package expr

import (
	"github.com/EchoUtopia/zerror"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type function struct {
	name string
	iFn  interface{}
	// the following will be set by register
	// if argsNumber is -1, no limit set
	argsNumber int
	isVariadic bool
	fn         reflect.Value
	returnType reflect.Type
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func checkFunction(name string, i interface{}) (reflect.Type, reflect.Value, error) {
	v := reflect.ValueOf(i)
	t := v.Type()
	if t.Kind() != reflect.Func {
		return nil, reflect.Value{}, zerror.BadRequest.Errorf(`[%s] is not func`, name)
	}
	// TODO: support variadic
	if t.IsVariadic() {
		return nil, reflect.Value{}, zerror.BadRequest.WithMsg(`variadic function not supported`)
	}
	for i := 0; i < t.NumIn(); i++ {
		nik := t.In(i).Kind()
		if nik != reflect.Int64 && nik != reflect.Float64 && nik != reflect.Bool && nik != reflect.String {
			return nil, reflect.Value{}, zerror.BadRequest.WithMsg(`function variables only support float64, int64, string and bool`)
		}
	}
	if t.NumOut() > 2 {
		return nil, reflect.Value{}, zerror.BadRequest.WithMsg(`function returns more than 2 vars`)
	} else if t.NumOut() == 2 {
		if t.Out(1) != errorType {
			return nil, reflect.Value{}, zerror.BadRequest.WithMsg(`function second return var must be error`)
		}
	} else if t.NumOut() == 1 {
		if t.Out(0) == errorType {
			return nil, reflect.Value{}, zerror.BadRequest.WithMsg(`func first return var can not be error`)
		}
	}
	returnType := t.Out(0)
	returnKind := returnType.Kind()
	returnBool := returnKind == reflect.Bool
	if !returnBool && returnKind != reflect.Int64 && returnKind != reflect.Float64 && returnKind != reflect.String {
		return nil, reflect.Value{}, zerror.BadRequest.WithMsg(`function first return val type must be one of string, int64, float64 and bool`)
	}
	return returnType, v, nil
}

var builtinFuncs = map[string]interface{}{
	`contains`:        contains,
	`endsWith`:        endsWith,
	`startsWith`:      startsWith,
	`length`:          length,
	`toLower`:         toLower,
	`toUpper`:         toUpper,
	`trim`:            trim,
	`concat`:          concat,
	`geoWithin2d`:     geoWithin2d,
	`hasIntersection`: hasIntersection,
	`timestampBefore`: timestampBefore,
	`now`:             now,
}

func now() string {
	return time.Now().Format(time.RFC3339)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func endsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func startsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func length(s string) int64 {
	return int64(len(s))
}

func toLower(s string) string {
	return strings.ToLower(s)
}

func toUpper(s string) string {
	return strings.ToUpper(s)
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func concat(sa, sb string) string {
	return sa + sb
}

func geoWithin2d(s string, bottomLeftPoint string, topRightPoint string) (bool, error) {
	type point struct {
		x, y float64
	}
	s2p := func(s string) (*point, error) {
		parts := strings.Split(s, ",")
		if len(parts) != 2 {
			return nil, zerror.BadRequest.Errorf("got invalid point string %v, expected: 'x, y'", s)
		}
		sx, sy := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		x, err := strconv.ParseFloat(sx, 64)
		if err != nil {
			return nil, zerror.BadRequest.Errorf("got invalid point x value: %v, err: %v", sx, err)
		}
		y, err := strconv.ParseFloat(sy, 64)
		if err != nil {
			return nil, zerror.BadRequest.Errorf("got invalid point y value: %v, err: %v", sy, err)
		}
		return &point{
			x: x,
			y: y,
		}, nil
	}
	p, err := s2p(s)
	if err != nil {
		return false, err
	}
	bl, err := s2p(bottomLeftPoint)
	if err != nil {
		return false, err
	}
	tr, err := s2p(topRightPoint)
	if err != nil {
		return false, err
	}
	return bl.x <= p.x && p.x <= tr.x && bl.y <= p.y && p.y <= tr.y, nil
}

func hasIntersection(superset, set string) bool {
	if superset == "" || set == "" {
		return false
	}

	parts := strings.Split(superset, ",")
	supersetMap := make(map[string]struct{})
	for _, part := range parts {
		key := strings.TrimSpace(part)
		if key == "" {
			continue
		}
		supersetMap[key] = struct{}{}
	}

	parts = strings.Split(set, ",")
	setMap := make(map[string]struct{})
	for _, part := range parts {
		key := strings.TrimSpace(part)
		if key == "" {
			continue
		}
		setMap[key] = struct{}{}
	}
	if len(supersetMap) == 0 || len(setMap) == 0 {
		return false
	}
	for key := range setMap {
		if _, ok := supersetMap[key]; ok {
			return true
		}
	}
	return false
}

func timestampBefore(sec int64) int64 {
	return time.Now().Unix() - sec
}
