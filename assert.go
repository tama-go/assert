package assert

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// CommonOptions is common options which applied for all types.
var CommonOptions = []cmp.Option{
	cmpopts.IgnoreInterfaces(struct{ context.Context }{}),
}

// AddCommon adds common cmp.Options for all types.
func AddCommon(opts ...cmp.Option) {
	CommonOptions = append(CommonOptions, opts...)
}

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

func join(oo ...[]cmp.Option) []cmp.Option {
	n := 0
	for _, o := range oo {
		n += len(o)
	}
	opts := make([]cmp.Option, 0, n)
	for _, o := range oo {
		if len(o) == 0 {
			continue
		}
		opts = append(opts, o...)
	}
	return opts
}

func check(t testing.TB, expected, actual interface{}) bool {
	var opts2 []cmp.Option
	texp, tact := reflect.TypeOf(expected), reflect.TypeOf(actual)
	opts := get(texp)
	if texp != tact {
		opts2 = get(tact)
	}

	diff := cmp.Diff(expected, actual, join(CommonOptions, opts, opts2)...)
	if diff == "" {
		return true
	}
	t.Logf("not equal\nName: %s\nDiff:\n--- Expected\n+++ Actual\n%s", t.Name(), diff)
	return false
}

// Equal checks equality between "expected" and "actual".
// Test will be aborted soon if it wasn't match.
func Equal(tb testing.TB, expected, actual interface{}) {
	tb.Helper()
	if check(tb, expected, actual) {
		return
	}
	tb.FailNow()
}

// FatalEqual checks equality between expected and actual.
func FatalEqual(tb testing.TB, expected, actual interface{}) {
	tb.Helper()
	if check(tb, expected, actual) {
		return
	}
	tb.FailNow()
}
