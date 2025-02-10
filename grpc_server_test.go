package lightning

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/viebiz/lit/grpcclient/testdata"
)

func TestGRPCServer_UnaryRun(t *testing.T) {
	const srvAddr = "localhost:50051"
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
		},
	}
	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			// Start a new GRPC server for testing
			go func(ctx context.Context) {
				weatherSvc := new(weatherService)
				weatherSvc.On("GetWeatherInfo", mock.Anything, tc.mockData.inReq).
					Return(tc.mockData.outRes, tc.mockData.outErr)

				srv, err := NewGRPCServerWithOptions(ctx, srvAddr)
				require.NoError(t, err)
				testdata.RegisterWeatherServiceServer(srv.Registrar(), weatherSvc)

				require.NoError(t, srv.Run())
			}(context.Background())

			// When
			// Create a client connection to gRPC server
			conn, err := grpc.NewClient(srvAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()

			weatherClient := testdata.NewWeatherServiceClient(conn)
			resp, err := weatherClient.GetWeatherInfo(context.Background(), tc.givenReq)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				requireEqual(t, *tc.expResp, *resp,
					cmpopts.IgnoreUnexported(testdata.WeatherResponse{}, testdata.WeatherDetail{}),
				)
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
