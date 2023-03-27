package oneitempool

import (
	"sync"
	"testing"
)

const sliceSize = 1000

func BenchmarkFreelist(b *testing.B) {
	freelist := New([]int{})
	for i := 0; i < b.N; i++ {
		data := freelist.Get()[:0]
		for x := 0; x < sliceSize; x += 1 {
			data = append(data, x)
		}
		freelist.Put(data)
	}
}

func BenchmarkSyncPool(b *testing.B) {
	pool := &sync.Pool{
		New: func() interface{} {
			s := []int{}
			return &s
		},
	}
	for i := 0; i < b.N; i++ {
		data := pool.Get().(*[]int)
		*data = (*data)[:0]
		for x := 0; x < sliceSize; x += 1 {
			*data = append(*data, x)
		}
		pool.Put(data)
	}
}

func BenchmarkMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := []int{}
		for x := 0; x < sliceSize; x += 1 {
			data = append(data, x)
		}
		_ = data
	}
}

func BenchmarkMakeReuse(b *testing.B) {

	data := []int{}
	for i := 0; i < b.N; i++ {
		data = data[:0]
		for x := 0; x < sliceSize; x += 1 {
			data = append(data, x)
		}
		_ = data
	}
}
