package oneitempool

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	value := 42
	freelist := New(&value)
	freelist.Get()
	freelist.Put(&value)
	result := freelist.Get()
	require.Equal(t, 42, *result)
}

func TestPanic(t *testing.T) {
	value := 42
	freelist := New(&value)
	freelist.Get()
	defer func() { _ = recover() }()
	freelist.Get()
	t.Errorf("did not panic")
}

func TestSlice(t *testing.T) {
	value := []int{1, 2, 3}
	pool := New(value)
	x := pool.Get()
	x = append(x, 109)
	x = append(x, 109)
	x = append(x, 109)
	x = append(x, 109)
	x = append(x, 109)
	x = append(x, 109)
	x = append(x, 109)
	x = append(x, 109)
	pool.Put(x)
	y := pool.Get()
	require.Equal(t, cap(x), cap(y))
}
