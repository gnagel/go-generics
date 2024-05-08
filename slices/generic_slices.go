package values

import "github.com/go-generics-playground/generics/values"

// AllocSlice is a generic function for creating a new slice of type T.
type AllocSlice[T any] func(defaultLen, defaultCap int) []T

// ResetSlice is a generic function for resetting an existing slice of type T.
// This returns an empty slice of length 0, but the capacity can be > 0
type ResetSlice[T any] func([]T) []T

// DeallocSlice is a generic function for releasing an existing slice of type T.
// This is the same as slices.Clip, except with the added call to dealloc for each element
// This returns an empty slice of length 0 & capacity 0
type DeallocSlice[T any] func([]T) []T

// AllocSlice, DeallocSlice, and ResetSlice all implement the generic interface
var _ AllocSlice[any] = DefaultAllocSlice[any](nil)
var _ ResetSlice[any] = DefaultResetSlice[any](nil)
var _ DeallocSlice[any] = DefaultDeallocSlice[any](nil)

// DefaultAllocSlice implements the AllocSlice[T] interface for generic values
func DefaultAllocSlice[T any](alloc values.Alloc[T]) AllocSlice[T] {
	if nil == alloc {
		var t T
		alloc = values.DefaultAlloc[T](t)
	}

	return func(defaultLen, defaultCap int) []T {
		output := make([]T, defaultLen, defaultCap)
		for index := range output[:defaultLen] {
			output[index] = alloc()
		}
		return output
	}
}

// DefaultResetSlice implements the ResetSlice[T] interface for generic values
func DefaultResetSlice[T any](reset values.Reset[T]) ResetSlice[T] {
	if nil == reset {
		reset = values.DefaultReset[T]
	}
	return func(input []T) []T {
		for index, value := range input {
			input[index] = reset(value)
		}
		return input[:0:cap(input)]
	}
}

// DefaultDeallocSlice implements the DeallocSlice[T] interface for generic values
func DefaultDeallocSlice[T any](dealloc values.Dealloc[T]) DeallocSlice[T] {
	if nil == dealloc {
		dealloc = values.DefaultDealloc[T]
	}

	return func(input []T) []T {
		for index, value := range input {
			dealloc(value)
			var t T
			input[index] = t
		}
		return input[:0:0]
	}
}
