# oneitempool

OneItemPool is a very basic implementation of a pool to save memory allocations.
It is NOT SAFE FOR CONCURRENT USE.

Why having a pool with a single item?

Say that you have a function that gets called many times and
that requires some memory allocations.

Instead of allocating the memory each time you call the function,
a standard solution might be to use sync.Pool.

A more performant alternative is to allocate memory
before calling the function you want to optimise,
and pass the allocated variable as a parameter to the function.
This works great, but you might end up with a shared mutable object,
and you need to make sure you are not changing it in multiple places.

Example:

```go

allocated := &[]float64{}
doSomething(allocated)

func doSomething(allocated *[]float64) {

	*allocated = append((*allocated)[:0], 0.1)
	// you accidentally use allocated in another function.
	doSomethingElse(allocated)
	*allocated = append(*allocated, 0.1)

	// now allocated is [10.0, 0.1]
	// this is probably a bug, or at at least it's confusing.

}

func doSomethingElse(allocated *[]float64) {
	*allocated = append(*allocated[:0], 10.0)
	// do something else
}

```

OneItemPool will not make the item available until we put it back.
Simply follow the rule that after using Put(), you should not hold any reference to the item, 
because it might be changed later.

```go

allocated := New([]float64{})
doSomething(allocated)

func doSomething(pool *OneItemPool[[]float64]) {
	allocated := pool.Get()[:0]
	allocated = append(allocated[:0], 0.1)
	// after using Put, we should not refer to allocated anymore.
	pool.Put(allocated)
}

```

Now allocated is protected from accidental changes: 

```go


allocated := New([]float64{})
doSomething(allocated)

func doSomething(pool *OneItemPool[[]float64]) {

	allocated := pool.Get()[:0]

	allocated = append(allocated[:0], 0.1)
	// you accidentally use the pool in another function:
	// it will panic.
	doSomethingElse(pool)
	allocated = append(allocated, 0.1)

	pool.Put(allocated)
}

func doSomethingElse(pool *OneItemPool[[]float64]) {
	allocated := pool.Get()[:0]
	allocated = append(allocated, 0.1)
	pool.Put(allocated)
}

```


This is clearly more limited than sync.Pool, but if you need performance, it might be faster.
Here is a (small and artificial) benchmark:

```
BenchmarkOneItemPool-8   	 2877622	       416.1 ns/op
BenchmarkSyncPool-8      	  931642	      1285 ns/op
```


