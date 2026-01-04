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
	"google.golang.org/grpc/credentials/insecure"

	"github.com/viebiz/lit/grpcclient/testdata"
	"github.com/viebiz/lit/monitoring"
)

// weatherSvc implements the WeatherService for testing
type weatherSvc struct {
	testdata.UnimplementedWeatherServiceServer
}

func (s *weatherSvc) GetWeatherInfo(ctx context.Context, req *testdata.WeatherRequest) (*testdata.WeatherResponse, error) {
	// Echo back a response using the incoming request date so we can assert through the client
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

func TestDialOptions_ExternalServiceInfoAndInterceptor(t *testing.T) {
	// Given: start a gRPC server
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srvAddr := lis.Addr().String()

	grpcServer := grpc.NewServer()
	testdata.RegisterWeatherServiceServer(grpcServer, &weatherSvc{})

	done := make(chan struct{})
	go func() {
		_ = grpcServer.Serve(lis)
		close(done)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		<-done
	})

	// And: prepare monitoring to capture logs
	logBuffer := new(bytes.Buffer)
	m, err := monitoring.New(monitoring.Config{
		ServerName:  "lightning",
		Environment: "dev",
		Version:     "1.0.0",
		Writer:      logBuffer,
	})
	require.NoError(t, err)
	ctx := monitoring.SetInContext(context.Background(), m)

	// When: create a client conn using the common dial options and call the RPC
	svcInfo := monitoring.NewExternalServiceInfo(srvAddr)
	conn, err := grpc.NewClient(
		srvAddr,
		append(
			commonUnaryClientDialOptions(svcInfo),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)...,
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client := testdata.NewWeatherServiceClient(conn)
	req := &testdata.WeatherRequest{Date: "M41.993.32"}
	resp, err := client.GetWeatherInfo(ctx, req)

	// Then: RPC succeeds
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetWeatherDetails(), 1)
	require.Equal(t, "M41.993.32", resp.GetWeatherDetails()[0].GetDate())

	// And: the interceptor ran with the external service info option applied
	entry := findOutgoingRequestLog(t, logBuffer)
	require.Equal(t, "grpc.outgoing_request", entry["msg"])
	require.Equal(t, `{"date":"M41.993.32"}`, entry["grpc.request"])
	require.Equal(t, "grpc", entry["rpc.system"])
	require.Equal(t, "weather.WeatherService", entry["rpc.service"])
	require.Equal(t, "GetWeatherInfo", entry["rpc.method"])
	require.Equal(t, svcInfo.Hostname+":"+svcInfo.Port, entry["server.address"])
}

// findOutgoingRequestLog parses the log buffer and returns the first entry for grpc.outgoing_request.
// It ignores the initial initialization logs and any empty lines.
func findOutgoingRequestLog(t *testing.T, buf *bytes.Buffer) map[string]string {
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
