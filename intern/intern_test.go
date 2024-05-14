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

	inputs := randomStringInputs(1_00)
	unique := make(map[string]string, len(inputs))
	for _, k := range inputs {
		unique[k] = k
	}

	keys := map[uint64]uint64{}
	for _, input := range inputs {
		key := i.Insert(input)
		keys[key] = key
	}
	require.Len(t, keys, len(unique))
	require.EqualValues(t, i.Len(), len(unique))
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

	inputs := randomStringInputs(1_00)
	unique := make(map[string]string, len(inputs))
	for _, k := range inputs {
		unique[k] = k
	}

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
	require.EqualValues(t, i.Len(), len(unique))
}

func BenchmarkIntern(b *testing.B) {
	inputs := randomStringInputs(10_000)

	i := New[string]()

	for index := 0; index < b.N; index++ {
		i.Deduplicate(inputs[index%len(inputs)])
	}
}

func BenchmarkSafeIntern(b *testing.B) {
	inputs := randomStringInputs(10_000)

	i := NewSafe[string]()

	rows := make(chan string)

	var wg sync.WaitGroup
	for idx := 0; idx < 10; idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rows {
				i.Deduplicate(row)
			}
		}()
	}

	for index := 0; index < b.N; index++ {
		rows <- inputs[index%len(inputs)]
	}
	close(rows)
	wg.Wait()
}

func randomStringInputs(maxLen int) []string {
	inputs := make([]string, maxLen)
	for index := 0; index < len(inputs); index++ {
		inputs[index] = faker.Sentence()
	}
	rand.Shuffle(len(inputs), func(i, j int) {
		inputs[i] = inputs[j]
	})
	return inputs
}
