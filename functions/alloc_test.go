package functions

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

func TestDefaultAllocSlice(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		alloc := DefaultAllocSlice[int](DefaultAlloc[int](123))

		output := alloc(0, 10)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 10)

		output = alloc(1, 10)
		require.NotNil(t, output)
		require.Len(t, output, 1)
		require.EqualValues(t, cap(output), 10)
		require.EqualValues(t, output[0], 123)
	})

	t.Run("time.Time", func(t *testing.T) {
		alloc := DefaultAllocSlice[time.Time](DefaultAlloc[time.Time](time.Unix(0, 0).UTC()))

		output := alloc(0, 10)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 10)

		output = alloc(1, 10)
		require.NotNil(t, output)
		require.Len(t, output, 1)
		require.EqualValues(t, cap(output), 10)
		require.EqualValues(t, output[0], time.Unix(0, 0).UTC())
	})
}
