package grpcclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
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
				if diff := cmp.Diff(*tc.expResp, *resp, cmpopts.IgnoreUnexported(testdata.WeatherResponse{}, testdata.WeatherDetail{})); diff != "" {
					t.Errorf("unexpected response (-want, +got) = %v", diff)
				}
			}

			pasedLogs, err := parseLog(logBuffer.Bytes())
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

func parseLog(b []byte) ([]map[string]string, error) {
	var result []map[string]string
	for _, s := range strings.Split(string(b), "\n") {
		if s == "" {
			break
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
