package instrumentgrpc

//func TestStartUnaryIncomingCall(t *testing.T) {
//	// Given
//	tp := mocktracer.Start()
//	defer tp.Stop()
//
//	logBuffer := bytes.NewBuffer(nil)
//	m, err := monitoring.New(monitoring.Config{Writer: logBuffer})
//	require.NoError(t, err)
//
//	reqCtx := context.Background()
//	reqCtx = peer.NewContext(reqCtx, &peer.Peer{
//		Addr: &net.TCPAddr{
//			Port: 50051,
//		},
//	})
//
//	reqCtx = metadata.NewIncomingContext(reqCtx, metadata.New(map[string]string{
//		"traceparent": "00-deadbeefcafebabefeedfacebadc0de1-abad1dea0ddba11c-01",
//		"tracestate":  "test=test-value",
//		"baggage":     "user_id=1234,role=admin",
//	}))
//
//	// When
//	ctx, reqMeta, end := StartUnaryIncomingCall(reqCtx, m, "/weather.WeatherService/GetWeatherInfo", &testdata.WeatherRequest{
//		Date: "M41.993.32",
//	})
//
//	// Then
//	monitoring.FromContext(ctx).Infof("Got incoming request")
//	expectedLogs := []map[string]interface{}{
//		{
//			"level":    "info",
//			"msg":      "Got incoming request",
//			"span_id":  "0000000000000000",                 // Random generated value
//			"trace_id": "deadbeefcafebabefeedfacebadc0de1", // Should sample with incoming request
//		},
//	}
//	requireEqual(t, expectedLogs, monitor.GetLogs(t), cmpopts.IgnoreMapEntries(func(key string, value interface{}) bool {
//		return key == "ts" || key == "span_id"
//	}))
//
//	// 3.2. Validate request metadata
//	requireEqual(t, RequestMetadata{
//		ServiceMethod: "/weather.WeatherService/GetWeatherInfo",
//		BodyToLog:     []uint8(`{"date":"M41.993.32"}`),
//	}, reqMeta, cmpopts.IgnoreFields(RequestMetadata{}, "ContextData"))
//	require.ElementsMatch(t, []string{"user_id=1234", "role=admin"}, reqMeta.ContextData)
//
//	// 3.3. Simulated end instrument
//	require.NotNil(t, end)
//	end(nil)
//
//	// 3.4. Validate trace attributes
//	expectedAttributes := []attribute.KeyValue{
//		semconv.RPCSystemGRPC,
//		semconv.NetworkPeerAddress(":50051"),
//		semconv.NetworkTransportTCP,
//		semconv.RPCService("weather.WeatherService"),
//		semconv.RPCMethod("GetWeatherInfo"),
//		semconv.RPCGRPCStatusCodeOk,
//	}
//	require.ElementsMatch(t, expectedAttributes, monitor.GetSpans().Snapshots()[0].Attributes())
//}
