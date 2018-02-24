package assert

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var repos = map[reflect.Type][]cmp.Option{}

// Add adds options for the type of value to option repository.
func Add(v interface{}, opts ...cmp.Option) {
	if v == nil || len(opts) == 0 {
		return
	}
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(v)
	}
	curr, ok := repos[typ]
	if !ok {
		curr = make([]cmp.Option, 0, len(opts))
	}
	repos[typ] = append(curr, opts...)
}

func get(typ reflect.Type) []cmp.Option {
	if typ == nil {
		return nil
	}
	return repos[typ]
}

type logger interface {
	Helper()
	Logf(format string, args ...interface{})
}

func check(l logger, x, y interface{}) bool {
	l.Helper()

	tx, ty := reflect.TypeOf(x), reflect.TypeOf(y)
	opts := get(tx)
	if tx != ty {
		opts2 := get(ty)
		if len(opts2) > 0 {
			tmp := make([]cmp.Option, 0, len(opts)+len(opts2))
			opts = append(append(tmp, opts...), opts2...)
		}
	}

	diff := cmp.Diff(x, y, opts...)
	if diff == "" {
		return true
	}
	l.Logf("not equal\n%s", diff)
	return false
}

// Equal checks equality between x and y.
// When not equal tb.Fail() is called.
func Equal(tb testing.TB, x, y interface{}) {
	tb.Helper()
	if check(tb, x, y) {
		return
	}
	tb.Fail()
}

// FatalEqual checks equality between x and y.
// When not equal tb.FailNow() is called.
func FatalEqual(tb testing.TB, x, y interface{}) {
	tb.Helper()
	if check(tb, x, y) {
		return
	}
	tb.FailNow()
}
