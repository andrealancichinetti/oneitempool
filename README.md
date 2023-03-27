# oneitempool

OneItemPool is a very basic implementation of a pool to save memory allocations.
It is NOT SAFE FOR CONCURRENT USE.

Why having a pool with a single item?

Say that you have a function that gets called many times and
that requires some memory allocations.
Instead of allocating the memory each time you call the function,
you can use sync.Pool, right?
Yes, but you get some overhead. Also, you might not want to
lose the object if the GC runs.

If that doesn't work for you, you can create your object
before calling the function you want to optimise,
and pass it there.
This works great, but I found two inconveniences:
 1. if you have a slice, you need to make sure to update it in the caller.
 2. you end up with a shared mutable object, which needs special care:
    you want to re-use the same object, but not in multiple places at the same time.

If you are using a single thread,
this simple struct will make sure that your code will get access to the data only once
and will make updates convenient.

```

pool := New([]float64{})
frenquentlyCalledFunction(pool)
frenquentlyCalledFunction(pool)
frenquentlyCalledFunction(pool)
frenquentlyCalledFunction(pool)

floats := pool.Get()[:0]
fmt.Printf("cap(floats) >= 1000? %v\n", cap(floats) >= 1000)
pool.Put(floats)

// Output: cap(floats) >= 1000? true


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


```


This is clearly more limited than sync.Pool, but if you need performance, it might be faster.
Here is a (small and artificial) benchmark:

```
BenchmarkOneItemPool-8   	 2877622	       416.1 ns/op
BenchmarkSyncPool-8      	  931642	      1285 ns/op
```


