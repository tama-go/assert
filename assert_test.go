package assert

import "testing"

func TestNil(t *testing.T) {
	Equal(t, nil, nil)
}
