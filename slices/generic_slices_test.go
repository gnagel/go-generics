package slices

import (
	"github.com/go-generics-playground/generics/values"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDefaultAllocSlice(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		alloc := DefaultAllocSlice[int](values.DefaultAlloc[int](123))

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
		alloc := DefaultAllocSlice[time.Time](values.DefaultAlloc[time.Time](time.Unix(0, 0).UTC()))

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

func TestDefaultResetSlice(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		reset := DefaultResetSlice[int](nil)

		input := make([]int, 0, 10)
		output := reset(input)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 10)

		input = []int{123}
		output = reset(input)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 1)
		require.EqualValues(t, input[0], 0)
	})

	t.Run("time.Time", func(t *testing.T) {
		reset := DefaultResetSlice[time.Time](nil)

		input := make([]time.Time, 0, 10)
		output := reset(input)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 10)

		input = []time.Time{time.Now()}
		output = reset(input)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 1)
		require.EqualValues(t, input[0], time.Time{})
	})
}

func TestDefaultDeallocSlice(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		dealloc := DefaultDeallocSlice[int](nil)

		input := make([]int, 0, 10)
		output := dealloc(input)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 0)

		input = []int{123}
		output = dealloc(input)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 0)
		require.EqualValues(t, input[0], 0)
	})

	t.Run("time.Time", func(t *testing.T) {
		dealloc := DefaultDeallocSlice[time.Time](nil)

		input := make([]time.Time, 0, 10)
		output := dealloc(input)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 0)

		input = []time.Time{time.Now()}
		output = dealloc(input)
		require.NotNil(t, output)
		require.Len(t, output, 0)
		require.EqualValues(t, cap(output), 0)
		require.EqualValues(t, input[0], time.Time{})
	})
}
