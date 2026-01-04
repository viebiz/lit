package grpcclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/viebiz/lit/grpcclient/testdata"
	"github.com/viebiz/lit/monitoring"
)

type testWeatherServer struct {
	testdata.UnimplementedWeatherServiceServer
}

func (s *testWeatherServer) GetWeatherInfo(ctx context.Context, req *testdata.WeatherRequest) (*testdata.WeatherResponse, error) {
	return &testdata.WeatherResponse{
		WeatherDetails: []*testdata.WeatherDetail{
			{
				Location:    "Hive City, Necromunda",
				Date:        req.GetDate(),
				Description: "Toxic smog with occasional acid rain",
				Temperature: 42.7,
			},
		},
	}, nil
}

func TestNewUnauthenticatedConnection_Success(t *testing.T) {
	// Given: start a gRPC server on an ephemeral port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srvAddr := lis.Addr().String()

	grpcServer := grpc.NewServer()
	testdata.RegisterWeatherServiceServer(grpcServer, &testWeatherServer{})

	done := make(chan struct{})
	go func() {
		_ = grpcServer.Serve(lis)
		close(done)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		<-done
	})

	// And: setup monitoring to capture logs from interceptor
	logBuffer := new(bytes.Buffer)
	m, mErr := monitoring.New(monitoring.Config{
		ServerName:  "lightning",
		Environment: "dev",
		Version:     "1.0.0",
		Writer:      logBuffer,
	})
	require.NoError(t, mErr)
	ctx := monitoring.SetInContext(context.Background(), m)

	// When: create unauthenticated connection and perform RPC
	conn, cErr := NewUnauthenticatedConnection(context.Background(), srvAddr)
	require.NoError(t, cErr)

	client := testdata.NewWeatherServiceClient(conn)
	req := &testdata.WeatherRequest{Date: "M41.993.32"}
	resp, rErr := client.GetWeatherInfo(ctx, req)

	// Then: RPC succeeds with expected response
	require.NoError(t, rErr)
	require.NotNil(t, resp)
	require.Len(t, resp.GetWeatherDetails(), 1)
	require.Equal(t, "M41.993.32", resp.GetWeatherDetails()[0].GetDate())

	// And: interceptor emitted outgoing request log with expected fields
	entry := findOutgoingRequestLogForTest(t, logBuffer)
	require.Equal(t, "grpc.outgoing_request", entry["msg"])
	require.Equal(t, `{"date":"M41.993.32"}`, entry["grpc.request"])
	require.Equal(t, "grpc", entry["rpc.system"])
	require.Equal(t, "weather.WeatherService", entry["rpc.service"])
	require.Equal(t, "GetWeatherInfo", entry["rpc.method"])
}

// findOutgoingRequestLogForTest parses the log buffer and returns the first entry for grpc.outgoing_request.
// It ignores initialization logs and blank lines.
func findOutgoingRequestLogForTest(t *testing.T, buf *bytes.Buffer) map[string]string {
	t.Helper()

	lines := strings.Split(buf.String(), "\n")
	for _, s := range lines {
		if strings.TrimSpace(s) == "" {
			continue
		}
		var m map[string]string
		if err := json.Unmarshal([]byte(s), &m); err != nil {
			continue
		}
		if m["msg"] == "grpc.outgoing_request" {
			return m
		}
	}
	t.Fatalf("did not find grpc.outgoing_request log entry in logs:\n%s", buf.String())
	return nil
}
