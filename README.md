# oneitempool
A bare-bones free list for Go, not safe for concurrent use.

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
