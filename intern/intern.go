package intern

import (
	"golang.org/x/exp/maps"
	"sync"
)

type GenericIntern[T comparable] interface {
	// Deduplicate deduplicates the value by inserting it to the map & reading it back out.
	// If the value is already present, then the cached value is returned (interned value)
	// If the value is new, then a new key is added and the original value is returned after being interned.
	Deduplicate(input T) (output T)

	// Insert attempts to insert key to the map, to get the value back later call Value(id) to retrieve the interned value.
	// If the value is already present in the map, then it returns the existing unique id.
	// If the value is new, then it increments the unique id counter & returns the new unique id.
	Insert(input T) (uniqueID uint64)

	// Value returns the interned value for the given unique id
	// If the key is present in the map then the output T will be the interned value, and `ok=true`.
	// If the key is not present, then `ok=false` will be returned.
	Value(uniqueID uint64) (output T, ok bool)

	// Len returns the size of the map (number of keys)
	Len() int

	// Clear deletes the interned values & reset the counter back to 0.
	// This allows you to free the interned values & put the Intern struct back into a safe state for a sync.Pool or garbage collection.
	Clear()
}

var _ GenericIntern[int] = &genericIntern[int]{}
var _ GenericIntern[int] = &safeGeneric[int]{}

// genericIntern implements the GenericIntern interface for any comparable type
type genericIntern[T comparable] struct {
	keys    map[T]uint64
	values  map[uint64]T
	counter uint64
}

// New creates a new GenericIntern[T] instance
func New[T comparable]() GenericIntern[T] {
	return &genericIntern[T]{
		keys:    map[T]uint64{},
		values:  map[uint64]T{},
		counter: 0,
	}
}

func (i *genericIntern[T]) Deduplicate(input T) T {
	uniqueId := i.Insert(input)
	output, _ := i.Value(uniqueId)
	return output
}

func (i *genericIntern[T]) Insert(input T) uint64 {
	if uniqueId, ok := i.keys[input]; ok {
		return uniqueId
	}

	i.counter++
	uniqueId := i.counter
	i.keys[input] = uniqueId
	i.values[uniqueId] = input
	return uniqueId
}

func (i *genericIntern[T]) Value(uniqueID uint64) (output T, ok bool) {
	output, ok = i.values[uniqueID]
	return output, ok
}

func (i *genericIntern[T]) Len() int {
	return int(i.counter)
}

func (i *genericIntern[T]) Clear() {
	maps.Clear(i.keys)
	maps.Clear(i.values)
	i.counter = 0
}

// safeGeneric wraps the genericIntern struct & makes it thread safe
type safeGeneric[T comparable] struct {
	intern GenericIntern[T]
	sync.RWMutex
}

// NewSafe creates a new GenericIntern[T] instance
func NewSafe[T comparable]() GenericIntern[T] {
	return &safeGeneric[T]{intern: New[T]()}
}

func (i *safeGeneric[T]) Deduplicate(input T) (output T) {
	i.RWMutex.Lock()
	output = i.intern.Deduplicate(input)
	i.RWMutex.Unlock()
	return output
}

func (i *safeGeneric[T]) Insert(input T) (uniqueID uint64) {
	i.RWMutex.Lock()
	uniqueID = i.intern.Insert(input)
	i.RWMutex.Unlock()
	return uniqueID
}

func (i *safeGeneric[T]) Value(index uint64) (output T, ok bool) {
	i.RWMutex.RLock()
	output, ok = i.intern.Value(index)
	i.RWMutex.RUnlock()
	return output, ok
}

func (i *safeGeneric[T]) Len() (output int) {
	i.RWMutex.RLock()
	output = i.intern.Len()
	i.RWMutex.RUnlock()
	return output
}

func (i *safeGeneric[T]) Clear() {
	i.Lock()
	i.intern.Clear()
	i.Unlock()
}
