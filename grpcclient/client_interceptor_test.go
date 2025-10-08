package grpcclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/testutil"
	"google.golang.org/grpc"

	"github.com/viebiz/lit/grpcclient/testdata"
	"github.com/viebiz/lit/monitoring"
)

func TestClientConn_Invoke(t *testing.T) {
	const srvAddr = "localhost:50052"
	type mockData struct {
		inReq  *testdata.WeatherRequest
		outRes *testdata.WeatherResponse
		outErr error
	}

	tcs := map[string]struct {
		givenReq *testdata.WeatherRequest
		mockData mockData
		expResp  *testdata.WeatherResponse
		expErr   error
		expLog   []map[string]string
	}{
		"success": {
			givenReq: &testdata.WeatherRequest{
				Date: "M41.993.32",
			},
			mockData: mockData{
				inReq: &testdata.WeatherRequest{
					Date: "M41.993.32",
				},
				outRes: &testdata.WeatherResponse{
					WeatherDetails: []*testdata.WeatherDetail{
						{
							Location:    "Hive City, Necromunda",
							Date:        "M41.993.32",
							Description: "Toxic smog with occasional acid rain",
							Temperature: 42.7,
						},
						{
							Location:    "Macragge's Northern Hemisphere",
							Date:        "M41.874.21",
							Description: "Freezing winds with snowstorms",
							Temperature: -20.5,
						},
					},
				},
			},
			expResp: &testdata.WeatherResponse{
				WeatherDetails: []*testdata.WeatherDetail{
					{
						Location:    "Hive City, Necromunda",
						Date:        "M41.993.32",
						Description: "Toxic smog with occasional acid rain",
						Temperature: 42.7,
					},
					{
						Location:    "Macragge's Northern Hemisphere",
						Date:        "M41.874.21",
						Description: "Freezing winds with snowstorms",
						Temperature: -20.5,
					},
				},
			},
			expLog: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "grpc.outgoing_request", "grpc.request": `{"date":"M41.993.32"}`, "outgoing_span_id": "0000000000000000", "outgoing_trace_id": "00000000000000000000000000000000", "rpc.method": "GetWeatherInfo", "rpc.service": "weather.WeatherService", "rpc.system": "grpc", "server.address": "localhost:50052", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
			},
		},
	}
	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			logBuffer := new(bytes.Buffer)
			m, _ := monitoring.New(monitoring.Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: logBuffer})
			reqCtx := monitoring.SetInContext(context.Background(), m)

			// Start a new GRPC server for testing
			go func() {
				weatherSvc := new(weatherService)
				weatherSvc.On("GetWeatherInfo", mock.Anything, tc.mockData.inReq).
					Return(tc.mockData.outRes, tc.mockData.outErr)

				lis, err := net.Listen("tcp", srvAddr)
				require.NoError(t, err)

				grpcServer := grpc.NewServer()
				testdata.RegisterWeatherServiceServer(grpcServer, weatherSvc)

				require.NoError(t, grpcServer.Serve(lis))
			}()

			// When
			conn, err := NewUnauthenticatedConnection(context.Background(), srvAddr)
			require.NoError(t, err)

			weatherClient := testdata.NewWeatherServiceClient(conn)
			resp, err := weatherClient.GetWeatherInfo(reqCtx, tc.givenReq)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				testutil.Equal(t, tc.expResp, resp, testutil.IgnoreUnexported[*testdata.WeatherResponse](testdata.WeatherResponse{}, testdata.WeatherDetail{}))
			}

			pasedLogs, err := parseLog(logBuffer.Bytes(), 2)
			require.NoError(t, err)
			testutil.Equal(t, tc.expLog, pasedLogs, testutil.IgnoreSliceMapEntries(func(k string, v string) bool {
				if k == "ts" {
					return true
				}

				if k == "error.stack" {
					return true
				}

				if v == "Caught a panic" {
					return true
				}

				return false
			}))
		})
	}
}

