package values

// Alloc is a generic function for creating a new instance of type T.
type Alloc[T any] func() T

// Reset zeros out a value and resets it back to zero state.
type Reset[T any] func(T) T

// Dealloc releases a value back to the garbage collector and releases any internals.
type Dealloc[T any] func(T)

// Alloc, Dealloc, and Reset all implement the generic interface
var _ Alloc[any] = DefaultAlloc[any](nil)
var _ Dealloc[any] = DefaultDealloc[any]
var _ Reset[any] = DefaultReset[any]

// DefaultAlloc implements the Alloc[T] interface for generic values
func DefaultAlloc[T any](defaultValue T) Alloc[T] {
	return func() T {
		var t = defaultValue
		return t
	}
}

// DefaultReset implements the Reset[T] interface for generic values
func DefaultReset[T any](t T) T {
	var v T
	t = v
	return t
}

// DefaultDealloc implements the Dealloc[T] interface for generic values
func DefaultDealloc[T any](_ T) {}
