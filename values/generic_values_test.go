package values

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDefaultAlloc(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		alloc := DefaultAlloc[int](0)
		output := alloc()
		require.Zero(t, output)

		alloc = DefaultAlloc[int](123)
		output = alloc()
		require.EqualValues(t, output, 123)
	})

	t.Run("time.Time", func(t *testing.T) {
		alloc := DefaultAlloc[time.Time](time.Unix(0, 0).In(time.UTC))
		output := alloc()
		require.Zero(t, output.Unix())
		require.False(t, output.IsZero(), output)

		alloc = DefaultAlloc[time.Time](time.Time{})
		output = alloc()
		require.NotZero(t, output.Unix())
		require.True(t, output.IsZero(), output)
	})
}

func TestDefaultReset(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		reset := DefaultReset[int]
		output := reset(123)
		require.Zero(t, output)
	})

	t.Run("time.Time", func(t *testing.T) {
		reset := DefaultReset[time.Time]
		value := time.Unix(0, 0).In(time.UTC)
		require.Zero(t, value.Unix())
		require.False(t, value.IsZero(), value)

		output := reset(value)
		require.NotZero(t, output.Unix())
		require.True(t, output.IsZero(), output)

		value = time.Now()
		output = reset(value)
		require.NotZero(t, output.Unix())
		require.True(t, output.IsZero(), output)
	})
}

func TestDefaultDealloc(t *testing.T) {
	t.Run("int ptr", func(t *testing.T) {
		var reset Reset[**int] = func(i **int) **int {
			*i = nil
			return i
		}
		var value int = 123
		var ptr *int = &value
		var prtptr **int = &ptr

		output := reset(prtptr)
		require.NotNil(t, output)
		require.Nil(t, *output)
	})
}
