package foo

import (
	"math/rand"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tama-go/assert"
)

func init() {
	assert.Add(&Foo{}, cmpopts.IgnoreFields(Foo{}, "Rnd"))
}

type Foo struct {
	Str string
	Int int
	Rnd int
}

func New(s string, n int) *Foo {
	return &Foo{
		Str: s,
		Int: n,
		Rnd: rand.Int(),
	}
}
