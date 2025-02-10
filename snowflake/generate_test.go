package snowflake

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	// Given:
	gen, err := New()
	require.NoError(t, err)

	// When:
	id, err := gen.Generate()

	// Then:
	require.NoError(t, err)
	require.NotZero(t, id)

	// When:
	id, err = gen.Generate()

	// Then:
	require.NoError(t, err)
	require.NotZero(t, id)

	// Given:
	gen, err = New(StartTime(time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)), MachineID(12345))
	require.NoError(t, err)

	// When:
	id, err = gen.Generate()

	// Then:
	require.NoError(t, err)
	require.NotZero(t, id)

	// When:
	id, err = gen.Generate()

	// Then:
	require.NoError(t, err)
	require.NotZero(t, id)
}

func TestGenerator_Generate_CheckUnique(t *testing.T) {
	// Given:
	gen, err := New()
	require.NoError(t, err)

	ids := make([]int64, 10)
	wg := &sync.WaitGroup{}
	wg.Add(10)

	// When:
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			id, err := gen.Generate()
			require.NoError(t, err)
			ids[i] = id
		}(i)
	}
	wg.Wait()

	// Then:
	for i := 0; i < 9; i++ {
		for j := i + 1; j < 10; j++ {
			if ids[i] == ids[j] {
				require.FailNow(t, "has duplicates")
			}
		}
	}
}