func parseLog(b []byte, skip int) ([]map[string]string, error) {
	var result []map[string]string
	for idx, s := range strings.Split(string(b), "\n") {
		if s == "" {
			break
		}
		if idx < skip {
			continue
		}
		var r map[string]string
		if err := json.Unmarshal([]byte(s), &r); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

type weatherService struct {
	testdata.UnimplementedWeatherServiceServer
	mock.Mock
}

func (s *weatherService) GetWeatherInfo(ctx context.Context, req *testdata.WeatherRequest) (*testdata.WeatherResponse, error) {
	args := s.Called(ctx, req)

	return args.Get(0).(*testdata.WeatherResponse), args.Error(1)
}

func TestUnaryClientInterceptor_Success(t *testing.T) {
	t.Parallel()

	// Given: monitoring context to capture logs
	logBuffer := new(bytes.Buffer)
	m, err := monitoring.New(monitoring.Config{
		ServerName:  "lightning",
		Environment: "dev",
		Version:     "1.0.0",
		Writer:      logBuffer,
	})
	require.NoError(t, err)
	ctx := monitoring.SetInContext(context.Background(), m)

	// Given: request/reply and method
	method := "/weather.WeatherService/GetWeatherInfo"
	req := &testdata.WeatherRequest{Date: "M41.993.32"}
	var reply testdata.WeatherResponse

	// And: external service info option so interceptor enriches context and logs
	svcInfo := monitoring.NewExternalServiceInfo("127.0.0.1:4317")
	callOpt := externalServiceInfoOption{info: svcInfo}

	// And: invoker that simulates a successful RPC
	invoker := func(ctx context.Context, _ string, _ interface{}, reply interface{}, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		// Write some response to verify reply is passed-through correctly
		wr := reply.(*testdata.WeatherResponse)
		wr.WeatherDetails = []*testdata.WeatherDetail{
			{
				Location:    "Hive City, Necromunda",
				Date:        "M41.993.32",
				Description: "Toxic smog with occasional acid rain",
				Temperature: 42.7,
			},
		}
		return nil
	}

	// When
	err = unaryClientInterceptor(ctx, method, req, &reply, nil, invoker, callOpt)

	// Then: interceptor returns nil and reply is populated by invoker
	require.NoError(t, err)
	require.Len(t, reply.GetWeatherDetails(), 1)
	require.Equal(t, "M41.993.32", reply.GetWeatherDetails()[0].GetDate())

	// And: interceptor logged outgoing request with expected fields
	entry := findOutgoingRequestLogFromBuffer(t, logBuffer)
	require.Equal(t, "grpc.outgoing_request", entry["msg"])
	require.Equal(t, `{"date":"M41.993.32"}`, entry["grpc.request"])
	require.Equal(t, "grpc", entry["rpc.system"])
	require.Equal(t, "weather.WeatherService", entry["rpc.service"])
	require.Equal(t, "GetWeatherInfo", entry["rpc.method"])
	require.Equal(t, svcInfo.Hostname+":"+svcInfo.Port, entry["server.address"])
}

func TestUnaryClientInterceptor_Error(t *testing.T) {
	t.Parallel()

	// Given: monitoring context to capture logs
	logBuffer := new(bytes.Buffer)
	m, err := monitoring.New(monitoring.Config{
		ServerName:  "lightning",
		Environment: "dev",
		Version:     "1.0.0",
		Writer:      logBuffer,
	})
	require.NoError(t, err)
	ctx := monitoring.SetInContext(context.Background(), m)

	// Given: request/reply and method
	method := "/weather.WeatherService/GetWeatherInfo"
	req := &testdata.WeatherRequest{Date: "M41.874.21"}
	var reply testdata.WeatherResponse

	// And: external service info option so interceptor enriches context and logs
	svcInfo := monitoring.NewExternalServiceInfo("localhost:9999")
	callOpt := externalServiceInfoOption{info: svcInfo}

	// And: invoker that simulates an RPC error
	expErr := errors.New("rpc failed")
	invoker := func(ctx context.Context, _ string, _ interface{}, _ interface{}, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		return expErr
	}

	// When
	err = unaryClientInterceptor(ctx, method, req, &reply, nil, invoker, callOpt)

	// Then: interceptor returns error from invoker
	require.EqualError(t, err, expErr.Error())

	// And: interceptor still logged outgoing request with expected fields
	entry := findOutgoingRequestLogFromBuffer(t, logBuffer)
	require.Equal(t, "grpc.outgoing_request", entry["msg"])
	require.Equal(t, `{"date":"M41.874.21"}`, entry["grpc.request"])
	require.Equal(t, "grpc", entry["rpc.system"])
	require.Equal(t, "weather.WeatherService", entry["rpc.service"])
	require.Equal(t, "GetWeatherInfo", entry["rpc.method"])
	require.Equal(t, svcInfo.Hostname+":"+svcInfo.Port, entry["server.address"])
}

// findOutgoingRequestLogFromBuffer parses the log buffer and returns the first entry for grpc.outgoing_request.
// It ignores initialization logs and blank lines.
func findOutgoingRequestLogFromBuffer(t *testing.T, buf *bytes.Buffer) map[string]string {
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
