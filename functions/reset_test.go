package functions

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
