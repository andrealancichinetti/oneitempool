package oneitempool

import "fmt"

type ExampleStruct struct {
	SomeIDs   []int
	SomeNames []string
}

// a pool with several objects that we want to re-use
type Pool struct {
	pointerExample *OneItemPool[*ExampleStruct]
	example        *OneItemPool[ExampleStruct]
	someFloats     *OneItemPool[[]float64]
}

func ExamplePool() {
	pool := Pool{
		pointerExample: New(&ExampleStruct{}),
		example:        New(ExampleStruct{}),
		someFloats:     New(make([]float64, 0, 1000)), // let's preallocate some
	}
	hotSpot(pool)
	hotSpot(pool)
	hotSpot(pool)
	cap1, cap2, cap3 := hotSpot(pool)

	p := pool.pointerExample.Get()
	defer pool.pointerExample.Put(p)
	fmt.Printf("saved allocations1 %v, ", cap1 == cap(p.SomeIDs))

	s := pool.example.Get()
	defer pool.example.Put(s)
	fmt.Printf("saved allocations2 %v, ", cap2 == cap(s.SomeNames))

	floats := pool.someFloats.Get()

	defer pool.someFloats.Put(floats)
	fmt.Printf("saved allocations3 %v", cap3 == cap(floats))

	// Output: saved allocations1 true, saved allocations2 true, saved allocations3 true

}

func hotSpot(pool Pool) (int, int, int) {

	p := pool.pointerExample.Get()
	p.SomeIDs = p.SomeIDs[:0]
	p.SomeNames = p.SomeNames[:0]

	s := pool.example.Get()
	s.SomeIDs = s.SomeIDs[:0]
	s.SomeNames = s.SomeNames[:0]

	floats := pool.someFloats.Get()[:0]

	// defer pool.someFloats.Put(floats) would evaluate floats immediately,
	// which we don't want.
	defer func() {
		pool.pointerExample.Put(p)
		pool.example.Put(s)
		pool.someFloats.Put(floats)
	}()

	for i := 0; i < 10000; i += 1 {
		p.SomeIDs = append(p.SomeIDs, i)
		p.SomeNames = append(p.SomeNames, "a")

		s.SomeIDs = append(p.SomeIDs, i)
		s.SomeNames = append(p.SomeNames, "a")

		floats = append(floats, 0.1)
	}

	return cap(p.SomeIDs), cap(s.SomeNames), cap(floats)
}
