package functions

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
