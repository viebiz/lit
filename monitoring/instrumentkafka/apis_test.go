package instrumentkafka

import (
	"bytes"
	"context"
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mocktracer"
	"github.com/viebiz/lit/monitoring"
)

func TestInjectAsyncSegmentInfo(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	// GIVEN
	buf := bytes.NewBuffer(nil)
	m, err := monitoring.New(monitoring.Config{
		ServerName:  "lightning",
		Environment: "test",
		Version:     "1.0.0",
		Writer:      buf,
	})
	require.NoError(t, err)

	ctx := context.Background()
	_, span := tracer.Start(ctx, "test-parent-span")
	seg := buildPublishSegment(span)

	pm := &sarama.ProducerMessage{
		Key:       sarama.StringEncoder("test-key"),
		Partition: 0,
		Offset:    100,
	}

	var sendErr error

	// WHEN
	monitor := InjectAsyncSegmentInfo(m, seg, pm, sendErr)

	// THEN
	require.NotNil(t, monitor)
}

func TestExtractSpanIDFromProducerMessage(t *testing.T) {
	tcs := map[string]struct {
		msg       *sarama.ProducerMessage
		expSpanID string
	}{
		"no traceparent header": {
			msg: &sarama.ProducerMessage{
				Headers: []sarama.RecordHeader{
					{
						Key: []byte("key"),
					},
				},
			},
		},
		"invalid traceparent header": {
			msg: &sarama.ProducerMessage{
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("traceparent"),
						Value: []byte("invalid"),
					},
				},
			},
		},
		"valid traceparent header": {
			msg: &sarama.ProducerMessage{
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("traceparent"),
						Value: []byte("00-1234567890abcdef-1234567890abcdef-01"),
					},
				},
			},
			expSpanID: "1234567890abcdef",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// When
			actualSpanID := ExtractSpanIDFromProducerMessage(tc.msg)

			// Then
			require.Equal(t, tc.expSpanID, actualSpanID)
		})
	}
}
