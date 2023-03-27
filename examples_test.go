package oneitempool

import "fmt"

type ExampleStruct struct {
	SomeIDs   []int
	SomeNames []string
}

func ExampleOneItemPool() {
	pool := New([]float64{})
	frenquentlyCalledFunction(pool)
	frenquentlyCalledFunction(pool)
	frenquentlyCalledFunction(pool)
	frenquentlyCalledFunction(pool)

	floats := pool.Get()[:0]
	fmt.Printf("cap(floats) >= 1000? %v\n", cap(floats) >= 1000)
	pool.Put(floats)

	// Output: cap(floats) >= 1000? true
}

func frenquentlyCalledFunction(pool *OneItemPool[[]float64]) {

	floats := pool.Get()[:0]
	// defer pool.Put(floats) would evaluate floats immediately, which we don't want.
	defer func() {
		pool.Put(floats)
	}()
	for i := 0; i < 1000; i += 1 {
		floats = append(floats, 0.1)
	}
}

// pools with several objects that we want to re-use
type Pools struct {
	pointerExample *OneItemPool[*ExampleStruct]
	example        *OneItemPool[ExampleStruct]
	someFloats     *OneItemPool[[]float64]
}

func ExamplePools() {
	pools := Pools{
		pointerExample: New(&ExampleStruct{}),
		example:        New(ExampleStruct{}),
		someFloats:     New(make([]float64, 0, 1000)), // let's preallocate some
	}
	frenquentlyCalledFunction2(pools)
	frenquentlyCalledFunction2(pools)
	frenquentlyCalledFunction2(pools)
	cap1, cap2, cap3 := frenquentlyCalledFunction2(pools)

	p := pools.pointerExample.Get()
	defer pools.pointerExample.Put(p)
	fmt.Printf("saved allocations1 %v, ", cap1 == cap(p.SomeIDs))

	s := pools.example.Get()
	defer pools.example.Put(s)
	fmt.Printf("saved allocations2 %v, ", cap2 == cap(s.SomeNames))

	floats := pools.someFloats.Get()

	defer pools.someFloats.Put(floats)
	fmt.Printf("saved allocations3 %v", cap3 == cap(floats))

	// Output: saved allocations1 true, saved allocations2 true, saved allocations3 true

}

func frenquentlyCalledFunction2(pools Pools) (int, int, int) {

	p := pools.pointerExample.Get()
	p.SomeIDs = p.SomeIDs[:0]
	p.SomeNames = p.SomeNames[:0]

	s := pools.example.Get()
	s.SomeIDs = s.SomeIDs[:0]
	s.SomeNames = s.SomeNames[:0]

	floats := pools.someFloats.Get()[:0]

	defer func() {
		pools.pointerExample.Put(p)
		pools.example.Put(s)
		pools.someFloats.Put(floats)
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
