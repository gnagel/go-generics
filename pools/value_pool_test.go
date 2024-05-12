package pools

import (
	"github.com/go-generics-playground/generics/functions"
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestValuePool(t *testing.T) {
	now := time.Now().UTC()
	var alloc = functions.DefaultAlloc[time.Time](now)
	var reset functions.Reset[time.Time] = functions.DefaultReset[time.Time]

	pool := NewValuePool[time.Time](alloc, reset)
	require.Zero(t, pool.Cap())
	require.Zero(t, pool.Len())

	output := pool.Get()
	require.EqualValues(t, output, now)
	time.Sleep(time.Millisecond)
	require.EqualValues(t, pool.Cap(), 1)
	require.EqualValues(t, pool.Len(), 1)

	output = output.Add(time.Minute)
	pool.Put(output)
	time.Sleep(time.Millisecond)
	require.EqualValues(t, pool.Cap(), 1)
	require.EqualValues(t, pool.Len(), 0)

	output = pool.Get()
	require.EqualValues(t, output, now)
	time.Sleep(time.Millisecond)
	require.GreaterOrEqual(t, int(pool.Cap()), 1)
	require.EqualValues(t, pool.Len(), 1)
}

func TestValuePool_race(t *testing.T) {
	now := time.Now().UTC()
	var alloc functions.Alloc[*time.Time] = func() *time.Time {
		t := now
		return &t
	}
	var reset functions.Reset[*time.Time] = func(t *time.Time) *time.Time {
		*t = now
		return t
	}

	pool := NewValuePool[*time.Time](alloc, reset)
	require.Zero(t, pool.Cap())
	require.Zero(t, pool.Len())

	var wg sync.WaitGroup

	for index, maxIndex := 0, runtime.NumGoroutine(); index < maxIndex; index++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value := pool.Get()
			*value = time.Now().UTC()
			pool.Put(value)
		}()
	}

	wg.Wait()
	time.Sleep(time.Millisecond)
	require.NotZero(t, pool.Cap())
	require.Zero(t, pool.Len())
}
