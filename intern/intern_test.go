package intern

import (
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
)

func TestIntern(t *testing.T) {
	i := New[string]()
	require.Zero(t, i.Len())

	output, ok := i.Value(0)
	require.False(t, ok)
	require.Zero(t, output)

	index := i.Insert("value")
	require.NotZero(t, index)

	output, ok = i.Value(index)
	require.True(t, ok)
	require.NotZero(t, output)
	require.EqualValues(t, output, "value")

	index2 := i.Insert("value")
	require.NotZero(t, index2)
	require.EqualValues(t, index, index2)

	i.Clear()
	require.Zero(t, i.Len())

	inputs := []string{faker.Word(), faker.Word(), faker.Word()}
	inputs = append(inputs, inputs...)
	inputs = append(inputs, inputs...)
	rand.Shuffle(len(inputs), func(i, j int) {
		inputs[i] = inputs[j]
	})

	keys := map[uint64]uint64{}
	for _, input := range inputs {
		key := i.Insert(input)
		keys[key] = key
	}
	require.Len(t, keys, 3)
	require.EqualValues(t, i.Len(), 3)
}

func TestSafeIntern(t *testing.T) {
	i := NewSafe[string]()
	output, ok := i.Value(0)
	require.False(t, ok)
	require.Zero(t, output)

	index := i.Insert("value")
	require.NotZero(t, index)

	output, ok = i.Value(index)
	require.True(t, ok)
	require.NotZero(t, output)
	require.EqualValues(t, output, "value")

	index2 := i.Insert("value")
	require.NotZero(t, index2)
	require.EqualValues(t, index, index2)

	i.Clear()
	require.Zero(t, i.Len())

	inputs := []string{faker.Word(), faker.Word(), faker.Word()}
	inputs = append(inputs, inputs...)
	inputs = append(inputs, inputs...)
	rand.Shuffle(len(inputs), func(i, j int) {
		inputs[i] = inputs[j]
	})

	var wg sync.WaitGroup
	for idx := 0; idx < 10; idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, input := range inputs {
				i.Deduplicate(input)
			}
		}()
	}
	wg.Wait()
	require.EqualValues(t, i.Len(), 3)
}
