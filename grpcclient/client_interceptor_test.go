package grpcclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
		expLog   []map[string]interface{}
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
			expLog: []map[string]interface{}{
				{
					"grpc.request":      `{"date":"M41.993.32"}`,
					"level":             "info",
					"msg":               "grpc.outgoing_request",
					"outgoing_span_id":  "0000000000000000",
					"outgoing_trace_id": "00000000000000000000000000000000",
				},
			},
		},
	}
	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			logBuffer := new(bytes.Buffer)
			logger := monitoring.NewLoggerWithWriter(logBuffer)
			reqCtx := monitoring.SetInContext(context.Background(), logger)

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

			if diff := cmp.Diff(tc.expLog, parseLogs(t, *logBuffer), cmpopts.IgnoreMapEntries(func(key string, value interface{}) bool {
				if key == "ts" {
					return true
				}

				return false
			})); diff != "" {
				t.Errorf("unexpected log (-want, +got) = %v", diff)
			}
		})
	}
}

type weatherService struct {
	testdata.UnimplementedWeatherServiceServer
	mock.Mock
}

func (s *weatherService) GetWeatherInfo(ctx context.Context, req *testdata.WeatherRequest) (*testdata.WeatherResponse, error) {
	args := s.Called(ctx, req)

	return args.Get(0).(*testdata.WeatherResponse), args.Error(1)
}

func parseLogs(t require.TestingT, buf bytes.Buffer) []map[string]interface{} {
	var logs []map[string]interface{}

	lines := bytes.Split(buf.Bytes(), []byte("\n")) // \n is end of line
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var msg map[string]interface{}
		require.NoError(t, json.Unmarshal(line, &msg))
		logs = append(logs, msg)
	}

	return logs
}
