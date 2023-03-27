package oneitempool

/*
OneItemPool is a very basic implementation of a pool to save memory allocations.
It is NOT SAFE FOR CONCURRENT USE.
It simply forces you to get an item, and put it back before being able to get it again.
*/
type OneItemPool[T any] struct {
	accessed bool
	data     T
}

/*
New return a pointer to a new pool.
Example:
pool := New(make([]float64, 0, 1000))
If you are passing the pool as a parameter to a function, you should use the pointer.
*/
func New[T any](t T) *OneItemPool[T] {
	return &OneItemPool[T]{data: t}
}

/*
Get your item. You cannot get an item twice.
You should initialize it before using it.
*/
func (p *OneItemPool[T]) Get() T {
	if p.accessed {
		panic("requesting data twice")
	}
	p.accessed = true
	return p.data
}

/*
Put the item back in the pool This is necessary to be able to Get it again.
You shouldn't keep any reference to the item after putting it back.
*/
func (p *OneItemPool[T]) Put(t T) {
	p.accessed = false
	p.data = t
}
