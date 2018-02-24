package assert

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type logger interface {
	Logf(format string, args ...interface{})
}

var repos = map[reflect.Type][]cmp.Option{}

// Add adds options for the type of value to option repository.
func Add(v interface{}, opts ...cmp.Option) {
	if len(opts) == 0 {
		return
	}
	typ, ok := v.(reflect.Type)
	if !ok {
		typ = reflect.ValueOf(v).Type()
	}
	curr, ok := repos[typ]
	if !ok {
		curr = make([]cmp.Option, 0, len(opts))
		repos[typ] = curr
	}
	curr = append(curr, opts...)
}

var empty = make([]cmp.Option, 0)

func check(l logger, x, y interface{}) bool {
	tx := reflect.ValueOf(x).Type()
	ty := reflect.ValueOf(y).Type()
	var opts, ox, oy []cmp.Option
	ox = repos[tx]
	if ty != tx {
		oy = repos[ty]
	}
	if len(ox) == 0 {
		if len(oy) == 0 {
			opts = empty
		} else {
			opts = oy
		}
	} else if len(oy) == 0 {
		opts = ox
	} else {
		opts = make([]cmp.Option, len(ox)+len(oy))
		copy(opts, ox)
		copy(opts[len(ox):], oy)
	}

	diff := cmp.Diff(x, y, opts...)
	if diff == "" {
		return true
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "[N/A]"
		line = -1
	}
	l.Logf("not euqual in file %s at line %d: %s", file, line, diff)
	return false
}

// Equal checks equality between x and y.
// When not equal tb.Fail() is called.
func Equal(tb testing.TB, x, y interface{}) {
	if check(tb, x, y) {
		return
	}
	tb.Fail()
}

// FatalEqual checks equality between x and y.
// When not equal tb.FailNow() is called.
func FatalEqual(tb testing.TB, x, y interface{}) {
	if check(tb, x, y) {
		return
	}
	tb.FailNow()
}
