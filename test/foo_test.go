package foo

import (
	"testing"

	"github.com/tama-go/assert"
)

type foo Foo

func TestNew(t *testing.T) {
	a := New("aaa", 123)
	b := New("aaa", 123)
	assert.FatalEqual(t, a, b)
}
