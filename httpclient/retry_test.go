package httpclient

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeoutAndRetryOption_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		opt     timeoutAndRetryOption
		wantErr bool
	}{
		{
			name: "valid",
			opt: timeoutAndRetryOption{
				maxRetries:         1,
				maxWaitPerTry:      1 * time.Second,
				maxWaitInclRetries: 2 * time.Second,
				onTimeout:          false,
				onStatusCodes:      map[int]bool{},
			},
			wantErr: false,
		},
		{
			name: "maxWaitPerTry > maxWaitInclRetries",
			opt: timeoutAndRetryOption{
				maxRetries:         1,
				maxWaitPerTry:      3 * time.Second,
				maxWaitInclRetries: 2 * time.Second,
				onTimeout:          false,
				onStatusCodes:      map[int]bool{},
			},
			wantErr: true,
		},
		{
			name: "maxWaitPerTry < 0",
			opt: timeoutAndRetryOption{
				maxRetries:         1,
				maxWaitPerTry:      -1 * time.Second,
				maxWaitInclRetries: 2 * time.Second,
				onTimeout:          false,
				onStatusCodes:      map[int]bool{},
			},
			wantErr: true,
		},
		{
			name: "maxWaitInclRetries < 0",
			opt: timeoutAndRetryOption{
				maxRetries:         1,
				maxWaitPerTry:      1 * time.Second,
				maxWaitInclRetries: -1 * time.Second,
				onTimeout:          false,
				onStatusCodes:      map[int]bool{},
			},
			wantErr: true,
		},
		{
			name: "onTimeout true but maxRetries 0",
			opt: timeoutAndRetryOption{
				maxRetries:         0,
				maxWaitPerTry:      1 * time.Second,
				maxWaitInclRetries: 2 * time.Second,
				onTimeout:          true,
				onStatusCodes:      map[int]bool{},
			},
			wantErr: true,
		},
		{
			name: "onStatusCodes not empty but maxRetries 0",
			opt: timeoutAndRetryOption{
				maxRetries:         0,
				maxWaitPerTry:      1 * time.Second,
				maxWaitInclRetries: 2 * time.Second,
				onTimeout:          false,
				onStatusCodes:      map[int]bool{500: true},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opt.IsValid()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestExecWithRetry_Success(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	f := func() error {
		callCount++
		if callCount < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := execWithRetry(ctx, 3, 10*time.Second, f)
	require.NoError(t, err)
	require.Equal(t, 3, callCount)
}

func TestExecWithRetry_ExceedRetries(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	f := func() error {
		callCount++
		return errors.New("persistent error")
	}

	err := execWithRetry(ctx, 2, 10*time.Second, f)
	require.Error(t, err)
	require.Equal(t, 3, callCount) // initial + 2 retries
}

func TestExecWithRetry_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	f := func() error {
		cancel()
		return errors.New("error")
	}

	err := execWithRetry(ctx, 3, 1*time.Second, f)
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)
}
