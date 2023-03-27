package oneitempool

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
}

func hotSpot(pool Pool) {

	p := pool.pointerExample.Get()
	p.SomeIDs = p.SomeIDs[:0]
	p.SomeNames = p.SomeNames[:0]
	defer pool.pointerExample.Put(p)

	s := pool.example.Get()
	s.SomeIDs = p.SomeIDs[:0]
	s.SomeNames = p.SomeNames[:0]
	defer pool.example.Put(s)

	floats := pool.someFloats.Get()[:0]
	defer pool.someFloats.Put(floats)

	for i := 0; i < 10000; i += 1 {
		p.SomeIDs = append(p.SomeIDs, i)
		p.SomeNames = append(p.SomeNames, "a")

		s.SomeIDs = append(p.SomeIDs, i)
		s.SomeNames = append(p.SomeNames, "a")

		floats = append(floats, 0.1)
	}

}
