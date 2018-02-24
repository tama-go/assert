package foo

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tama-go/assert"
)

type foo Foo

func init() {
	assert.Add(&foo{}, cmpopts.IgnoreFields(foo{}, "Int", "Rnd"))
}

func TestStruct(t *testing.T) {
	a := New("aaa", 123)
	b := New("aaa", 123)
	assert.FatalEqual(t, a, b)
}

func TestTypeAlias(t *testing.T) {
	a := (*foo)(New("aaa", 123))
	b := (*foo)(New("aaa", 456))
	assert.FatalEqual(t, a, b)
}
