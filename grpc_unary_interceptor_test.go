package lit

//func Test_unaryServerInterceptor(t *testing.T) {
//	type mockServiceServer struct {
//		willPanic bool
//		in        *testdata.WeatherRequest
//		out       *testdata.WeatherResponse
//		inErr       error
//	}
//	tcs := map[string]struct {
//		givenRequest *testdata.WeatherRequest
//		mockSrv      mockServiceServer
//		expResult    *testdata.WeatherResponse
//		expErr       error
//		expLogs      []map[string]interface{}
//	}{
//		"success": {
//			givenRequest: &testdata.WeatherRequest{
//				Date: "M41.993.32",
//			},
//			mockSrv: mockServiceServer{
//				in: &testdata.WeatherRequest{
//					Date: "M41.993.32",
//				},
//				out: &testdata.WeatherResponse{
//					WeatherDetails: []*testdata.WeatherDetail{
//						{Location: "Hive City, Necromunda", Date: "M41.993.32", Description: "Toxic smog with occasional acid rain", Temperature: 42.7},
//					},
//				},
//			},
//			expResult: &testdata.WeatherResponse{
//				WeatherDetails: []*testdata.WeatherDetail{
//					{Location: "Hive City, Necromunda", Date: "M41.993.32", Description: "Toxic smog with occasional acid rain", Temperature: 42.7},
//				},
//			},
//			expLogs: []map[string]interface{}{
//				{
//					"grpc.request_body":   `{"date":"M41.993.32"}`,
//					"grpc.response_body":  `{"weatherDetails":[{"location":"Hive City, Necromunda","date":"M41.993.32","description":"Toxic smog with occasional acid rain","temperature":42.7}]}`,
//					"grpc.service_method": "/weather.WeatherService/GetWeatherInfo",
//					"level":               "info",
//					"msg":                 "grpc.unary_incoming_call",
//					"span_id":             "0000000000000001",
//					"trace_id":            "00000000000000000000000000000001",
//				},
//			},
//		},
//		"expected-error": {
//			givenRequest: &testdata.WeatherRequest{},
//			mockSrv: mockServiceServer{
//				in:  &testdata.WeatherRequest{},
//				inErr: errors.New("expected error"),
//			},
//			expErr: errors.New("expected error"),
//			expLogs: []map[string]interface{}{
//				{
//					"grpc.service_method": "/weather.WeatherService/GetWeatherInfo",
//					"level":               "info",
//					"msg":                 "grpc.unary_incoming_call",
//					"span_id":             "0000000000000001",
//					"trace_id":            "00000000000000000000000000000001",
//				},
//			},
//		},
//		"panic": {
//			givenRequest: &testdata.WeatherRequest{},
//			mockSrv: mockServiceServer{
//				willPanic: true,
//				in:        &testdata.WeatherRequest{},
//			},
//			expErr: ErrDefaultInternal,
//			expLogs: []map[string]interface{}{
//				{
//					"error":    "simulated panic",
//					"level":    "error",
//					"msg":      "Caught a panic: goroutine 7 [running]:\nruntime/debug.Stack()\n\t/Users/locdang/sdk/go1.23.3/src/runtime/debug/stack.go:26 +0x64\ngi",
//					"span_id":  "0000000000000001",
//					"trace_id": "00000000000000000000000000000001",
//				},
//			},
//		},
//	}
//
//	for scenario, tc := range tcs {
//		tc := tc
//		t.Run(scenario, func(t *testing.T) {
//			t.Parallel()
//
//			// Given
//			monitorTest, endTest := monitoring.NewMonitorTest()
//			defer endTest()
//
//			srv := new(weatherService)
//			if tc.mockSrv.willPanic {
//				srv.On("GetWeatherInfo", mock.Anything, tc.mockSrv.in).Panic("simulated panic")
//			} else {
//				srv.On("GetWeatherInfo", mock.Anything, tc.mockSrv.in).Return(tc.mockSrv.out, tc.mockSrv.inErr)
//			}
//			srvInfo := &grpc.UnaryServerInfo{
//				Server:     srv,
//				FullMethod: testdata.WeatherService_GetWeatherInfo_FullMethodName,
//			}
//
//			// When
//			intercept := unaryServerInterceptor(monitorTest.Context())
//			rs, inErr := intercept(context.Background(), tc.givenRequest, srvInfo, func(ctx context.Context, req interface{}) (interface{}, error) {
//				return srv.GetWeatherInfo(ctx, req.(*testdata.WeatherRequest))
//			})
//
//			// Then
//			if tc.expErr != nil {
//				require.EqualError(t, inErr, tc.expErr.Error())
//			} else {
//				require.NoError(t, inErr)
//				requireEqual(t, *tc.expResult, *rs.(*testdata.WeatherResponse),
//					cmpopts.IgnoreUnexported(testdata.WeatherResponse{}, testdata.WeatherDetail{}),
//				)
//			}
//
//			requireEqual(t, tc.expLogs, monitorTest.GetLogs(t),
//				cmpopts.IgnoreMapEntries(func(key string, value interface{}) bool {
//					return key == "ts" || key == "span_id" || strings.HasPrefix(value.(string), "Caught a panic")
//				}),
//				cmp.FilterPath(func(path cmp.Path) bool {
//					for _, ps := range path {
//						if mapIndex, ok := ps.(cmp.MapIndex); ok {
//							if key, ok := mapIndex.Key().Interface().(string); ok && key == "grpc.response_body" {
//								return true
//							}
//						}
//					}
//					return false
//				}, cmpopts.AcyclicTransformer("stripAllWhiteSpace", func(s string) string {
//					return strings.ReplaceAll(s, " ", "")
//				})), // grpc.response_body may vary due to protojson serialization inconsistencies
//			)
//		})
//	}
//}
