package snowflake

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	type mockFlakeData struct {
		useMock bool
		out     uint64
		err     error
	}
	tcs := map[string]struct {
		opts          []Option
		mockFlakeData mockFlakeData
		expErr        error
	}{
		"default generator": {
			opts: nil,
		},
		"generator with options": {
			opts: []Option{
				StartTime(time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)),
				MachineID(12345),
			},
		},
		"error from flake": {
			mockFlakeData: mockFlakeData{
				useMock: true,
				err:     errors.New("simulated error"),
			},
			expErr: errors.New("snowflake ID generation failed: simulated error"),
		},
		"flake return 0": {
			mockFlakeData: mockFlakeData{
				useMock: true,
				out:     0,
			},
			expErr: errors.New("snowflake ID is invalid: 0"),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given:
			gen, err := New(tc.opts...)
			require.NoError(t, err)

			if tc.mockFlakeData.useMock {
				mockIDProvider := NewMockidProvider(t)
				mockIDProvider.
					EXPECT().
					NextID().
					Return(tc.mockFlakeData.out, tc.mockFlakeData.err)
				gen.flake = mockIDProvider
			}

			// When:
			if tc.expErr != nil {
				_, err := gen.Generate()
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				id, err := gen.Generate()
				require.NoError(t, err)
				require.NotZero(t, id)

				id, err = gen.Generate()
				require.NoError(t, err)
				require.NotZero(t, id)
			}
		})
	}
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
